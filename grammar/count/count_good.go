package count

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var digitRegex = regexp.MustCompile(`^\d+$`)

type LogLine struct {
	Method  string
	RawPath string
	Status  int
}

func (l *LogLine) Key() string {
	return fmt.Sprintf("%s-%d", l.Method, l.Status)
}

type Node struct {
	PathPart string
	Counts   map[string]int
	Children []*Node // should be a pointer than a value or the nodes won't be udpated?
}

func (p *Node) FindInChildren(pathPart string) *Node {
	for _, child := range p.Children {
		if child.PathPart == pathPart {
			return child
		}
	}
	return nil
}

func CountPath(root *Node, logLine LogLine) {
	parts := strings.Split(logLine.RawPath, "/")
	if len(parts) == 0 {
		fmt.Println("corrupt line found in CountPath", logLine.RawPath)
		return
	}
	currentNode := root
	for _, part := range parts {
		if digitRegex.MatchString(part) {
			part = "#"
		}
		child := currentNode.FindInChildren(part)
		// not found, add new path and return
		if child == nil {
			fmt.Println("adding new path", part)
			child = &Node{
				PathPart: part,
				Children: make([]*Node, 0),
				Counts:   make(map[string]int),
			}
			currentNode.Children = append(currentNode.Children, child)
		}
		// if found, update continue
		currentNode = child
	}
	currentNode.Counts[logLine.Key()]++
}

type PathCount struct {
	Path  string
	Count int
	Key   string
}

func CollectPathCount(root *Node) []PathCount {
	names := make([]string, 0)
	pathCounts := make([]PathCount, 0)
	var dfs func(node *Node)
	dfs = func(node *Node) {
		names = append(names, node.PathPart)
		currentPath := strings.Join(names[1:], "/")
		for key, count := range node.Counts {
			if count > 0 {
				pc := PathCount{
					Path:  currentPath,
					Count: count,
					Key:   key,
				}
				pathCounts = append(pathCounts, pc)
			}
		}
		for _, child := range node.Children {
			dfs(child)
		}
		// pop the last name
		if len(names) > 0 {
			names = names[:len(names)-1]
		}
	}
	dfs(root)
	sort.Slice(pathCounts, func(i, j int) bool {
		return pathCounts[i].Count > pathCounts[j].Count
	})
	return pathCounts
}

func ParseLogLine(line string) *LogLine {
	n := len("[02/Nov/2018:21:46:43 +0000] ")
	// split the line wihtout the time: n: end of it
	// separator is " "
	logInfos := strings.Split(line[n:], " ")
	if len(logInfos) < 4 {
		fmt.Println("corrupt line", line)
		return nil
	}
	status, err := strconv.Atoi(logInfos[3])
	if err != nil {
		fmt.Println("corrupt line", line)
		return nil
	}
	res := LogLine{
		Method:  logInfos[0],
		RawPath: logInfos[1],
		Status:  status,
	}
	return &res
}

func Process(logs []string) {
	root := Node{
		PathPart: "/",
		Counts:   make(map[string]int),
		Children: make([]*Node, 0),
	}
	for _, log := range logs {
		logLine := ParseLogLine(log)
		if logLine != nil {
			CountPath(&root, *logLine)
		} else {
			fmt.Println("corrupt line", log)
		}
	}
	fmt.Println("root node", root.Children)
	pathCounts := CollectPathCount(&root)

	for _, pc := range pathCounts {
		fmt.Println(pc.Path, pc.Count, pc.Key)
	}
}
