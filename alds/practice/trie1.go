package practice

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
//  PUT   |   /users/#/locations | 204  ||  4
// POST   |               /rides | 202  ||  3
//  GET   |             /users/# | 200  ||  2
// POST   |               /users | 202  ||  2
//  PUT   |   /users/#/locations | 500  ||  1
//  GET   |      /prices/#/geo/# | 200  ||  1
//  PUT   |   /users/#/locations | 404  ||  1
//  GET   |      /rides/#/status | 200  ||  1
//  GET   |      /users/#/ride/# | 200  ||  1
// ```

type Log struct {
	EndPoint string
	Method   string
	Status   string
}

func (l *Log) GetKey() string {
	return fmt.Sprintf("%s-%s", l.Method, l.Status)
}

func (l *Log) GetParts() []string {
	return strings.Split(l.EndPoint, "/")
}

func NewLogFromKey(key string, endpoint string) *Log {
	keys := strings.Split(key, "-")
	return &Log{
		EndPoint: endpoint,
		Method:   keys[0],
		Status:   keys[1],
	}
}

// todo: error handling logic is omitted for the first version, add later
// so that we can focus on the core logic here

func logParser(line string) *Log {
	line = strings.Split(line, `]`)[1]
	parts := strings.Fields(line)
	log := &Log{
		Method:   parts[0],
		EndPoint: parsePath(parts[1]),
		Status:   parts[3],
	}
	return log
}

var isDigits = regexp.MustCompile(`^\d+$`)

func parsePath(rawPath string) string {
	parts := strings.Split(rawPath, `/`)
	for i, part := range parts {
		if isDigits.MatchString(part) {
			parts[i] = "#"
		}
	}
	return strings.Join(parts, "/")
}

// build trie tree to count
type LogNode struct {
	PathPart string
	Children map[string]*LogNode // key is the pathpart of the child
	Counts   map[string]int      // key is the log's key
}

func (ln *LogNode) GetChild(partName string) *LogNode {
	return ln.Children[partName]
}

func (ln *LogNode) AddChild(partName string) *LogNode {
	ln.Children[partName] = NewLogNode(partName)
	return ln.Children[partName]
}

func NewLogNode(partName string) *LogNode {
	return &LogNode{
		PathPart: partName,
		Children: make(map[string]*LogNode),
		Counts:   make(map[string]int),
	}
}

// couting
func ProcessOneLog(root *LogNode, line string) {
	node := root
	log := logParser(line)
	// fmt.Println("the log's path", strings.Join(log.PathParts, "/"))
	for _, part := range log.GetParts() {
		child := node.GetChild(part)
		if child == nil {
			child = node.AddChild(part)
		}
		node = child
	}
	node.Counts[log.GetKey()] += 1
}

func ProcessLogs(root *LogNode, lines []string) {
	for _, line := range lines {
		ProcessOneLog(root, line)
	}
}

// print
type Stat struct {
	*Log
	Count int
}

// Method |             Endpoint | Code || Count
// =============================================
//
//	PUT   |   /users/#/locations | 204  ||  4
//
// POST   |               /rides | 202  ||  3
func (s *Stat) Format() string {
	return "%-15s|%-50s|%-10s||%-10v"
}

func (s *Stat) Line() string {
	return strings.Repeat("=", 50+20+15)
}

func (s *Stat) String() string {
	return fmt.Sprintf(s.Format(), s.Method, s.EndPoint, s.Status, s.Count)
}

func (s *Stat) Header() string {
	return fmt.Sprintf(s.Format(), "Method", "Endpoint", "Code", "Count")
}

func CollectStats(root *LogNode) []*Stat {
	// collect stats
	parts := make([]string, 0, 16)
	stats := make([]*Stat, 0, 16)
	var recur func(node *LogNode)
	recur = func(node *LogNode) {
		if node == nil {
			return
		}
		parts = append(parts, node.PathPart)
		// add stats of current node
		if len(node.Counts) > 0 {
			for k, v := range node.Counts {
				endpoint := strings.Join(parts[1:], "/")
				log := NewLogFromKey(k, endpoint) // problem
				stat := &Stat{
					Log:   log,
					Count: v,
				}
				stats = append(stats, stat)
			}
		}
		for _, child := range node.Children {
			recur(child)
		}
		parts = parts[:len(parts)-1] // pop out the current pathPart
	}
	recur(root)
	// sort stats
	sort.Slice(stats, func(i, j int) bool { return stats[i].Count > stats[j].Count })
	return stats
}

func PrintStats(stats []*Stat) {
	for idx, stat := range stats {
		if idx == 0 {
			fmt.Println(stat.Header())
			fmt.Println(stat.Line())
		}
		fmt.Println(stat.String())
	}
}
