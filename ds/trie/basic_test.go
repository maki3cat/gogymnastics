package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie(t *testing.T) {
	root := NewNode('/')
	assert.False(t, Search(root, "hello"))
	Insert(root, "hello")
	assert.True(t, Search(root, "hello"))

	Insert(root, "world")

	assert.False(t, Search(root, "w"))
	assert.False(t, Search(root, "worldabc"))
	assert.False(t, Search(root, "helloworld"))
	assert.True(t, Search(root, "world"))
	assert.False(t, Search(root, "WORLD"))

	// print the trie
	Print(root)
}

func TestTrie_Print(t *testing.T) {
	root := NewNode('/')
	Insert(root, "cat")
	Insert(root, "caterpillar")
	Insert(root, "dog")
	Insert(root, "door")
	Insert(root, "doggy")
	Print(root)
}

func TestTrie_Delete(t *testing.T) {
	root := NewNode('/')
	Insert(root, "cat")
	Insert(root, "caterpillar")
	Insert(root, "door")
	Insert(root, "doggy")
	Print(root)
	assert.True(t, Delete(root, "cat"))
	assert.False(t, Delete(root, "dog"))
	assert.False(t, Delete(root, "cat"))
	Print(root)
}
