package trie

import "fmt"

type Node struct {
	character  rune
	terminated bool
	children   map[rune]*Node
}

func (node *Node) SearchChild(ch rune) *Node {
	return node.children[ch]
}

func (node *Node) AddChild(ch rune) *Node {
	node.children[ch] = NewNode(ch)
	return node.children[ch]
}

func NewNode(ch rune) *Node {
	return &Node{
		character:  ch,
		terminated: false,
		children:   make(map[rune]*Node),
	}
}

// methods over a root node

func Insert(root *Node, word string) {
	node := root
	for _, char := range word {
		child := node.SearchChild(char)
		if child == nil {
			child = node.AddChild(char)
		}
		node = child
	}
	node.terminated = true
}

// harder:
// if the word does not exist, return false
func Delete(root *Node, word string) bool {
	var helper func(node *Node, depth int) (*Node, bool)
	helper = func(node *Node, depth int) (*Node, bool) {
		// edge case: tree does not contain the word
		if node == nil {
			return nil, false
		}

		// edge case: the last character of the word
		if depth == len(word) {
			if node.terminated {
				node.terminated = false
				if len(node.children) == 0 {
					return nil, true
				}
				return node, true
			}
			return node, false // the word is not in the tree
		}

		// normal case
		ch := rune(word[depth])
		child := node.SearchChild(ch)
		if child == nil {
			return nil, false
		}
		child, found := helper(child, depth+1)
		if found && child == nil {
			delete(node.children, ch) // how to delete a key from a map?
		}
		return node, found
	}
	_, found := helper(root, 0)
	return found
}

func Search(root *Node, word string) bool {
	node := root
	for _, ru := range word {
		child := node.SearchChild(ru)
		if child == nil {
			return false
		}
		node = node.SearchChild(ru)
	}
	// should check the terminated flag
	return node.terminated
}

// medium: pre-order traversal in recursive way, stack cache;
func Print(root *Node) {
	fmt.Println("Printing Trie, the root node is:", string(root.character))
	cache := make([]rune, 0, 16)
	var recur func(node *Node)
	recur = func(node *Node) {
		cache = append(cache, node.character)
		if node.terminated {
			word := string(cache[1:])
			fmt.Println("WORD:", word)
		}
		for _, child := range node.children {
			recur(child)
		}
		cache = cache[:len(cache)-1]
	}
	recur(root)
}
