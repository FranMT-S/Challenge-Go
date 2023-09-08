package Helpers

type Node struct {
	value int
	next  *Node
}

type StackLinkedList struct {
	head *Node
	size int
}

// Size() function
func (s *StackLinkedList) Size() int {
	return s.size
}

// IsEmpty() function
func (s *StackLinkedList) IsEmpty() bool {
	return s.size == 0
}

// Peek() function
func (s *StackLinkedList) Peek() int {
	if s.IsEmpty() {
		return 0
	}
	return s.head.value
}

// Push() function
func (s *StackLinkedList) Push(value int) {
	s.head = &Node{value, s.head}
	s.size++
}

// Pop() function
func (s *StackLinkedList) Pop() (int, bool) {
	if s.IsEmpty() {

		return 0, false
	}
	value := s.head.value
	s.head = s.head.next
	s.size--
	return value, true
}
