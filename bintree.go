package forest

//
// bintree.go - implementation of the generic Binary Search Tree in Go.
//
//
// Jan 2015

import (
	"fmt"
)

// Node defines a single binary tree node
type Node struct {

	// Data, here a simple integer
	Data int

	// left and right children of the node and parent node
	left, right, parent *Node
}

// NewNode creates a new empty tree node.
func NewNode(data int) *Node { return &Node{data, nil, nil, nil} }

// String returns a human-readable representation of the node.
func (n *Node) String() string {
	s := fmt.Sprintf("Node: %d\n", n.Data)
	if n.left != nil {
		s = fmt.Sprintf("%s\tleft: %d\n", s, n.left.Data)
	} else {
		s = fmt.Sprintf("%s\tleft: empty\n", s)
	}
	if n.right != nil {
		s = fmt.Sprintf("%s\tright: %d\n", s, n.right.Data)
	} else {
		s = fmt.Sprintf("%s\tright: empty\n", s)
	}
	if n.parent != nil {
		s = fmt.Sprintf("%s\tparent: %d\n", s, n.parent.Data)
	} else {
		s = fmt.Sprintf("%s\tparent: empty\n", s)
	}
	return s
}

// BinTree defines the Binary Search Tree
type BinTree struct {

	// Root is a pointer to root node
	Root *Node

	// Len stores a number of elements in the tree (we keep this information separately)
	Len int
}

// NewBinTree creates a new empty BST.
func NewBinTree() *BinTree { return &BinTree{nil, 0} }

// aux method that finds the right element; returns nil if not found
func (bt *BinTree) find(root *Node, elem int) *Node {

	var cur = root
	for (cur != nil) && (cur.Data != elem) {
		if elem < cur.Data {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return cur
}

// Find searches for the element in the tree. Returns nil if not found.
func (bt *BinTree) Find(elem int) *Node {
	node := bt.find(bt.Root, elem)
	return node
}

// In checks whether the given element is a member of a tree.
func (bt *BinTree) In(elem int) bool {

	status := false

	if bt.Find(elem) != nil {
		status = true
	}
	return status
}

// Searches for the MIN element of the tree; it's the far left element.
func (bt *BinTree) findMinElem(node *Node) (*Node, bool) {

	var cur = node // we start at node

	// if tree is still empty, just return an error
	if cur == nil {
		return cur, false
	}

	for cur.left != nil {
		cur = cur.left
	}
	return cur, true
}

// Min returns the MIN element in the tree. This element is located as the far left element.
func (bt *BinTree) Min() (int, bool) {

	var cur *Node
	var ok bool
	if cur, ok = bt.findMinElem(bt.Root); !ok {
		return 0, false
	}
	return cur.Data, true
}

// Searches for the MAX element of the tree; it's the far right element.
func (bt *BinTree) findMaxElem(node *Node) (*Node, bool) {

	var cur = node // we start at node (usually root)

	// if tree is still empty, just return an error
	if cur == nil {
		return nil, false
	}

	for cur.right != nil {
		cur = cur.right
	}
	return cur, true

}

// Max Return the MAX element in the tree. This element is located as the far right element.
func (bt *BinTree) Max() (int, bool) {

	var cur *Node
	var ok bool
	if cur, ok = bt.findMaxElem(bt.Root); !ok {
		return 0, false
	}
	return cur.Data, true
}

// Traverse implements the tree traversing left iteratively
// FIXME: there's obviously a bug somewhere in this method...
func (bt *BinTree) Traverse() {

	var cur = bt.Root
	var prev *Node
	var next *Node

	fmt.Print("        Traversing: ")
	for cur != nil {
		fmt.Printf("%d  ", cur.Data)
		switch {
		case prev == cur.parent: // we are in parent node, we try to go left, then right, then back to parent
			if cur.left != nil {
				next = cur.left
			} else if cur.right != nil {
				next = cur.right
			} else {
				next = cur.parent
			}
		case prev == cur.left: // we are in left element: we try to go right, then back to parent
			if cur.right != nil {
				next = cur.right
			} else {
				next = cur.parent
			}
		default: // we are in right element, go back to parent
			next = cur.parent
		}
		prev = cur
		cur = next
	}
	fmt.Println()
}

// TraverseReverse implements tree traversing right (reverse order) iteratively.
// FIXME: there's obviously a bug somewhere in this method...
func (bt *BinTree) TraverseReverse() {

	var cur = bt.Root
	var prev *Node
	var next *Node

	fmt.Print("Reverse-Traversing: ")
	for cur != nil {

		fmt.Printf("%d  ", cur.Data)

		switch {
		case prev == cur.parent: // we are in parent node, we try to go right, then left, then back to parent
			if cur.right != nil {
				next = cur.right
			} else if cur.left != nil {
				next = cur.left
			} else {
				next = cur.parent
			}
		case prev == cur.right: // we are in right element: we try to go left, then back to parent
			if cur.left != nil {
				next = cur.left
			} else {
				next = cur.parent
			}
		default: // we are in left element, go back to parent
			next = cur.parent
		}
		prev = cur
		cur = next
	}
	fmt.Println()
}

// Insert inserts a new element into the tree.
func (bt *BinTree) Insert(node *Node) {

	// if node is empty, just assign the current node to root
	if bt.Root == nil {
		bt.Root = node
		bt.Len++
		return
	}

	// traverse the tree to find the right place to insert
	cur := bt.Root
	prev := bt.Root
	for cur != nil {

		if cur.Data == node.Data { // node already exists...
			return
		}

		prev = cur

		if node.Data < cur.Data {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}

	// now insert the new node
	switch {

	case node.Data < prev.Data:
		prev.left = node
		bt.Len++
		node.parent = prev

	case node.Data > prev.Data:
		prev.right = node
		bt.Len++
		node.parent = prev

	default: // if element already exists, do nothing...

	}
	return
}

// Add adds an element with value 'val' to tree.
func (bt *BinTree) Add(val int) {
	node := NewNode(val)
	bt.Insert(node)
}

// Recursive traverse from Min to Max value (ascending sort).
func (bt *BinTree) sortA(node *Node, sorted *[]int) {

	if node == nil {
		return
	}
	bt.sortA(node.left, sorted) // traverse left subtree
	*sorted = append(*sorted, node.Data)
	bt.sortA(node.right, sorted) // and  right subtree
}

// Recursive traverse from Max to Min value (descending sort).
func (bt *BinTree) sortD(node *Node, sorted *[]int) {

	if node == nil {
		return
	}
	bt.sortD(node.right, sorted)
	*sorted = append(*sorted, node.Data)
	bt.sortD(node.left, sorted)
}

// SortA sorts the elements in the tree in ascending order (from min element to max). Returns a slice of ints.
func (bt *BinTree) SortA() []int {

	var sorted []int
	bt.sortA(bt.Root, &sorted)
	return sorted
}

// SortD sorts the elements in the tree in descending order (from max element to min). Returns a slice of ints.
func (bt *BinTree) SortD() []int {

	var sorted []int
	bt.sortD(bt.Root, &sorted)
	return sorted
}

// Delete removes an element from the tree.
func (bt *BinTree) Delete(node *Node) {

	if elem := bt.find(bt.Root, node.Data); elem != nil {

		parent := elem.parent

		switch {

		case elem.left == nil && elem.right == nil: // found element is a leaf
			if parent.right == elem {
				parent.right = nil
			} else {
				parent.left = nil
			}

		case elem.left == nil: // found element has only right child
			parent.right = elem.right
			elem.right.parent = parent

		case elem.right == nil: // found element has only left child
			parent.left = elem.left
			elem.left.parent = parent

		default: // in general, found element has both children
			min, _ := bt.findMinElem(elem.right) // find MIN value in right subtree
			elem.Data = min.Data                 // just trade the MIN value of right subtree
			if min.parent.left == min {          // and we can safely delete the found MIN value node
				min.parent.left = nil
			} else {
				min.parent.right = nil
			}
		}
		bt.Len-- // we have one element less...
	}
	// if no element is found (elem == nil), just return; works for empty root, too
}

// Remove removes the element from the BST.
func (bt *BinTree) Remove(data int) {
	n := NewNode(data)
	bt.Delete(n)
}

// Balance rebalances the BST using the Day-Stout-Warren algorithm.
func (bt *BinTree) Balance() { Balance(bt) }
