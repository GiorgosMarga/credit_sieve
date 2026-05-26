package main

import "fmt"

var (
	ErrInvalidSize error = fmt.Errorf("item can't fit")
	ErrNotFound    error = fmt.Errorf("item not found")
)

type Sieve struct {
	dll      DoubleLinkedList
	capacity uint32
	size     uint32
	items    map[string]*Node

	hand *Node

	hits       int
	misses     int
	byteHits   uint32
	byteMisses uint32
}

func NewSieve(cap uint32) *Sieve {
	s := &Sieve{
		dll:      NewDLL(),
		capacity: cap,
		items:    make(map[string]*Node),
	}
	s.hand = s.dll.tail
	return s
}

func (s *Sieve) GetOrInsert(k string, size uint32) (any, error) {
	if size > s.capacity {
		return nil, ErrInvalidSize
	}
	node, exists := s.items[k]
	// cache hit
	if exists {
		s.hits++
		s.byteHits += node.size

		node.visited = true
		s.dll.moveToHead(node)

		return node.val, nil
	}
	s.misses++
	s.byteMisses += size

	return nil, s.Put(k, size)
}
func (s *Sieve) Get(k string, size uint32) (any, error) {
	node, exists := s.items[k]
	// cache hit
	if exists {
		s.hits++
		s.byteHits += node.size
		node.visited = true
		s.dll.moveToHead(node)
		return node.val, nil
	}
	s.misses++
	s.byteMisses += node.size
	return nil, ErrNotFound
}

func (s *Sieve) Put(k string, size uint32) error {
	if size > s.capacity {
		return ErrInvalidSize
	}
	// if cache is full
	for s.size+size > s.capacity {
		s.Evict()
	}
	// insert
	newNode := s.dll.insert(k, size)
	s.items[k] = newNode
	s.size += size
	return nil
}

func (s *Sieve) Evict() uint32 {
	if len(s.items) == 0 {
		return 0
	}
	o := s.hand
	//tail is a dummy node
	if o == s.dll.tail || o == s.dll.head {
		o = s.dll.tail.prev
	}

	for o.visited == true {
		o.visited = false
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
		panic("invalid key in sieve cache")
	}
	delete(s.items, deleteKey)
	return o.size
}
func (s *Sieve) GetName() string {
	return fmt.Sprintf("SIEVE_%d", s.capacity)
}

func (s *Sieve) GetHits() uint32 {
	return uint32(s.hits)
}

func (s *Sieve) GetByteHits() uint32 {
	return s.byteHits
}
