package urlrouter

import "strings"

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

func (n *Node) AddChild(name string) *Node {
	n.Children[name] = NewNode(name)
	return n.Children[name]
}

func (n *Node) AddMethod(method string, handler func()) {
	n.Methods[method] = handler
}

// ----------------build the trie----------------
var urlRoot *Node = NewNode("") // dummy root

func Register(path string, method string, handler func()) {
	parts := strings.Split(path, "/")
	node := urlRoot
	for _, part := range parts {
		child := node.FindChild(part)
		if child == nil {
			child = node.AddChild(part)
		}
		node = child
	}
	node.AddMethod(method, handler)
}

func Route(path string, method string) func() {
	node := urlRoot

}
