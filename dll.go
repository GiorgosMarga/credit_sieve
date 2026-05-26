package main

import "math"

type Node struct {
	next    *Node
	prev    *Node
	visited bool
	credit  int
	penalty int
	reward  int
	size    uint32
	val     any
	key     any
}

type DoubleLinkedList struct {
	head  *Node
	tail  *Node
	items int
}

func NewDLL() DoubleLinkedList {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head
	return DoubleLinkedList{
		head: head,
		tail: tail,
	}
}
func (dll *DoubleLinkedList) moveToHead(node *Node) {
	// remove
	node.prev.next = node.next
	node.next.prev = node.prev

	dll.head.next.prev = node
	node.next = dll.head.next
	node.prev = dll.head
	dll.head.next = node
}
func (dll *DoubleLinkedList) remove(node *Node) {
	if node == dll.head || node == dll.tail {
		panic("removing dummy")
	}
	node.prev.next = node.next
	node.next.prev = node.prev
	dll.items--
}
func (dll *DoubleLinkedList) insert(val string, size uint32) *Node {

	n := &Node{
		visited: false,
		val:     val,
		key:     val,
		size:    size,
		credit:  0,
		penalty: 1,
		reward:  int(math.Log2(float64(size))),
	}

	dll.head.next.prev = n
	n.next = dll.head.next
	n.prev = dll.head
	dll.head.next = n
	dll.items++
	return n
}
