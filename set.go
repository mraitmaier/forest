package forest

import ()

// IntSet represents the set of integers and is implemented as an AVL tree.
type IntSet struct {
    // AvlTree is embedded as set of integers   
    AvlTree
}

// NewIntSet creates a new instance of IntSet.
func NewIntSet() *IntSet { return &IntSet{ *NewAvlTree() } }

// Add insert new value into given set.
func (s *IntSet) Add(val int) { s.Add(val) }

// Remove
func (s *IntSet) Remove(val int) { s.Remove(val) }

// Count returns the number of members in set.
func (s *IntSet) Count() int { return s.Len }

//IsEmpty returns status if a set is an empty one (has 0 members).
func (s *IntSet) IsEmpty() bool {
    if s.Len == 0 {
        return true
    }
    return false
}

// In  returns the status whether the given value is a member of a set or not.
func (s *IntSet) In(val int) bool { return s.In(val) }

// Min returns the smallest member of the set.
func (s *IntSet) Min() int { return s.Min() }

// Max returns the largest member of the set.
func (s *IntSet) Max() int { return s.Max() }

//
