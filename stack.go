// Copyright (c) 2018 ef-ds
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package stack implements a very fast and efficient general purpose
// Last-In-First-Out (LIFO) stack data structure that is specifically
// optimized to perform when used by Microservices and serverless services
// running in production environments.
package stack

const (
	// firstSliceSize holds the size of the first slice.
	firstSliceSize = 4

	// sliceGrowthFactor determines by how much and how fast the first internal
	// slice should grow. A growth factor of 4, firstSliceSize = 4 and
	// maxInternalSliceSize = 256, the first slice will start with size 4,
	// then 16 (4*4), then 64 (16*4), then 256 (64*4), then 1024 (256*4).
	// The growth factor should be tweaked together with firstSliceSize and
	// maxInternalSliceSize and for maximum efficiency.
	// sliceGrowthFactor only applies to the very first slice creates. All other
	// subsequent slices are created with fixed size of maxInternalSliceSize.
	sliceGrowthFactor = 4

	// maxInternalSliceSize holds the maximum size of each internal slice.
	maxInternalSliceSize = 1024
)

// Stack implements an unbounded, dynamically growing Last-In-First-Out (LIFO)
// stack data structure.
// The zero value for stack is an empty stack ready to use.
type Stack struct {
	// Head points to the first node of the linked list.
	head *node

	// Tail points to the last node of the linked list.
	// In an empty stack, head and tail points to the same node.
	tail *node

	// Len holds the current stack values length.
	len int
}

// Node represents a stack node.
// Each node holds a slice of user managed values.
type node struct {
	// v holds the list of user added values in this node.
	v []interface{}

	// n points to the next node in the linked list.
	n *node

	// p points to the previous node in the linked list.
	p *node
}

// New returns an initialized stack.
func New() *Stack {
	return new(Stack)
}

// Init initializes or clears stack s.
func (s *Stack) Init() *Stack {
	*s = Stack{}
	return s
}

// Len returns the number of elements of stack s.
// The complexity is O(1).
func (s *Stack) Len() int { return s.len }

// Back returns the last element of stack d or nil if the stack is empty.
// The second, bool result indicates whether a valid value was returned;
// if the stack is empty, false will be returnes.
// The complexity is O(1).
func (s *Stack) Back() (interface{}, bool) {
	if s.len == 0 {
		return nil, false
	}
	return s.tail.v[len(s.tail.v)-1], true
}

// Push adds value v to the the back of the stack.
// The complexity is O(1).
func (s *Stack) Push(v interface{}) {
	switch {
	case s.head == nil:
		// No nodes present yet.
		h := &node{v: make([]interface{}, 0, firstSliceSize)}
		h.p = h
		s.head = h
		s.tail = h
	case len(s.tail.v) < cap(s.tail.v):
		// There's room in the tail slice.
	case cap(s.tail.v) < maxInternalSliceSize:
		// We're on the first slice and it hasn't grown large enough yet.
		l := len(s.tail.v)
		nv := make([]interface{}, l, l*sliceGrowthFactor)
		copy(nv, s.tail.v)
		s.tail.v = nv
	case s.tail.n != nil:
		// There's at least one unused slice between head and tail nodes.
		n := s.tail.n
		s.tail = n
	default:
		// No available nodes, so make one.
		n := &node{v: make([]interface{}, 0, maxInternalSliceSize)}
		n.p = s.tail
		s.tail.n = n
		s.tail = n
	}
	s.len++
	s.tail.v = append(s.tail.v, v)
}

// Pop retrieves and removes the current element from the back of the stack.
// The second, bool result indicates whether a valid value was returned;
// if the stack is empty, false will be returnes.
// The complexity is O(1).
func (s *Stack) Pop() (interface{}, bool) {
	if s.len == 0 {
		return nil, false
	}
	s.len--
	tp := len(s.tail.v) - 1
	vp := &s.tail.v[tp]
	v := *vp
	*vp = nil // Avoid memory leaks
	s.tail.v = s.tail.v[:tp]
	switch {
	case tp > 0:
		// There's space before tp.
	default:
		// Leave the slice unused as spare.
		s.tail = s.tail.p
	}
	return v, true
}
