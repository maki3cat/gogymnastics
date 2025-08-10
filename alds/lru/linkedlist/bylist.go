package lru

import "container/list"

type LRUCache struct {
	cap      int
	kv       map[int]*list.Element
	lru      *list.List
	nilValue int
}

type Pair struct {
	key int // we should store this in the list.Element, or we don't know which key to delete in the kv
	val int
}

func Constructor(capacity int) LRUCache {
	lru := list.New()
	kv := make(map[int]*list.Element)
	return LRUCache{
		cap:      capacity,
		kv:       kv,
		lru:      lru,
		nilValue: -1,
	}
}

func (this *LRUCache) Put(key int, value int) {
	// exists
	pair := this.get(key)
	if pair != nil {
		pair.val = value
		return
	}
	pair = &Pair{
		key: key,
		val: value,
	}
	// evict
	if len(this.kv) == this.cap {
		first := this.lru.Front()
		pair := first.Value.(*Pair)
		this.lru.Remove(first)
		delete(this.kv, pair.key)
	}
	// add to back
	e := this.lru.PushBack(pair)
	this.kv[key] = e
}

func (this *LRUCache) Get(key int) int {
	pair := this.get(key)
	if pair == nil {
		return this.nilValue
	}
	return pair.val
}

func (this *LRUCache) get(key int) *Pair {
	e, ok := this.kv[key]
	if !ok {
		return nil
	}
	pair := e.Value.(*Pair)
	this.lru.MoveToBack(e) // back is the newest
	return pair
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
