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
	for i := 1; i <= maxInternalSliceSize; i++ {
		pushValue++
		s.Push(pushValue)
	}

	// Push 1 extra item to force the creation of a new array
	pushValue++
	s.Push(pushValue)
	extraAddedItems++
	checkLinks(t, s, pushValue, maxInternalSliceSize)

	// Push another maxInternalSliceSize-1 to fill the second array
	for i := 1; i <= maxInternalSliceSize-1; i++ {
		pushValue++
		s.Push(pushValue)
		checkLinks(t, s, pushValue, maxInternalSliceSize)
	}

	// Push 1 extra item to force the creation of a new array (3 total)
	pushValue++
	s.Push(pushValue)
	checkLinks(t, s, pushValue, maxInternalSliceSize)

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
	checkLinks(t, s, popValue, maxInternalSliceSize)

	// Pop maxInternalSliceSize-1 items to empty the tail (middle) slice
	for i := 1; i <= maxInternalSliceSize-1; i++ {
		if v, ok := s.Pop(); !ok || v.(int) != popValue {
			t.Errorf("Expected: %d; Got: %d", popValue, v)
		}
		popValue--
		checkLinks(t, s, popValue, maxInternalSliceSize)
	}

	// Pop one extra item to force moving the tail to the head (first) slice. This also means the old tail
	// slice should have no items now.
	if v, ok := s.Pop(); !ok || v.(int) != popValue {
		t.Errorf("Expected: %d; Got: %d", popValue, v)
	}
	popValue--
	checkLinks(t, s, popValue, maxInternalSliceSize)

	// Pop maxFirstSliceSize-1 items to empty the head (first) slice
	for i := 1; i <= maxInternalSliceSize; i++ {
		if v, ok := s.Pop(); !ok || v.(int) != popValue {
			t.Errorf("Expected: %d; Got: %d", popValue, v)
		}
		popValue--
		checkLinks(t, s, popValue, maxInternalSliceSize)
	}

	// The stack shoud be empty
	if s.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, s.Len())
	}
	if _, ok := s.Back(); ok {
		t.Error("Expected: false; Got: true")
	}
	if cap(s.tail.v) != maxInternalSliceSize {
		t.Errorf("Expected: %d; Got: %d", maxInternalSliceSize, cap(s.tail.v))
	}
}

// Helper methods-----------------------------------------------------------------------------------

// Checks the internal slices and its links.
func checkLinks(t *testing.T, s *Stack, length, tailSliceSize int) {
	t.Helper()
	if s.Len() != length {
		t.Errorf("Unexpected length; Expected: %d; Got: %d", length, s.Len())
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
	if s.tail == nil {
		// Zero value.
		if s.tail != nil {
			fail("nil tail when zero", s.tail, nil)
		}
		if s.len != 0 {
			fail("zero length when zero", s.len, 0)
		}
		return
	}
	if t.Failed() {
		t.FailNow()
	}
}
