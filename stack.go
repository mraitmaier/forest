// Package stack is a go implementation of the generic stack using slices.
//
// This is general implementation of the stack using slices and empty interface (interface{}) elements. This is convenient as one
// can use elements of any type to build a stack as long as one is carefull. It is advisable to use elements of the same type,
// though...
// Example usage is seen from the commented main function.
package stack

// Miran R., Jan2015

/*
package main

import (
    "fmt"
)
*/

// Stack is a generic implementation of the stack.
type Stack []interface{}

// IsEmpty checks if stack is empty.
func (s Stack) IsEmpty() bool { return len(s) == 0 }

// Len returns the length of the stack.
func (s Stack) Len() int { return len(s) }

// Peek  peeks into stack and returns the last value (value is not removed from stack).
func (s Stack) Peek() interface{} { return s[len(s)-1] }

// Push pushes the new value onto a stack.
func (s *Stack) Push(i interface{}) { *s = append(*s, i) }

// Pop pops the last value from the stack.
func (s *Stack) Pop() interface{} {
	d := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return d
}

/*
func main() {

    var s Stack
    fmt.Printf("Is stack empty? %t\n", s.IsEmpty())
    fmt.Println(">> Pushing")
    s.Push("something")
    fmt.Printf("Is stack empty? %t\n", s.IsEmpty())
    fmt.Printf("Stack length: %d\n", s.Len())
    str := s.Peek()
    fmt.Printf("Peek element: %v\n", str)
    fmt.Println(">> Pushing")
    s.Push(0)
    elem := s.Peek()
    fmt.Printf("Stack length: %d\n", s.Len())
    fmt.Printf("Peek element: %v\n", elem)
    fmt.Println(">> Pushing")
    s.Push(3)
    s.Push(4)
    s.Push(5)
    fmt.Printf("Stack length: %d\n", s.Len())
    fmt.Printf("Peek element: %v\n", s.Peek())
    fmt.Println(">> Popping")
    elem = s.Pop()
    fmt.Printf("Stack length: %d\n", s.Len())
    fmt.Printf("Popped element: %d\n", elem)
    fmt.Printf("Peek element: %v\n", s.Peek())
    fmt.Println(">> Popping")
    elem = s.Pop()
    elem = s.Pop()
    elem = s.Pop()
    fmt.Printf("Stack length: %d\n", s.Len())
    fmt.Printf("Popped element: %d\n", elem)
    fmt.Printf("Peek element: %v\n", s.Peek())
    fmt.Println(">> Popping")
    elem = s.Pop()
    fmt.Printf("Popped element: %q\n", elem)
    if !s.IsEmpty() {
        fmt.Printf("Peek element: %v\n", s.Peek())
    } else {
        fmt.Printf("Stack IS empty!\n")
    }
}
*/
