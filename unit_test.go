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

package stack

import (
	"fmt"
	"testing"
)

const (
	refillCount = 3
	pushCount   = maxInternalSliceSize * 3 // Push to fill at least 3 internal slices
)

func TestNewShouldReturnInitiazedInstanceOfstack(t *testing.T) {
	s := New()
	assertInvariants(t, s, nil)
}

func TestInvariantsWhenEmptyInMiddleOfSlice(t *testing.T) {
	s := new(Stack)
	s.Push(0)
	assertInvariants(t, s, nil)
	s.Push(1)
	assertInvariants(t, s, nil)
	s.Pop()
	assertInvariants(t, s, nil)
	s.Pop()
	// At this point, the stack is empty and hp will
	// not be pointing at the start of the slice.
	assertInvariants(t, s, nil)
}

func TestPushPopShouldHaveAllInternalLinksInARing(t *testing.T) {
	s := New()
	pushValue, extraAddedItems := 0, 0

	// Push maxInternalSliceSize items to fill the first array
	expectedHeadSliceSize := firstSliceSize
	for i := 1; i <= maxInternalSliceSize; i++ {
		pushValue++
		s.Push(pushValue)

		if pushValue >= expectedHeadSliceSize {
			expectedHeadSliceSize *= sliceGrowthFactor
		}
	}

	// Push 1 extra item to force the creation of a new array
	pushValue++
	s.Push(pushValue)
	extraAddedItems++
	checkLinks(t, s, pushValue, maxInternalSliceSize, maxInternalSliceSize, s.tail, s.head, nil, s.head)

	// Push another maxInternalSliceSize-1 to fill the second array
	for i := 1; i <= maxInternalSliceSize-1; i++ {
		pushValue++
		s.Push(pushValue)
		checkLinks(t, s, pushValue, maxInternalSliceSize, maxInternalSliceSize, s.tail, s.head, nil, s.head)
	}

	// Push 1 extra item to force the creation of a new array (3 total)
	pushValue++
	s.Push(pushValue)
	checkLinks(t, s, pushValue, maxInternalSliceSize, maxInternalSliceSize, s.tail.p, s.head, nil, s.head.n)
	/// Check middle links
	if s.head.n.n != s.tail {
		t.Error("Expected: s.head.n.n == s.tail; Got: s.head.n.n != s.tail")
	}
	if s.head.n.p != s.head {
		t.Error("Expected: s.head.n.p == s.head; Got: s.head.n.p != s.head")
	}

	// Check final len after all pushes
	if s.Len() != maxInternalSliceSize+maxInternalSliceSize+extraAddedItems {
		t.Errorf("Expected: %d; Got: %d", maxInternalSliceSize+maxInternalSliceSize+extraAddedItems, s.Len())
	}

	// Pop one item to force moving the tail to the middle slice. This also means the old tail
	// slice should have no items now
	popValue := s.Len()
	if v, ok := s.Pop(); !ok || v.(int) != popValue {
		t.Errorf("Expected: %d; Got: %d", popValue, v)
	}
	popValue--
	checkLinks(t, s, popValue, maxInternalSliceSize, maxInternalSliceSize, s.tail, s.head, s.tail.n, s.head)
	//Check last slice links (not tail anymore; tail is the middle one)
	if s.tail.n.n != nil {
		t.Error("Expected: s.tail.n.n == nil; Got: s.tail.n.n != nil")
	}
	if s.tail.n.p != s.tail {
		t.Error("Expected: s.tail.n.p == s.tail; Got: s.tail.n.p != s.tail")
	}

	// Pop maxInternalSliceSize-1 items to empty the tail (middle) slice
	for i := 1; i <= maxInternalSliceSize-1; i++ {
		if v, ok := s.Pop(); !ok || v.(int) != popValue {
			t.Errorf("Expected: %d; Got: %d", popValue, v)
		}
		popValue--
		checkLinks(t, s, popValue, maxInternalSliceSize, maxInternalSliceSize, s.tail, s.head, s.tail.n, s.head)
		/// Check last slice links
		if s.tail.n.n != nil {
			t.Error("Expected: s.tail.n.n == nil; Got: s.tail.n.n != nil")
		}
		if s.tail.n.p != s.tail {
			t.Error("Expected: s.tail.n.p == s.tail; Got: s.tail.n.p != s.tail")
		}
	}

	// Pop one extra item to force moving the tail to the head (first) slice. This also means the old tail
	// slice should have no items now.
	if v, ok := s.Pop(); !ok || v.(int) != popValue {
		t.Errorf("Expected: %d; Got: %d", popValue, v)
	}
	popValue--
	checkLinks(t, s, popValue, maxInternalSliceSize, maxInternalSliceSize, s.tail.n, s.head, s.tail.n, s.head)
	/// Check middle links
	if s.head.n.n != s.tail.n.n {
		t.Error("Expected: s.head.n.n == s.tail.n.n; Got: s.head.n.n != s.tail.n.n")
	}
	if s.head.n.p != s.tail {
		t.Error("Expected: s.head.n.p == s.tail; Got: s.head.n.p != s.tail")
	}
	//Check last slice links (not tail anymore; tail is the first one)
	if s.head.n.n.n != nil {
		t.Error("Expected: s.head.n.n.n == nil; Got: s.head.n.n.n != nil")
	}
	if s.head.n.n.p != s.tail.n {
		t.Error("Expected: s.head.n.n.p == s.tail.n; Got: s.head.n.n.p != s.tail.n")
	}

	// Pop maxFirstSliceSize-1 items to empty the head (first) slice
	for i := 1; i <= maxInternalSliceSize; i++ {
		if v, ok := s.Pop(); !ok || v.(int) != popValue {
			t.Errorf("Expected: %d; Got: %d", popValue, v)
		}
		popValue--
		checkLinks(t, s, popValue, maxInternalSliceSize, maxInternalSliceSize, s.tail.n, s.head, s.tail.n, s.head)
		/// Check middle links
		if s.head.n.n != s.tail.n.n {
			t.Error("Expected: s.head.n.n == s.tail.n.n; Got: s.head.n.n != s.tail.n.n")
		}
		if s.head.n.p != s.tail {
			t.Error("Expected: s.head.n.p == s.tail; Got: s.head.n.p != s.tail")
		}
		//Check last slice links (not tail anymore; tail is the first one)
		if s.head.n.n.n != nil {
			t.Error("Expected: s.head.n.n.n == nil; Got: s.head.n.n.n != nil")
		}
		if s.head.n.n.p != s.tail.n {
			t.Error("Expected: s.head.n.n.p == s.tail.n; Got: s.head.n.n.p != s.tail.n")
		}
	}

	// The stack shoud be empty
	if s.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, s.Len())
	}
	if _, ok := s.Back(); ok {
		t.Error("Expected: false; Got: true")
	}
	if cap(s.head.v) != maxInternalSliceSize {
		t.Errorf("Expected: %d; Got: %d", maxInternalSliceSize, cap(s.head.v))
	}
	if cap(s.tail.v) != maxInternalSliceSize {
		t.Errorf("Expected: %d; Got: %d", maxInternalSliceSize, cap(s.tail.v))
	}
	if s.head.n == s.tail.p {
		t.Error("Expected: s.head.n != s.tail.p; Got: s.head.n == s.tail.p")
	}
}

// Helper methods-----------------------------------------------------------------------------------

// Checks the internal slices and its links.
func checkLinks(t *testing.T, s *Stack, length, headSliceSize, tailSliceSize int, headNext, headPrevious, tailNext, tailPrevious *node) {
	t.Helper()
	if s.Len() != length {
		t.Errorf("Unexpected length; Expected: %d; Got: %d", length, s.Len())
	}
	if cap(s.head.v) != headSliceSize {
		t.Errorf("Unexpected head size; Expected: %d; Got: %d", headSliceSize, len(s.head.v))
	}
	if s.head.n != headNext {
		t.Error("Unexpected head node; Expected: s.head.n == headNext; Got: s.head.n != headNext")
	}
	if s.head.p != headPrevious {
		t.Error("Unexpected head; Expected: s.head.p == headPrevious; Got: s.head.p != headPrevious")
	}
	if s.tail.n != tailNext {
		t.Error("Unexpected tailNext; Expected: s.tail.n == tailNext; Got: s.tail.n != tailNext")
	}
	if s.tail.p != tailPrevious {
		t.Error("Unexpected tailPrevious; Expected: s.tail.p == tailPrevious; Got: s.tail.p != tailPrevious")
	}
	if cap(s.tail.v) != tailSliceSize {
		t.Errorf("Unexpected tail size; Expected: %d; Got: %d", tailSliceSize, len(s.tail.v))
	}
	if t.Failed() {
		t.FailNow()
	}
}

// assertInvariants checks all the invariant conditions in d that we can think of.
// If val is non-nil it is used to find the expected value for an item at index
// i measured from the head of the stack.
func assertInvariants(t *testing.T, s *Stack, val func(i int) interface{}) {
	t.Helper()
	fail := func(what string, got, want interface{}) {
		t.Errorf("invariant fail: %s; got %v want %v", what, got, want)
	}
	if s == nil {
		fail("non-nil stack", s, "non-nil")
	}
	if s.head == nil {
		// Zero value.
		if s.tail != nil {
			fail("nil tail when zero", s.tail, nil)
		}
		if s.Len() != 0 {
			fail("zero length when zero", s.Len(), 0)
		}
		return
	}

	spareLinkCount := 0
	inStack := true
	elemCount := 0
	smallNodeCount := 0
	index := 0
	walkLinks(t, s, func(n *node) {
		if len(n.v) < maxInternalSliceSize {
			smallNodeCount++
			if len(n.v) > maxInternalSliceSize {
				fail("first node within bounds", len(n.v), maxInternalSliceSize)
			}
		}
		if len(n.v) > maxInternalSliceSize {
			fail("slice too big", len(n.v), maxInternalSliceSize)
		}
		for i, v := range n.v {
			failElem := func(what string, got, want interface{}) {
				fail(fmt.Sprintf("at elem %d, node %p, %s", i, n, what), got, want)
				t.FailNow()
			}
			if !inStack {
				if v != nil {
					failElem("all values outside queue nil", v, nil)
				}
				continue
			}
			if v != nil {
				if val != nil {
					want := val(index)
					if want != v {
						failElem(fmt.Sprintf("element %d has expected value", index), v, want)
					}
				}
				elemCount++
				index++
			}
		}
		if !inStack {
			spareLinkCount++
		}
		if n == s.tail {
			inStack = false
		}
	})
	if inStack {
		// We never encountered the tail pointer.
		t.Errorf("tail does not point to element in list")
	}
	if elemCount != s.Len() {
		fail("element count == s.Len()", elemCount, s.Len())
	}
	if smallNodeCount > 1 {
		fail("only one first node", smallNodeCount, 1)
	}
	if t.Failed() {
		t.FailNow()
	}
}

// walkLinks calls f for each node in the linked list.
// It also checks link invariants:
func walkLinks(t *testing.T, s *Stack, f func(n *node)) {
	t.Helper()
	fail := func(what string, got, want interface{}) {
		t.Errorf("link invariant %s fail; got %v want %v", what, got, want)
	}
	n := s.head
	for {
		if n.n != nil && n.n.p != n {
			fail("node.n.p == node", n.n.p, n)
		}
		if n.p.n != nil && n.p.n != n {
			fail("node.p.n == node", n.p.n, n)
		}
		f(n)
		n = n.n
		if n == nil {
			break
		}
	}
}
