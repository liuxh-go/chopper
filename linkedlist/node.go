package linkedlist

// Node 链表节点
type Node[T any] struct {
	data T
}

// SingleNode 单向链表节点
type SingleNode[T any] struct {
	*Node[T]

	next *SingleNode[T]
}

type DoubleNode[T any] struct {
	*Node[T]

	prev *DoubleNode[T]
	next *DoubleNode[T]
}
