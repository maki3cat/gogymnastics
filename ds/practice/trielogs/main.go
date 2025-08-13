package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ```
// Web server log sample:
// [02/Nov/2018:21:46:31 +0000] PUT /users/12345/locations HTTP/1.1 204 iphone-3
// [02/Nov/2018:21:46:31 +0000] PUT /users/6098/locations HTTP/1.1 204 iphone-3
// [02/Nov/2018:21:46:32 +0000] PUT /users/3911/locations HTTP/1.1 204 moto-x
// [02/Nov/2018:21:46:33 +0000] PUT /users/9933/locations HTTP/1.1 404 moto-x
// [02/Nov/2018:21:46:33 +0000] PUT /users/3911/locations HTTP/1.1 500 moto-x
// [02/Nov/2018:21:46:34 +0000] GET /rides/9943222/status HTTP/1.1 200 moto-x
// [02/Nov/2018:21:46:34 +0000] POST /rides HTTP/1.1 202 iphone-2
// [02/Nov/2018:21:46:35 +0000] POST /users HTTP/1.1 202 iphone-5
// [02/Nov/2018:21:46:35 +0000] POST /rides HTTP/1.1 202 iphone-5
// [02/Nov/2018:21:46:37 +0000] POST /rides HTTP/1.1 202 iphone-4
// [02/Nov/2018:21:46:38 +0000] GET /users/994/ride/16 HTTP/1.1 200 iphone-5
// [02/Nov/2018:21:46:39 +0000] POST /users HTTP/1.1 202 iphone-3
// [02/Nov/2018:21:46:40 +0000] PUT /users/8384721/locations HTTP/1.1 204 iphone-3
// [02/Nov/2018:21:46:41 +0000] GET /users/342111 HTTP/1.1 200 iphone-5
// [02/Nov/2018:21:46:42 +0000] GET /users/9933 HTTP/1.1 200 iphone-5
// [02/Nov/2018:21:46:43 +0000] GET /prices/20180103/geo/12 HTTP/1.1 200 iphone-5

// Take a log and output a table representing the number of occurrences of events in the
// log.

// Where the event is defined as (Method + Endpoint + HttpStatusCode).

// Order by Count, descending:

// Method |             Endpoint | Code || Count
// =============================================
//
//	PUT   |   /users/#/locations | 204  ||  4
//
// POST   |               /rides | 202  ||  3
//
//	GET   |             /users/# | 200  ||  2
//
// POST   |               /users | 202  ||  2
//
//	PUT   |   /users/#/locations | 500  ||  1
//	GET   |      /prices/#/geo/# | 200  ||  1
//	PUT   |   /users/#/locations | 404  ||  1
//	GET   |      /rides/#/status | 200  ||  1
//	GET   |      /users/#/ride/# | 200  ||  1
//
// ```

// parse the log
type Log struct {
	method   string
	endpoint string
	code     string
}

func (l *Log) GetKey() string {
	return fmt.Sprintf("%s-%s", l.method, l.code)
}

func (l *Log) GetEndpointParts() []string {
	return strings.Split(l.endpoint, "/")
}

func NewLog(key string, endpoint string) *Log {
	log := new(Log)
	log.endpoint = endpoint
	parts := strings.Split(key, "-")
	log.method = parts[0]
	log.code = parts[1]
	return log
}

// todo: skip the error handling for now, focus on main logic
// [02/Nov/2018:21:46:40 +0000] PUT /users/8384721/locations HTTP/1.1 204 iphone-3
func parseLog(line string) *Log {
	log := new(Log) // maki: shall make good use of new to skip the verbose of using &Log{}
	parts := strings.Split(line, "]")
	contents := strings.Fields(parts[1])
	log.method = contents[0]
	log.endpoint = contents[1]
	log.code = contents[3]
	// parse the real path
	log.endpoint = parseEndpoint(log.endpoint)
	fmt.Println(log.endpoint)
	return log
}

var isDigits = regexp.MustCompile(`^\d+$`)

func parseEndpoint(rawPath string) string {
	parts := strings.Split(rawPath, "/")
	for idx, part := range parts {
		if isDigits.MatchString(part) {
			parts[idx] = "#"
		}
	}
	return strings.Join(parts, "/")
}

// count the log, using trie
type TNode struct {
	part     string
	children map[string]*TNode
	counts   map[string]int
}

func NewTNode(part string) *TNode {
	node := new(TNode)
	node.part = part
	node.children = make(map[string]*TNode)
	node.counts = make(map[string]int)
	return node
}

func (tn *TNode) GetChild(part string) *TNode {
	return tn.children[part]
}

func (tn *TNode) AddChild(part string) *TNode {
	node := NewTNode(part)
	tn.children[part] = node
	return node
}

func CountLog(root *TNode, log *Log) {
	node := root
	parts := log.GetEndpointParts()
	for _, part := range parts {
		child := node.GetChild(part)
		if child == nil {
			child = node.AddChild(part)
		}
		node = child
	}
	node.counts[log.GetKey()] += 1
}

// collect stat
type Stat struct {
	*Log
	count int
}

func (s *Stat) Header() string {
	return strings.Repeat("=", 100)
}

// [02/Nov/2018:21:46:33 +0000] PUT /users/9933/locations HTTP/1.1 404 moto-x
func (s *Stat) FormatLine() string {
	return fmt.Sprintf("|%-15s|%-30s|%-15s||%-15v|", s.method, s.endpoint, s.code, s.count)
}

func CollectStat(root *TNode) []*Stat {
	parts := make([]string, 0)
	stats := make([]*Stat, 0)
	var recur func(node *TNode)
	recur = func(node *TNode) {
		if node == nil {
			return
		}
		parts = append(parts, node.part)
		// count the current node
		if len(node.counts) > 0 {
			for k, v := range node.counts {
				endpoint := strings.Join(parts[1:], "/")
				log := NewLog(k, endpoint)
				stat := new(Stat)
				stat.Log = log
				stat.count = v
				stats = append(stats, stat)
			}
		}
		// go on couting the children
		for _, child := range node.children {
			recur(child)
		}
		parts = parts[:len(parts)-1]
	}
	recur(root)
	return stats
}

// process logs
func ProcessLogs(lines []string) {
	root := NewTNode("/")
	// count
	for _, line := range lines {
		log := parseLog(line)
		CountLog(root, log)
	}
	// collect stats and sort
	stats := CollectStat(root)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})
	// print
	for idx, stat := range stats {
		if idx == 0 {
			fmt.Println(stat.Header())
		}
		fmt.Println(stat.FormatLine())
	}
}
