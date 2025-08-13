package urlrouter

import (
	"regexp"
	"strings"
)

// // url path
// basic case: precise matching
// - /users/{:userID:}/location, GET/DELETE
// - /users, GET
// - /users/{:userID:}/account, GET/PATCH
// step forward: adding * wildcard matching + precise matching has the priority
// - /cities/*
// - /cities/shanghai

// // api
// ```
// func register(path string, method string, handler func())
// func route(path string, method string) handler func()
// ```

// ----------------node of the trie----------------
type Node struct {
	Name     string
	Methods  map[string]func()
	Children map[string]*Node
}

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Methods:  make(map[string]func()),
		Children: make(map[string]*Node),
	}
}

func (n *Node) FindChild(name string) *Node {
	return n.Children[name]
}

func (n *Node) FindWildcardChild() *Node {
	return n.Children["*"]
}

func (n *Node) AddChild(name string) *Node {
	n.Children[name] = NewNode(name)
	return n.Children[name]
}

func (n *Node) AddMethod(method string, handler func()) {
	n.Methods[method] = handler
}

// ----------------build the trie----------------
var urlRoot *Node = NewNode("") // dummy root
var digitValue *regexp.Regexp = regexp.MustCompile(`^\d+$`)
var digitRegister *regexp.Regexp = regexp.MustCompile(`^{(\w+)}$`)
var digitPlaceholder = "#"

func Register(path string, method string, handler func()) {
	parts := strings.Split(path, "/")
	node := urlRoot
	for _, part := range parts {
		if digitRegister.MatchString(part) {
			part = digitPlaceholder
		}
		child := node.FindChild(part)
		if child == nil {
			child = node.AddChild(part)
		}
		node = child
	}
	node.AddMethod(method, handler)
}

func Route(path string, method string) (func(), bool) {
	node := urlRoot
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if digitValue.MatchString(part) {
			part = digitPlaceholder
		}
		child := node.FindChild(part)
		if child == nil {
			// try to find the wildcard child
			wildcardChild := node.FindWildcardChild()
			if wildcardChild == nil {
				return nil, false // not found
			}
			child = wildcardChild
		}
		node = child
	}
	handler, ok := node.Methods[method]
	return handler, ok
}
