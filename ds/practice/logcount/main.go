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
// [02/Nov/2018:21:46:34 +0000] POST /rides HTTP/1.1 202 iphone-2
// [02/Nov/2018:21:46:35 +0000] POST /users HTTP/1.1 202 iphone-5

// part 1: turn a log string into the structure
type LogLine struct {
	Method   string
	Endpoint string
	Code     string
	// maki: we merge the stat and logline together here
	Count int
}

func (ll *LogLine) Format() string {
	// maki: without -, it will be right alignment
	return fmt.Sprintf("%-10s|%-30s|%-10s||%-5d", ll.Method, ll.Endpoint, ll.Code, ll.Count)
}

func (ll *LogLine) GetCountingKey() string {
	return strings.Join([]string{ll.Method, ll.Code}, "-")
}

func (ll *LogLine) GetParts() []string {
	return strings.Split(ll.Endpoint, `/`)
}

func FromCountingKey(key string, endpointParts []string, count int) *LogLine {
	parts := strings.Split(key, "-")
	logline := new(LogLine)
	logline.Method = parts[0]
	logline.Code = parts[1]
	logline.Count = count
	// maki: this is important we use a unchangeable string
	// instead of the endpointParts whose pointer may be shared outside the current scope
	logline.Endpoint = strings.Join(endpointParts, `/`)
	return logline
}

func Parse(line string) *LogLine {
	line = strings.Split(line, `]`)[1]
	line = strings.Trim(line, " ")
	parts := strings.Fields(line)
	logline := new(LogLine)
	logline.Method = parts[0]
	logline.Code = parts[3]
	logline.Count = 1
	logline.Endpoint = InterpretEndpoint(parts[1])
	return logline
}

var digitRegexp = regexp.MustCompile(`^\d+$`)

func isDigit(s string) bool {
	return digitRegexp.MatchString(s)
}

func InterpretEndpoint(path string) string {
	parts := strings.Split(path, `/`)
	for idx, part := range parts {
		if isDigit(part) {
			parts[idx] = `#`
		}
	}
	return strings.Join(parts, `/`)
}

// part 2: couting the basic structure
type Node struct {
	EndpointPart string
	Children     map[string]*Node
	Count        map[string]int
}

func NewNode(endpointPart string) *Node {
	n := new(Node)
	n.EndpointPart = endpointPart
	n.Children = make(map[string]*Node)
	n.Count = make(map[string]int, 0)
	return n
}

func (n *Node) GetChild(endpointPart string) *Node {
	return n.Children[endpointPart]
}

// maki: like the linkedlist, always return the element we add/push
func (n *Node) AddChild(endpointPart string) *Node {
	node := NewNode(endpointPart)
	n.Children[endpointPart] = node
	return node
}

// maki: new design we merge the method here
func (n *Node) GetOrAddChild(endpointPart string) *Node {
	child := n.GetChild(endpointPart)
	if child != nil {
		return child
	}
	return n.AddChild(endpointPart)
}

// part3: the counting api
func CountLogline(root *Node, log *LogLine) {
	node := root
	for _, part := range log.GetParts() {
		node = node.GetOrAddChild(part)
	}
	node.Count[log.GetCountingKey()] += 1
}

func CountLogs(root *Node, logs []string) {
	for _, line := range logs {
		logline := Parse(line)
		CountLogline(root, logline)
	}
}

// part4: retrieving the stats
// returns: the sorted logline with correct couting
// maki: pass the root in is better than using a global var
// maki: even if we use a global var, we can still pass it in
func RetrieveStats(root *Node) (stats []*LogLine) {

	// maki:
	// we have 2 in-function nonlocal var for the closure, stats and endpoint
	// this is the best part of dfs
	endpointCache := make([]string, 0)
	// we use dfs for retrieval
	var recur func(*Node)
	recur = func(node *Node) {
		if node == nil {
			return
		}
		endpointCache = append(endpointCache, node.EndpointPart)
		if len(node.Count) > 0 {
			for key, count := range node.Count {
				// maki: pay attention to the dummy
				stat := FromCountingKey(key, endpointCache[1:], count)
				stats = append(stats, stat)
			}
		}
		for _, child := range node.Children {
			recur(child)
		}
		// maki: dfs, don't forget to pop out the stack
		// pop out the current node
		endpointCache = endpointCache[:len(endpointCache)-1]
	}
	recur(root)
	// sort the stats
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})
	return
}

// todo: we skip the error handling and data integrity checking for now
func main() {
	logs := []string{"[02/Nov/2018:21:46:31 +0000] PUT /users/12345/locations HTTP/1.1 204 iphone-3",
		"[02/Nov/2018:21:46:31 +0000] PUT /users/6098/locations HTTP/1.1 204 iphone-3",
		"[02/Nov/2018:21:46:32 +0000] PUT /users/3911/locations HTTP/1.1 204 moto-x",
		"[02/Nov/2018:21:46:33 +0000] PUT /users/9933/locations HTTP/1.1 404 moto-x",
		"[02/Nov/2018:21:46:33 +0000] PUT /users/3911/locations HTTP/1.1 500 moto-x",
		"[02/Nov/2018:21:46:34 +0000] GET /rides/9943222/status HTTP/1.1 200 moto-x",
		"[02/Nov/2018:21:46:34 +0000] POST /rides HTTP/1.1 202 iphone-2",
		"[02/Nov/2018:21:46:35 +0000] POST /users HTTP/1.1 202 iphone-5",
		"[02/Nov/2018:21:46:35 +0000] POST /rides HTTP/1.1 202 iphone-5",
		"[02/Nov/2018:21:46:37 +0000] POST /rides HTTP/1.1 202 iphone-4",
		"[02/Nov/2018:21:46:38 +0000] GET /users/994/ride/16 HTTP/1.1 200 iphone-5",
		"[02/Nov/2018:21:46:39 +0000] POST /users HTTP/1.1 202 iphone-3",
		"[02/Nov/2018:21:46:40 +0000] PUT /users/8384721/locations HTTP/1.1 204 iphone-3",
		"[02/Nov/2018:21:46:41 +0000] GET /users/342111 HTTP/1.1 200 iphone-5",
		"[02/Nov/2018:21:46:42 +0000] GET /users/9933 HTTP/1.1 200 iphone-5",
		"[02/Nov/2018:21:46:43 +0000] GET /prices/20180103/geo/12 HTTP/1.1 200 iphone-5"}
	root := NewNode("/")
	CountLogs(root, logs)
	stats := RetrieveStats(root)
	for _, stat := range stats {
		fmt.Println(stat.Format())
	}
}
