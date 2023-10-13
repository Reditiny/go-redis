package database

type LinkNode struct {
	prev  *LinkNode
	next  *LinkNode
	value string
}

type List struct {
	head *LinkNode
	tail *LinkNode
	len  int
}

// ListConsumer 遍历函数
type ListConsumer func(node *LinkNode) bool

// MakeList makes a new list
func makeList() *List {
	return &List{}
}

func (list *List) rPush(args [][]byte) int {
	for _, arg := range args {
		node := &LinkNode{value: string(arg)}
		list.pushBack(node)
	}
	return list.len
}

// pushFront push node to front
func (list *List) pushFront(node *LinkNode) int {
	if list.len == 0 {
		list.head = node
		list.tail = node
	} else {
		node.next = list.head
		list.head.prev = node
		list.head = node
	}
	list.len++
	return list.len
}

// pushBack push node to back
func (list *List) pushBack(node *LinkNode) int {
	if list.len == 0 {
		list.head = node
		list.tail = node
	} else {
		node.prev = list.tail
		list.tail.next = node
		list.tail = node
	}
	list.len++
	return list.len
}

// popFront pop node from front
func (list *List) popFront() *LinkNode {
	if list.len == 0 {
		return nil
	}
	node := list.head
	list.head = list.head.next
	if list.head != nil {
		list.head.prev = nil
	}
	list.len--
	return node
}

// popBack pop node from back
func (list *List) popBack() *LinkNode {
	if list.len == 0 {
		return nil
	}
	node := list.tail
	list.tail = list.tail.prev
	if list.tail != nil {
		list.tail.next = nil
	}
	list.len--
	return node
}

// get node at index
func (list *List) get(index int) *LinkNode {
	if index < 0 || index >= list.len {
		return nil
	}
	node := list.head
	for i := 0; i < index; i++ {
		node = node.next
	}
	return node
}

// forEach traverse list
func (list *List) forEach(consumer ListConsumer) {
	node := list.head
	for node != nil {
		if !consumer(node) {
			break
		}
		node = node.next
	}
}
