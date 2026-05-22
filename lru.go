package main

type LRU struct {
	dll         DoubleLinkedList
	items       map[string]*Node
	cacheSize   uint32
	currentSize uint32
}

func NewLRU(c uint32) *LRU {
	return &LRU{
		dll:         NewDLL(),
		items:       make(map[string]*Node),
		cacheSize:   c,
		currentSize: 0,
	}
}
func (lru *LRU) Insert(key string, size uint32) {
	n := lru.dll.insert(key, size)
	if lru.currentSize+size > lru.cacheSize {
		_, _ = lru.Evict()
	}
	lru.items[key] = n
	lru.currentSize += size
}
func (lru *LRU) Get(key string) (any, bool) {
	n, exists := lru.items[key]
	if !exists {
		return "", false
	}

	lru.dll.moveToHead(n)
	return n.val, true
}
func (lru *LRU) Remove(key string) (any, bool) {
	n, exists := lru.items[key]
	if !exists {
		return "", false
	}
	lru.dll.remove(n)
	lru.currentSize -= n.size
	delete(lru.items, key)
	return n.val, true
}
func (lru *LRU) Contains(key string) bool {
	_, exists := lru.items[key]
	return exists
}
func (lru *LRU) Evict() (any, any) {
	if lru.currentSize == 0 {
		return "", ""
	}
	toEvict := lru.dll.tail.prev

	lru.dll.remove(toEvict)

	evictKey, ok := toEvict.key.(string)
	if !ok {
		panic("invalid key")
	}
	delete(lru.items, evictKey)
	lru.currentSize -= toEvict.size
	return toEvict.key, toEvict.val
}
func (lru *LRU) Len() uint32 {
	return lru.currentSize
}
