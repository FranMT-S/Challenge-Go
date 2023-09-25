package Helpers

//

import (
	"sync"
)

type queuenode struct {
	data string
	next *queuenode
}

// A go-routine safe FIFO (first in first out) data stucture safe to concurrent.
type QueueSafe struct {
	head  *queuenode
	tail  *queuenode
	count int
	lock  *sync.Mutex
}

// returns a new pointer to a new queue safe to concurrent.
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

// structure basic of queue
type QueueBasic struct {
	head  *queuenode
	tail  *queuenode
	count int
}

// return a structure basic of queue
func NewQueueBasic() *QueueBasic {
	q := &QueueBasic{}
	return q
}

// return length of the queue
func (q *QueueBasic) Len() int {
	return q.count
}

// insert a element in the queue
func (q *QueueBasic) Push(item string) {
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

// returns and extracts the last element from the head of the queue
func (q *QueueBasic) Poll() string {
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

// returns witout extracts the last element from the head of the queue
func (q *QueueBasic) Peek() string {
	n := q.head
	if n == nil {
		return ""
	}

	return n.data
}
