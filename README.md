# stack [![Build Status](https://travis-ci.com/ef-ds/stack.svg?branch=master)](https://travis-ci.com/ef-ds/stack) [![codecov](https://codecov.io/gh/ef-ds/stack/branch/master/graph/badge.svg)](https://codecov.io/gh/ef-ds/stack) [![Go Report Card](https://goreportcard.com/badge/github.com/ef-ds/stack)](https://goreportcard.com/report/github.com/ef-ds/stack)  [![GoDoc](https://godoc.org/github.com/ef-ds/stack?status.svg)](https://godoc.org/github.com/ef-ds/stack)

Package stack implements a very fast and efficient general purpose Last-In-First-Out (LIFO) stack data structure that is specifically optimized to perform when used by Microservices and serverless services running in production environments. Internally, stack stores the elements in a dynamic growing semi-circular inverted singly linked list of arrays.


## Install
From a configured [Go environment](https://golang.org/doc/install#testing):
```sh
go get -u github.com/ef-ds/stack
```

If you are using dep:
```sh
dep ensure -add github.com/ef-ds/stack@1.0.1
```

We recommend to target only released versions for production use.


## How to Use
```go
package main

import (
	"fmt"

	"github.com/ef-ds/stack"
)

func main() {
	var s stack.Stack

	for i := 1; i <= 5; i++ {
		s.Push(i)
	}
	for s.Len() > 0 {
		v, _ := s.Pop()
		fmt.Println(v)
	}
}
```

Output:
```
5
4
3
2
1
```

Also refer to the [integration](integration_test.go) and [API](api_test.go) tests.



## Tests
Besides having 100% code coverage, stack has an extensive set of [unit](unit_test.go), [integration](integration_test.go) and [API](api_test.go) tests covering all happy, sad and edge cases.

When considering all tests, stack has over 4x more lines of testing code when compared to the actual, functional code.

Performance and efficiency are major concerns, so stack has an extensive set of benchmark tests as well comparing the stack performance with a variety of high quality open source stack implementations.

See the [benchmark tests](https://github.com/ef-ds/stack-bench-tests/blob/master/BENCHMARK_TESTS.md) for details.


## Performance
Stack has constant time (O(1)) on all its operations (Push/Pop/Back/Len). It's not amortized constant because stack never copies more than 256 (maxInternalSliceSize/2) items and when it expands or grow, it never does so by more than 512 (maxInternalSliceSize) items in a single operation.

Stack offers either the best or very competitive performance across all test sets, suites and ranges.

As a general purpose LIFO stack, stack offers, by far, the most balanced and consistent performance of all tested data structures.

See [performance](https://github.com/ef-ds/stack-bench-tests/blob/master/PERFORMANCE.md) for details.


## Design
The Efficient Data Structures (ef-ds) stack employs a new, modern stack design: a dynamic growing semi-circular inverted singly linked list of slices.

That means the [LIFO stack](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)) is a [singly-linked list](https://en.wikipedia.org/wiki/Singly_linked_list) where each node value is a fixed size [slice](https://tour.golang.org/moretypes/7). It is inverted singly linked list because each node points only to the previous one (instead of next) and it is semi-circular in shape because the first node in the linked list points to itself, but the last one points to nil.

![ns/op](testdata/stack.jpg?raw=true "stack Design")


### Design Considerations
Stack uses linked slices as its underlying data structure. The reason for the choice comes from two main observations of pure slice based stacks:

1. When the stack needs to expand to accommodate new values, [a new, larger slice needs to be allocated](https://en.wikipedia.org/wiki/Dynamic_array#Geometric_expansion_and_amortized_cost) and used
2. Allocating and managing large slices is expensive, especially in an overloaded system with little available physical memory

To help clarify the scenario, below is what happens when a slice based stack that already holds, say 1bi items, needs to expand to accommodate a new item.

Slice based implementation.

- Allocate a new, twice the size of the previous allocated one, say 2 billion positions slice
- Copy over all 1 billion items from the previous slice into the new one
- Add the new value into the first unused position in the new slice, position 1000000001

The same scenario for stack plays out like below.

- Allocate a new 1024 size slice
- Set the previous and next pointers
- Add the value into the first position of the new slice, position 0

The decision to use linked slices was also the result of the observation that slices goes to great length to provide predictive, indexed positions. A hash table, for instance, absolutely need this property, but not a stack. So stack completely gives up this property and focus on what really matters: add and retrieve from the edges (front/back). No copying around and repositioning of elements is needed for that. So when a slice goes to great length to provide that functionality, the whole work of allocating new arrays, copying data around is all wasted work. None of that is necessary. And this work costs dearly for large data sets as observed in the tests.


## Supported Data Types
Similarly to Go's standard library list, [list](https://github.com/golang/go/tree/master/src/container/list),
[ring](https://github.com/golang/go/tree/master/src/container/ring) and [heap](https://github.com/golang/go/blob/master/src/container/heap/heap.go) packages, stack supports "interface{}" as its data type. This means it can be used with any Go data types, including int, float, string and any user defined structs and pointers to interfaces.

The data types pushed into the stack can even be mixed, meaning, it's possible to push ints, floats and struct instances into the same stack.


## Safe for Concurrent Use
Stack is not safe for concurrent use. However, it's very easy to build a safe for concurrent use version of the stack. Impl7 design document includes an example of how to make impl7 safe for concurrent use using a mutex. stack can be made safe for concurret use using the same technique. Impl7 design document can be found [here](https://github.com/golang/proposal/blob/master/design/27935-unbounded-queue-package.md).


## Range Support
Just like the current container data structures such as [list](https://github.com/golang/go/tree/master/src/container/list),
[ring](https://github.com/golang/go/tree/master/src/container/ring) and [heap](https://github.com/golang/go/blob/master/src/container/heap/heap.go), stack doesn't support the range keyword for navigation.

However, the API offers two ways to iterate over the stack items. Either use "PopFront"/"PopBack" to retrieve the first current element and the second bool parameter to check for an empty queue.

```go
for v, ok := s.Pop(); ok; v, ok = s.Pop() {
    // Do something with v
}
```

Or use "Len" and "Pop" to check for an empty stack and retrieve the first current element.
```go
for s.Len() > 0 {
    v, _ := s.Pop()
    // Do something with v
}
```



## Why
We feel like this world needs improving. Our goal is to change the world, for the better, for everyone.

As software engineers at ef-ds, we feel like the best way we can contribute to a better world is to build amazing systems,
systems that solve real world problems, with unheard performance and efficiency.

We believe in challenging the status-quo. We believe in thinking differently. We believe in progress.

What if we could build queues, stacks, lists, arrays, hash tables, etc that are much faster than the current ones we have? What if we had a dynamic array data structure that offers near constant time deletion (anywhere in the array)? Or that could handle 1 million items data sets using only 1/3 of the memory when compared to all known current implementations? And still runs 2x as fast?

One sofware engineer can't change the world him/herself, but a whole bunch of us can! Please join us improving this world. All the work done here is made 100% transparent and is 100% free. No strings attached. We only require one thing in return: please consider benefiting from it; and if you do so, please let others know about it.


## Competition
We're extremely interested in improving stack and we're on an endless quest for better efficiency and more performance. Please let us know your suggestions for possible improvements and if you know of other high performance stacks not tested here, let us know and we're very glad to benchmark them.


## Releases
We're committed to a CI/CD lifecycle releasing frequent, but only stable, production ready versions with all proper tests in place.

We strive as much as possible to keep backwards compatibility with previous versions, so breaking changes are a no-go.

For a list of changes in each released version, see [CHANGELOG.md](CHANGELOG.md).


## Supported Go Versions
See [supported_go_versions.md](https://github.com/ef-ds/docs/blob/master/supported_go_versions.md).


## License
MIT, see [LICENSE](LICENSE).

"Use, abuse, have fun and contribute back!"


## Contributions
See [CONTRIBUTING.md](CONTRIBUTING.md).


## Roadmap
- Build tool to help find out the combination of firstSliceSize, sliceGrowthFactor and maxInternalSliceSize that will yield the best performance
- Find the fastest open source stacks and add them the bench tests
- Improve stack performance and/or efficiency by improving its design and/or implementation
- Build a high performance safe for concurrent use version of stack


## Contact
Suggestions, bugs, new queues to benchmark, issues with the stack, please let us know at ef-ds@outlook.com.
