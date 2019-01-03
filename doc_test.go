package stack_test

import (
	"fmt"

	"github.com/ef-ds/stack"
)

func Example() {
	var s stack.Stack

	for i := 1; i <= 5; i++ {
		s.Push(i)
	}
	for s.Len() > 0 {
		v, _ := s.Pop()
		fmt.Print(v)
	}
	// Output: 54321
}
