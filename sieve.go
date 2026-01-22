package main

import "fmt"

var (
	ErrInvalidSize error = fmt.Errorf("item can't fit")
	ErrNotFound    error = fmt.Errorf("item not found")
)

type EvictFn func(*Sieve) int

type SieveOpts struct {
	Evict EvictFn
}

type Sieve struct {
	SieveOpts
	dll      DoubleLinkedList
	capacity int
	size     int
	items    map[string]*Node

	hand *Node

	hits       int
	misses     int
	byteHits   int
	byteMisses int
}

func NewSieve(cap int, opts SieveOpts) Sieve {
	// if opts.Evict == nil {
	// 	opts.Evict =
	// }
	s := Sieve{
		dll:       NewDLL(),
		capacity:  cap,
		items:     make(map[string]*Node, cap),
		SieveOpts: opts,
	}
	s.hand = s.dll.tail
	return s
}

func (s *Sieve) GetOrInsert(k string) (any, error) {
	if len(k) > s.capacity {
		return nil, ErrInvalidSize
	}
	node, exists := s.items[k]
	// cache hit
	if exists {
		s.hits++
		s.byteHits += len(k)
		node.visited = true
		node.credit += node.penalty
		node.penalty /= 2
		s.dll.moveToHead(node)
		return node.val, nil
	}
	s.misses++
	s.byteMisses += len(k)

	return nil, s.Put(k, "")
}
func (s *Sieve) Get(k string) (any, error) {
	if len(k) > s.capacity {
		return nil, ErrInvalidSize
	}
	node, exists := s.items[k]
	// cache hit
	if exists {
		s.hits++
		s.byteHits += len(k)
		node.visited = true
		s.dll.moveToHead(node)
		return node.val, nil
	}
	s.misses++
	s.byteMisses += len(k)
	return nil, ErrNotFound
}


func (s *Sieve) Put(k, v string) error {
	// if cache is full
	for s.size >= s.capacity {
		s.SieveOpts.Evict(s)
	}
	// insert
	newNode := s.dll.insert(k)
	s.items[k] = newNode
	s.size += len(k)
	return nil
}
