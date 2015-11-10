package spmc

import "sync/atomic"
import "unsafe"
import "runtime"

type Node struct {
	next unsafe.Pointer
	v    interface{}
}

var dummy = unsafe.Pointer(&Node{})

type SPMCQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func New() *SPMCQueue {
	return &SPMCQueue{head: dummy, tail: dummy}
}

func (q *SPMCQueue) Push(v interface{}) {
	n := unsafe.Pointer(&Node{v: v})
repeat:
	prev := atomic.LoadPointer(&q.head)
	if !atomic.CompareAndSwapPointer(&q.head, prev, n) {
		goto repeat
	}

	(*Node)(prev).next = n

}

func (q *SPMCQueue) Pop() interface{} {
repeat:
	tail := (*Node)(q.tail)
	if tail.next != nil {
		q.tail = tail.next
		return (*Node)(q.tail).v
	}

	runtime.Gosched()
	goto repeat

	return nil
}
