package Helpers

//

import (
	"sync"
)

type queuenode struct {
	data string
	next *queuenode
}

// A go-routine safe FIFO (first in first out) data stucture.
type QueueSafe struct {
	head  *queuenode
	tail  *queuenode
	count int
	lock  *sync.Mutex
}

// Creates a new pointer to a new queue.
func NewQueueSafe() *QueueSafe {
	q := &QueueSafe{}
	q.lock = &sync.Mutex{}
	return q
}

// Returns the number of elements in the queue (i.e. size/length)
// go-routine safe.
func (q *QueueSafe) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

// Pushes/inserts a value at the end/tail of the queue.
// Note: this function does mutate the queue.
// go-routine safe.
func (q *QueueSafe) Push(item string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := &queuenode{data: item}

	if q.tail == nil {
		q.tail = n
		q.head = n
	} else {
		q.tail.next = n
		q.tail = n
	}
	q.count++
}

// Returns the value at the front of the queue.
// i.e. the oldest value in the queue.
// Note: this function does mutate the queue.
// go-routine safe.
func (q *QueueSafe) Poll() string {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return ""
	}

	n := q.head
	q.head = n.next

	if q.head == nil {
		q.tail = nil
	}
	q.count--

	return n.data
}

// Returns a read value at the front of the queue.
// i.e. the oldest value in the queue.
// Note: this function does NOT mutate the queue.
// go-routine safe.
func (q *QueueSafe) Peek() string {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := q.head
	if n == nil {
		return ""
	}

	return n.data
}

/*
 End Quote Safe
*/

type QuoteBasic struct {
	head  *queuenode
	tail  *queuenode
	count int
}

func NewQueueBasic() *QuoteBasic {
	q := &QuoteBasic{}
	return q
}

func (q *QuoteBasic) Len() int {
	return q.count
}

func (q *QuoteBasic) Push(item string) {
	n := &queuenode{data: item}

	if q.tail == nil {
		q.tail = n
		q.head = n
	} else {
		q.tail.next = n
		q.tail = n
	}
	q.count++
}

func (q *QuoteBasic) Poll() string {
	if q.head == nil {
		return ""
	}

	n := q.head
	q.head = n.next

	if q.head == nil {
		q.tail = nil
	}
	q.count--

	return n.data
}

func (q *QuoteBasic) Peek() string {
	n := q.head
	if n == nil {
		return ""
	}

	return n.data
}
