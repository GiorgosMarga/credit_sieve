package main

import "math"

type GMR struct {
	dll      DoubleLinkedList
	capacity uint32
	size     uint32
	items    map[string]*Node
	k        int

	hand *Node

	hits       int
	misses     int
	byteHits   uint32
	byteMisses uint32
}

func NewGMR(cap uint32) *GMR {
	s := &GMR{
		dll:      NewDLL(),
		capacity: cap,
		items:    make(map[string]*Node),
		k:        2,
	}
	s.hand = s.dll.tail
	return s
}

func (s *GMR) GetOrInsert(k string, size uint32) (any, error) {
	if size > s.capacity {
		return nil, ErrInvalidSize
	}
	node, exists := s.items[k]
	// cache hit
	if exists {
		s.hits++
		s.byteHits += node.size

		node.credit += int(math.Log2(float64(node.size)))
		// node.penalty = max(1, node.penalty>>1)
		s.dll.moveToHead(node)

		return node.val, nil
	}
	s.misses++
	s.byteMisses += size

	return nil, s.Put(k, size)
}

func (s *GMR) Put(k string, size uint32) error {
	if size > s.capacity {
		return ErrInvalidSize
	}
	// if cache is full
	for s.size+size > s.capacity {
		s.EvictByGravity()
	}

	// insert
	newNode := s.dll.insert(k, size)
	s.items[k] = newNode
	s.size += size
	return nil
}

func (s *GMR) EvictByGravity() uint32 {
	if len(s.items) == 0 {
		return 0
	}
	o := s.hand
	//tail is a dummy node
	if o == s.dll.tail || o == s.dll.head {
		o = s.dll.tail.prev
	}

	for o.credit > 0 {
		o.credit -= o.penalty
		o.penalty <<= 1

		o = o.prev
		// head is a dummy node
		// wrap around list
		if o == s.dll.head {
			o = s.dll.tail.prev
		}
	}
	s.hand = o.prev
	s.dll.remove(o) // eviction
	s.size -= o.size

	deleteKey, ok := o.key.(string)
	if !ok {
		panic("invalid key on gmr cache")
	}
	delete(s.items, deleteKey)
	return o.size
}

func (gmr *GMR) GetName() string {
	return "GMR"
}

func (gmr *GMR) GetHits() uint32 {
	return uint32(gmr.hits)
}

func (gmr *GMR) GetByteHits() uint32 {
	return gmr.byteHits
}
