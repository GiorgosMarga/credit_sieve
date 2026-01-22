package main

type Node struct {
	next    *Node
	prev    *Node
	visited bool
	credit  int
	penalty int
	val     string
	key     string
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
func (dll *DoubleLinkedList) insert(val string) *Node {
	n := &Node{
		visited: false,
		val:     val,
		key:     val,
		credit:  len(val),
		penalty: 1,
	}

	dll.head.next.prev = n
	n.next = dll.head.next
	n.prev = dll.head
	dll.head.next = n
	dll.items++
	return n
}
