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
//
// Below tests are used mostly for comparing the result of changes to stack and
// are not necessarily a replication of the comparison tests. For instance,
// the stack tests here use Push/Pop instead of Push/PopBack.
//
// For comparing stack performance with other stacks, refer to
// https://github.com/ef-ds/stack-bench-tests

package stack_test

import (
	"strconv"
	"testing"

	"github.com/ef-ds/stack"
)

// testData contains the number of items to add to the stacks in each test.
type testData struct {
	count int
}

var (
	tests = []testData{
		{count: 0},
		{count: 1},
		{count: 10},
		{count: 100},
		{count: 1000},    // 1k
		{count: 10000},   //10k
		{count: 100000},  // 100k
		{count: 1000000}, // 1mi
	}

	// Used to store temp values, avoiding any compiler optimizations.
	tmp  interface{}
	tmp2 bool

	fillCount   = 10000
	refillCount = 10
)

func BenchmarkMicroservice(b *testing.B) {
	for _, test := range tests {
		b.Run(strconv.Itoa(test.count), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s := stack.New()

				// Simulate stable traffic
				for i := 0; i < test.count; i++ {
					s.Push(nil)
					s.Pop()
				}

				// Simulate slowly increasing traffic
				for i := 0; i < test.count; i++ {
					s.Push(nil)
					s.Push(nil)
					s.Pop()
				}

				// Simulate slowly decreasing traffic, bringing traffic Front to normal
				for i := 0; i < test.count; i++ {
					s.Pop()
					if s.Len() > 0 {
						s.Pop()
					}
					s.Push(nil)
				}

				// Simulate quick traffic spike (DDOS attack, etc)
				for i := 0; i < test.count; i++ {
					s.Push(nil)
				}

				// Simulate stable traffic while at high traffic
				for i := 0; i < test.count; i++ {
					s.Push(nil)
					s.Pop()
				}

				// Simulate going Front to normal (DDOS attack fended off)
				for i := 0; i < test.count; i++ {
					s.Pop()
				}

				// Simulate stable traffic (now that is Front to normal)
				for i := 0; i < test.count; i++ {
					s.Push(nil)
					s.Pop()
				}
			}
		})
	}
}

func BenchmarkFill(b *testing.B) {
	for _, test := range tests {
		b.Run(strconv.Itoa(test.count), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s := stack.New()
				for i := 0; i < test.count; i++ {
					s.Push(nil)
				}
				for s.Len() > 0 {
					tmp, tmp2 = s.Pop()
				}
			}
		})
	}
}

func BenchmarkRefill(b *testing.B) {
	for _, test := range tests {
		b.Run(strconv.Itoa(test.count), func(b *testing.B) {
			q := stack.New()
			for n := 0; n < b.N; n++ {
				for n := 0; n < refillCount; n++ {
					for i := 0; i < test.count; i++ {
						q.Push(nil)
					}
					for q.Len() > 0 {
						tmp, tmp2 = q.Pop()
					}
				}
			}
		})
	}
}

func BenchmarkRefillFull(b *testing.B) {
	s := stack.New()
	for i := 0; i < fillCount; i++ {
		s.Push(nil)
	}

	for _, test := range tests {
		b.Run(strconv.Itoa(test.count), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for k := 0; k < refillCount; k++ {
					for i := 0; i < test.count; i++ {
						s.Push(nil)
					}
					for i := 0; i < test.count; i++ {
						tmp, tmp2 = s.Pop()
					}
				}
			}
		})
	}

	for s.Len() > 0 {
		tmp, tmp2 = s.Pop()
	}
}

func BenchmarkStable(b *testing.B) {
	s := stack.New()
	for i := 0; i < fillCount; i++ {
		s.Push(nil)
	}

	for _, test := range tests {
		b.Run(strconv.Itoa(test.count), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for i := 0; i < test.count; i++ {
					s.Push(nil)
					tmp, tmp2 = s.Pop()
				}
			}
		})
	}

	for s.Len() > 0 {
		tmp, tmp2 = s.Pop()
	}
}

func BenchmarkSlowIncrease(b *testing.B) {
	for _, test := range tests {
		b.Run(strconv.Itoa(test.count), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s := stack.New()
				for i := 0; i < test.count; i++ {
					s.Push(nil)
					s.Push(nil)
					tmp, tmp2 = s.Pop()
				}
				for s.Len() > 0 {
					tmp, tmp2 = s.Pop()
				}
			}
		})
	}
}

func BenchmarkSlowDecrease(b *testing.B) {
	s := stack.New()
	for _, test := range tests {
		items := test.count / 2
		for i := 0; i <= items; i++ {
			s.Push(nil)
		}
	}

	for _, test := range tests {
		b.Run(strconv.Itoa(test.count), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				for i := 0; i < test.count; i++ {
					s.Push(nil)
					tmp, tmp2 = s.Pop()
					if s.Len() > 0 {
						tmp, tmp2 = s.Pop()
					}
				}
			}
		})
	}

	for s.Len() > 0 {
		tmp, tmp2 = s.Pop()
	}
}
