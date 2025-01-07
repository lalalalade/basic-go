package types

import "io"

type List interface {
	Add(index int, val any)
	Append(val any)
	Delete(index int)
}

type node struct {
	next *node
}

var _ List = &LinkedList{}

// LinkedList 是一个链表
type LinkedList struct {
	head *node
}

func (l *LinkedList) Add(index int, val any) {
	var p = l.head
	// 找到前驱节点
	for i := 0; i < index && p.next != nil; i++ {
		p = p.next
	}
	if newNode, ok := val.(node); ok {
		newNode.next = p.next
		p.next = &newNode
	}
	panic("implement me")
}

func (l *LinkedList) Append(val any) {
	p := l.head
	for p.next != nil {
		p = p.next
	}
	if newNode, ok := val.(node); ok {
		p.next = &newNode
	}
	panic("implement me")
}

func (l *LinkedList) Delete(index int) {
	p := l.head
	// 找到前驱节点
	for i := 0; i < index && p.next != nil; i++ {
		p = p.next
	}
	p.next = p.next.next
	panic("implement me")
}

func UserList() {
	l1 := LinkedList{}
	l1Ptr := &l1
	var l2 LinkedList = *l1Ptr
	println(l2)

	// 这个是nil
	var l3Ptr *LinkedList
	println(l3Ptr)
}

type ListV1[T any] interface {
	Add(index int, val T)
	Append(val T)
	Delete(index int)
}

type nodeV1[T any] struct {
	val  T
	next *nodeV1[T]
}

type LinkedListV1[T any] struct {
	head *nodeV1[T]
}

func (l *LinkedListV1[T]) Add(index int, val T) {
	newNode := &nodeV1[T]{val: val}
	p := l.head
	for i := 0; i < index && p.next != nil; i++ {
		p = p.next
	}
	newNode.next = p.next
	p.next = newNode
	panic("implement me")
}

func (l *LinkedListV1[T]) Append(val T) {
	p := l.head
	for p.next != nil {
		p = p.next
	}
	p.next = &nodeV1[T]{val: val}
	panic("implement me")
}

func (l *LinkedListV1[T]) Delete(index int) {
	p := l.head
	for i := 0; i < index && p.next != nil; i++ {
		p = p.next
	}
	p.next = p.next.next
	panic("implement me")
}

type nodeV2[T io.Closer] struct {
	val  T
	next *nodeV1[T]
}

func (n nodeV2[T]) Use() {
	n.val.Close()
}

type Integer int

type Number interface {
	~int | uint | uint8
}

func Sum[T Number](vals ...T) T {
	var t T
	for _, elem := range vals {
		t += elem
	}
	return t
}

func UseTypeP() {
	sum1 := Sum[int](2, 3, 4)
	println(sum1)
	sum2 := Sum[uint](2, 3, 4)
	println(sum2)

	list := &LinkedListV1[string]{}
	list.Append("hello")
}
