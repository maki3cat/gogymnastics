package logcounter

import (
	"fmt"
	"regexp"
	"strings"
)

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

// ------------------------------parsing the log------------------------------

type Log struct {
	Method string
	Path   string
	Status string
	Count  int
}

// under the same path, the key can be different
func (l *Log) GetKey() string {
	return fmt.Sprintf("%s-%s", l.Method, l.Status)
}

func (l *Log) GetPathParts() []string {
	return strings.Split(l.Path, "/")
}

func ReconstructLog(key string, path string) *Log {
	parts := strings.Split(key, "-")
	return &Log{
		Method: parts[0],
		Path:   path,
		Status: parts[1],
	}
}

// assumes the log is valid
// todo: skip for now add errors when the log is corrupt data
func parseLog(log string) *Log {
	timestampLen := len("[02/Nov/2018:21:46:31 +0000]")
	log = log[timestampLen:]
	log = strings.TrimSpace(log)
	parts := strings.Split(log, " ")

	rawPath := parts[1]
	realPath := findRealPath(rawPath)

	return &Log{
		Method: parts[0],
		Path:   realPath,
		Status: parts[3],
	}
}

// var re = regexp.MustCompile("^[0-9]+$") // should use mustcompile as the global var
// todo: what is the difference between `^\d+$` and ""
var re2 = regexp.MustCompile(`^\d+$`) // should use mustcompile as the global var

func findRealPath(rawPath string) string {
	parts := strings.Split(rawPath, "/")
	for idx, part := range parts {
		if re2.MatchString(part) {
			parts[idx] = "#"
		}
	}
	return strings.Join(parts, "/")
}

// ------------------------------counting the log------------------------------
type Node struct {
	Part     string
	Counts   map[string]int   // key = log.GetKey()
	Children map[string]*Node // key = the part of the url
}

func (n *Node) GetChild(part string) *Node {
	return n.Children[part]
}

func (n *Node) AddChild(part string) *Node {
	n.Children[part] = NewNode(part)
	return n.Children[part]
}

func NewNode(part string) *Node {
	return &Node{
		Part:     part,
		Counts:   make(map[string]int), // key = log.GetKey()
		Children: make(map[string]*Node),
	}
}

// build the trie-tree
// root is the dummy one
func CountLog(root *Node, log *Log) {
	parts := log.GetPathParts()
	node := root
	for _, part := range parts {
		child := node.GetChild(part)
		if child == nil {
			child = node.AddChild(part)
		}
		node = child
	}
	node.Counts[log.GetKey()]++
}

// ----------------printing and formatting----------------
type LogStat struct {
	Log   *Log
	Count int
}

func NewLogStat(log *Log, count int) *LogStat {
	return &LogStat{
		Log:   log,
		Count: count,
	}
}

func (l *LogStat) String() string {
	return fmt.Sprintf("%-10s | %-30s | %-10s | %-5d", l.Log.Method, l.Log.Path, l.Log.Status, l.Count)
}

func (l *LogStat) Header() string {
	return fmt.Sprintf("%-10s | %-30s | %-10s | %-5s", "Method", "Path", "Status", "Count")
}

// printing the trie-tree
func PrintCounts(root *Node) {
	if root == nil {
		return
	}
	// header
	var stat *LogStat
	fmt.Println(stat.Header())

	// real counts
	cache := make([]string, 0, 16)
	var recur func(node *Node)
	recur = func(node *Node) {
		// self
		cache = append(cache, node.Part)
		if len(node.Counts) > 0 {
			for key, count := range node.Counts {
				fullPath := strings.Join(cache[1:], "/") // skip the dummy one
				log := ReconstructLog(key, fullPath)
				stat = NewLogStat(log, count)
				fmt.Println(stat.String())
			}
		}
		// children from left to right
		for _, child := range node.Children {
			recur(child)
		}
		// pop the last part
		cache = cache[:len(cache)-1]
	}
	recur(root)
}
