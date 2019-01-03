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

package stack_test

import (
	"testing"

	"github.com/ef-ds/stack"
)

func TestPopWithZeroValueShouldReturnReadyToUsestack(t *testing.T) {
	var s stack.Stack
	s.Push(1)
	s.Push(2)

	v, ok := s.Back()
	if !ok || v.(int) != 2 {
		t.Errorf("Expected: 2; Got: %d", v)
	}
	v, ok = s.Pop()
	if !ok || v.(int) != 2 {
		t.Errorf("Expected: 2; Got: %d", v)
	}
	v, ok = s.Back()
	if !ok || v.(int) != 1 {
		t.Errorf("Expected: 1; Got: %d", v)
	}
	v, ok = s.Pop()
	if !ok || v.(int) != 1 {
		t.Errorf("Expected: 1; Got: %d", v)
	}
	_, ok = s.Back()
	if ok {
		t.Error("Expected: empty slice (ok=false); Got: ok=true")
	}
	_, ok = s.Pop()
	if ok {
		t.Error("Expected: empty slice (ok=false); Got: ok=true")
	}
}

func TestWithZeroValueAndEmptyShouldReturnAsEmpty(t *testing.T) {
	var s stack.Stack

	if _, ok := s.Back(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if _, ok := s.Back(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if _, ok := s.Pop(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if l := s.Len(); l != 0 {
		t.Errorf("Expected: 0 as the queue is empty; Got: %d", l)
	}
}

func TestInitShouldReturnEmptystack(t *testing.T) {
	var s stack.Stack
	s.Push(1)

	s.Init()

	if _, ok := s.Back(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if _, ok := s.Pop(); ok {
		t.Error("Expected: false as the queue is empty; Got: true")
	}
	if l := s.Len(); l != 0 {
		t.Errorf("Expected: 0 as the queue is empty; Got: %d", l)
	}
}

func TestPopWithNilValuesShouldReturnAllValuesInOrder(t *testing.T) {
	s := stack.New()
	s.Push(1)
	s.Push(nil)
	s.Push(2)
	s.Push(nil)

	v, ok := s.Pop()
	if !ok || v != nil {
		t.Errorf("Expected: nil; Got: %d", v)
	}
	v, ok = s.Pop()
	if !ok || v.(int) != 2 {
		t.Errorf("Expected: 2; Got: %d", v)
	}
	v, ok = s.Pop()
	if !ok || v != nil {
		t.Errorf("Expected: nil; Got: %d", v)
	}
	v, ok = s.Pop()
	if !ok || v.(int) != 1 {
		t.Errorf("Expected: 1; Got: %d", v)
	}
	_, ok = s.Pop()
	if ok {
		t.Error("Expected: empty slice (ok=false); Got: ok=true")
	}
}
