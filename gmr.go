package main

type GMR struct {
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

func NewGMR(cap uint32) *GMR {
	s := &GMR{
		dll:      NewDLL(),
		capacity: cap,
		items:    make(map[string]*Node),
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
		node.credit += node.penalty
		node.penalty = max(1, node.penalty/2)
		s.dll.moveToHead(node)
		return node.val, nil
	}
	s.misses++
	s.byteMisses += size

	return nil, s.Put(k, size)
}
func (s *GMR) Get(k string, size uint32) (any, error) {
	if size > s.capacity {
		return nil, ErrInvalidSize
	}
	node, exists := s.items[k]
	// cache hit
	if exists {
		s.hits++
		s.byteHits += node.size
		s.dll.moveToHead(node)
		return node.val, nil
	}
	s.misses++
	s.byteMisses += node.size
	return nil, ErrNotFound
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
		o.penalty *= 2
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
