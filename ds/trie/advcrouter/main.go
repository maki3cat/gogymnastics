package main

import (
	"fmt"
	"regexp"
	"strings"
)

// feature1: register normal url, handler
// feature2: register url with parameters of pure digits
// feature3: register url with wildcard but has low priority

type Node struct {
	part     string
	children map[string]*Node
	handler  *Handler // todo: should this be pointer or the function?
}

func NewNode(part string) *Node {
	n := new(Node)
	n.part = part
	n.children = make(map[string]*Node)
	return n
}

func (n *Node) GetChild(part string) *Node {
	return n.children[part]
}

// return: the node added
func (n *Node) AddChild(part string) *Node {
	c := NewNode(part)
	n.children[part] = c
	return c
}

var matchDigitPatten = regexp.MustCompile(`^{\w+}$`)
var isDigigt = regexp.MustCompile(`^\d+$`)

const digitPlaceHolder = "#"

type Handler func(string) string

func NewRouter() *Router {
	var dummyRoot *Node = NewNode("/")
	r := new(Router)
	r.dummyRoot = dummyRoot
	return r
}

type Router struct {
	dummyRoot *Node
}

// todo: feature2: accomodate wildcard
// todo: save the error handling for later
func (r *Router) RegisterFunc(path string, handler Handler) {
	path = strings.TrimSpace(path)
	parts := strings.Split(path, "/")
	node := r.dummyRoot
	for _, part := range parts {
		// maki: feature1a: accommodate formal parameter
		if matchDigitPatten.MatchString(part) {
			part = digitPlaceHolder
		}
		child := node.GetChild(part)
		if child == nil {
			child = node.AddChild(part)
		}
		node = child
	}
	node.handler = &handler
}

func (r *Router) findHandler(path string) (param string, handler *Handler) {
	path = strings.TrimSpace(path)
	parts := strings.Split(path, "/")
	node := r.dummyRoot
	param = ""
	handler = nil
	for _, part := range parts {
		// maki: feature1b: extract the actual parameter
		if isDigigt.MatchString(part) {
			param = part
			part = digitPlaceHolder
		}
		child := node.GetChild(part)
		if child == nil {
			// not found
			return param, nil
		}
		node = child
	}
	return param, node.handler // may still be none
}

func (r *Router) HandleRequest(path string) {
	param, handler := r.findHandler(path)
	if handler == nil {
		fmt.Println("path not found")
		return
	}
	(*handler)(param)
}
