package forest
/*
 * set.go - example implementation of the set using AVL tree.
 */

//import ()

// IntSet represents the set of integers and is implemented as an AVL tree.
type IntSet struct {
    // AvlTree is embedded as set of integers   
    AvlTree
}

// NewIntSet creates a new instance of IntSet.
func NewIntSet() *IntSet { return &IntSet{ *NewAvlTree() } }

// Add insert new value into given set.
func (s *IntSet) Add(val int) { s.AvlTree.Add(val) }

// Remove
func (s *IntSet) Remove(val int) { s.AvlTree.Remove(val) }

// Count returns the number of members in set.
func (s *IntSet) Len() int { return s.AvlTree.Len }

//IsEmpty returns status if a set is an empty one (has 0 members).
func (s *IntSet) IsEmpty() bool {
    if s.AvlTree.Len == 0 {
        return true
    }
    return false
}

// In  returns the status whether the given value is a member of a set or not.
func (s *IntSet) In(val int) bool { return s.AvlTree.In(val) }

// Min returns the smallest member of the set.
func (s *IntSet) Min() int {
    m, _ := s.AvlTree.Min()
    return m
}

// Max returns the largest member of the set.
func (s *IntSet) Max() int {
    m , _ := s.AvlTree.Max()
    return m
}

// 
func (s *IntSet) Members() []int { return s.AvlTree.Sorted() }
//
//func (s *IntSet) Copy() *IntSet {
//}

func (s *IntSet) Display() { s.AvlTree.Traverse() }

