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

const (
	pushCount = 1024 * 3 // Push to fill at least 3 internal slices
)

func TestFillStackShouldRetrieveAllElementsInOrder(t *testing.T) {
	var d stack.Stack

	for i := 0; i < pushCount; i++ {
		d.Push(i)
	}
	//fmt.Println(spew.Sdump(d))
	for i := pushCount - 1; i >= 0; i-- {
		if v, ok := d.Pop(); !ok || v.(int) != i {
			t.Errorf("Expected: %d; Got: %d", i, v)
		}
	}
	if d.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, d.Len())
	}
}

func TestRefillStackShouldRetrieveAllElementsInOrder(t *testing.T) {
	var d stack.Stack

	for i := 0; i < refillCount; i++ {
		for j := 0; j < pushCount; j++ {
			d.Push(j)
		}
		for j := pushCount - 1; j >= 0; j-- {
			if v, ok := d.Pop(); !ok || v.(int) != j {
				t.Errorf("Expected: %d; Got: %d", i, v)
			}
		}
		if d.Len() != 0 {
			t.Errorf("Expected: %d; Got: %d", 0, d.Len())
		}
	}
}

func TestRefillFullStackShouldRetrieveAllElementsInOrder(t *testing.T) {
	var d stack.Stack
	for i := 0; i < pushCount; i++ {
		d.Push(i)
	}

	for i := 0; i < refillCount; i++ {
		for j := 0; j < pushCount; j++ {
			d.Push(j)
		}
		for j := pushCount - 1; j >= 0; j-- {
			if v, ok := d.Pop(); !ok || v.(int) != j {
				t.Errorf("Expected: %d; Got: %d", j, v)
			}
		}
		if d.Len() != pushCount {
			t.Errorf("Expected: %d; Got: %d", pushCount, d.Len())
		}
	}
}

func TestSlowIncreaseStackShouldRetrieveAllElementsInOrder(t *testing.T) {
	var d stack.Stack

	count := 0
	for i := 0; i < pushCount; i++ {
		count++
		d.Push(count)
		count++
		d.Push(count)
		if v, ok := d.Pop(); !ok || v.(int) != count {
			t.Errorf("Expected: %d; Got: %d", count, v)
		}
	}
	if d.Len() != pushCount {
		t.Errorf("Expected: %d; Got: %d", pushCount, d.Len())
	}
}

func TestSlowDecreaseStackShouldRetrieveAllElementsInOrder(t *testing.T) {
	var d stack.Stack
	push := 0
	for i := 0; i < pushCount; i++ {
		d.Push(push)
		push++
	}

	count := push
	for i := 0; i < pushCount-1; i++ {
		count--
		if v, ok := d.Pop(); !ok || v.(int) != count {
			t.Errorf("Expected: %d; Got: %d", count, v)
		}
		count--
		if v, ok := d.Pop(); !ok || v.(int) != count {
			t.Errorf("Expected: %d; Got: %d", count, v)
		}

		d.Push(count)
		count++
	}
	count--
	if v, ok := d.Pop(); !ok || v.(int) != count {
		t.Errorf("Expected: %d; Got: %d", count, v)
	}
	if d.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, d.Len())
	}
}

func TestStableStackShouldRetrieveAllElementsInOrder(t *testing.T) {
	var d stack.Stack

	for i := 0; i < pushCount; i++ {
		d.Push(i)
		if v, ok := d.Pop(); !ok || v.(int) != i {
			t.Errorf("Expected: %d; Got: %d", i, v)
		}
	}
	if d.Len() != 0 {
		t.Errorf("Expected: %d; Got: %d", 0, d.Len())
	}
}

func TestStableFullStackShouldRetrieveAllElementsInOrder(t *testing.T) {
	var d stack.Stack
	for i := 0; i < pushCount; i++ {
		d.Push(i)
	}

	count := 0
	for i := 0; i < pushCount; i++ {
		d.Push(i)
		if v, ok := d.Pop(); !ok || v.(int) != count {
			t.Errorf("Expected: %d; Got: %d", count, v)
		}
		count++
	}
	if d.Len() != pushCount {
		t.Errorf("Expected: %d; Got: %d", pushCount, d.Len())
	}
}

func TestPushFrontPopRefillWith0ToPushCountItemsShouldReturnAllValuesInOrder(t *testing.T) {
	var d stack.Stack

	for i := 0; i < refillCount; i++ {
		for k := 0; k < pushCount; k++ {
			for j := 0; j < k; j++ {
				d.Push(j)
			}
			for j := k; j > 0; j-- {
				v, ok := d.Pop()
				if !ok || v == nil || v.(int) != j-1 {
					t.Errorf("Expected: %d; Got: %d", j-1, v)
				}
			}
			if d.Len() != 0 {
				t.Errorf("Expected: %d; Got: %d", 0, d.Len())
			}
		}
	}
}
