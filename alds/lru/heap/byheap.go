package lru

import "container/heap"

type LRUCache struct {
    cap int
    kv map[int]*Node
    lru *LRUHeap
    weight int
}

type Node struct{
    Key int
    Weight int
    Val int
    Index int
}

func NewNode(key int, val int, weight int)*Node{
    return &Node{
        Key: key,
        Weight: weight,
        Val: val,
        Index: -1,
    }
}

type LRUHeap []*Node
func (h LRUHeap) Len() int           { return len(h) }
func (h LRUHeap) Less(i, j int) bool { return h[i].Weight < h[j].Weight } //this is big heap or?
func (h LRUHeap) Swap(i, j int){ 
    h[i], h[j] = h[j], h[i]
    h[i].Index, h[j].Index = h[j].Index, h[i].Index
}
func (h *LRUHeap) Push(x any) {
    idx := len(*h)
    x.(*Node).Index = idx
	*h = append(*h, x.(*Node))
}
func (h *LRUHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

    old[n-1] = nil // for gc
    x.Index = -1 // for safety
	return x
}

func Constructor(capacity int) LRUCache {
    lru := LRUHeap{}
    heap.Init(&lru)
    return LRUCache{
        cap: capacity,
        kv: make(map[int]*Node, 64),
        lru: &lru,
    }
}

func (this *LRUCache) evict() {
    node := heap.Pop(this.lru).(*Node)
    delete(this.kv, node.Key)
}

func (this *LRUCache) getWeight() int{
    this.weight+=1
    return this.weight
}

func (this *LRUCache) Get(key int) int {
    node := this.get(key)
    if node == nil{
        return -1
    }
    return node.Val
}

func (this *LRUCache) get(key int) *Node{
    node, ok := this.kv[key]
    // fmt.Println("get kv", key, ok)
    if ok {
        // update LRU
        node.Weight = this.getWeight()
        heap.Fix(this.lru, node.Index)
        return node
    }
    return nil
}


func (this *LRUCache) Put(key int, value int)  {
    node := this.get(key)
    if node != nil {
        node.Val = value
        return
    }
    // assume key exists for now
    if len(this.kv) == this.cap {
        // fmt.Println("key triggers evict", key, value)
        this.evict()
    }
    node = NewNode(key, value, this.getWeight())
    heap.Push(this.lru, node)
    this.kv[key] = node
    // fmt.Println("put kv", key, value)
}


/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */