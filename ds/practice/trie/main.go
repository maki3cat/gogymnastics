package main

import (
	"fmt"
	"strings"
)

type Node struct {
	ch       rune
	children map[rune]*Node
	terminal bool
}

func NewNode(ch rune) *Node {
	node := new(Node)
	node.ch = ch
	node.children = make(map[rune]*Node)
	return node
}

func (n *Node) GetChild(ch rune) *Node {
	return n.children[ch]
}

func (n *Node) AddChild(ch rune) *Node {
	child := NewNode(ch)
	n.children[ch] = child
	return child
}

// todo: error checking to be handled later, focus on main logic now
func AddWord(root *Node, word string) {
	node := root
	for _, ch := range word {
		child := node.GetChild(ch)
		if child == nil {
			child = node.AddChild(ch)
		}
		node = child
	}
	node.terminal = true
}

func SearchWord(root *Node, word string) bool {
	node := root
	for _, ch := range word {
		child := node.GetChild(ch)
		if child == nil {
			return false
		}
		node = child
	}
	return node.terminal
}

// root: it is a dummy node, we find from chilren instead of the current node
// returns: true = found and deleted successfully; false = not found;
// DFS, with the trick that aside from the return value
// there is a global value to be updated which is the final answer
func DeleteWord(root *Node, word string) bool {
	// idx should be found in the children
	// bool returns true is we should delete the child
	found := false
	var recur func(node *Node, idx int) bool

	// returns: the node should be deleted
	recur = func(node *Node, idx int) bool {
		// EDGY NODEs, base conditions
		if node == nil {
			return false
		}
		// the node containing the last ch
		if idx == len(word) {
			found = true
			node.terminal = false
			// Logic-A
			return len(node.children) == 0
		}

		// MIDDLE POINTs, repetition conditions
		ch := rune(word[idx])
		child := node.GetChild(ch)
		if child == nil {
			// stop the recursion, no node to delete
			return false
		}
		deleted := recur(child, idx+1)
		if deleted {
			delete(node.children, ch)
			// maki: Logic-B:
			// this is very important that if this node itself is a terminal, it cannot be deleted!!
			// the middle points are different a bit!!
			return len(node.children) == 0 && !node.terminal
		} else {
			return false
		}
	}
	recur(root, 0)
	return found
}

func PrintNode(dummyRoot *Node) {
	cache := make([]rune, 0)
	var recur func(node *Node)
	recur = func(node *Node) {
		if node == nil {
			return
		}
		cache = append(cache, node.ch)
		if node.terminal {
			fmt.Println("word:", string(cache[1:]))
		}
		for _, child := range node.children {
			recur(child)
		}
		cache = cache[:len(cache)-1]
	}
	recur(dummyRoot)
}

func main() {
	words := []string{"cat", "cattle", "dog", "bird"}
	root := NewNode('/')
	for _, word := range words {
		fmt.Println("adding word: ", word)
		AddWord(root, word)
	}
	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("starting print the tree")
	PrintNode(root)
	fmt.Println(strings.Repeat("=", 30))
	for _, word := range words {
		found := SearchWord(root, word)
		fmt.Println("found word: ", found)
	}
	fmt.Println(strings.Repeat("=", 30))
	notExistWords := []string{"banana", "doggy", "ca"}
	for _, word := range notExistWords {
		found := SearchWord(root, word)
		fmt.Println("found word: ", found)
	}
	fmt.Println(strings.Repeat("=", 30))
	word := "cattle"
	found := SearchWord(root, word)
	fmt.Println("found word cattle: ", found)
	DeleteWord(root, word)
	found = SearchWord(root, word)
	fmt.Println("after deletion, found word cattle: ", found)
	found = SearchWord(root, "cat")
	fmt.Println("after deletion, found word cat: ", found)
}
