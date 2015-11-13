package mpsc

import "sync/atomic"
import "unsafe"
import "runtime"

type Node struct {
	next unsafe.Pointer
	v    interface{}
}

var dummy = unsafe.Pointer(&Node{})

type MPSCQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	size int64
}

func New() *MPSCQueue {
	return &MPSCQueue{head: dummy, tail: dummy}
}

func (q *MPSCQueue) Push(v interface{}) {
	n := unsafe.Pointer(&Node{v: v})
repeat:
	prev := atomic.LoadPointer(&q.head)
	if !atomic.CompareAndSwapPointer(&q.head, prev, n) {
		goto repeat
	}

	(*Node)(prev).next = n
	atomic.AddInt64(&q.size, 1)

}

func (q *MPSCQueue) Pop() interface{} {
repeat:
	tail := (*Node)(q.tail)
	if tail.next != nil {
		q.tail = tail.next
		atomic.AddInt64(&q.size, -1)
		return (*Node)(q.tail).v
	}

	runtime.Gosched()
	goto repeat

	return nil
}

func (q *MPSCQueue) Size() int64 {
	return atomic.LoadInt64(&q.size)
}
