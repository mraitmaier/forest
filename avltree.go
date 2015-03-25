//

package forest

import (
	"fmt"
)

// Define the AVL Tree node
type AvlNode struct {

	// Data
	Data int

	// left, right and parent node
	left, right, parent *AvlNode

	// height of the node (in th tree)
	height int
}

// Create a new AVL tree node.
func NewAvlNode(data int) *AvlNode { return &AvlNode{data, nil, nil, nil, 1} }

// A string representation of the node.
func (n *AvlNode) String() string {

	s := fmt.Sprintf("Node: %d, h=%d ", n.Data, n.height)
	if n.left != nil {
		s = fmt.Sprintf("%s left=&v ", s, n.left)
	} else {
		s = fmt.Sprintf("%s left=Empty ", s)
	}
	if n.right != nil {
		s = fmt.Sprintf("%s right=&v ", s, n.right)
	} else {
		s = fmt.Sprintf("%s right=Empty ", s)
	}
	if n.parent != nil {
		s = fmt.Sprintf("%s parent=&v ", s, n.right)
	} else {
		s = fmt.Sprintf("%s parent=Empty ", s)
	}
	return s
}

// Get balance factor
func (n *AvlNode) Balance() int { return height(n.left) - height(n.right) }

// Define the AVL Tree
type AvlTree struct {

	// The root of the tree
	Root *AvlNode

	// Number of nodes in the tree
	Len int
}

// Create a new empty tree.
func NewTree() *AvlTree { return &AvlTree{nil, 0} }

// Return a found node, otherwise return nil.
func (t *AvlTree) find(node *AvlNode, data int) *AvlNode {

	var cur *AvlNode = node
	for cur != nil {
		if cur.Data == data {
			break
		} else if cur.Data > data {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return cur
}

// Search for a given data in the tree
func (t *AvlTree) SearchFor(data int) *AvlNode { return t.find(t.Root, data) }

// Check if given element is a member of a tree.
func (t *AvlTree) In(data int) bool {

	if t.SearchFor(data) != nil {
		return true
	}
	return false
}

// Rotate right operation
func (t *AvlTree) rotateR(node *AvlNode) *AvlNode {

	left := node.left
	// rotate
	node.left = left.right
	left.right = node
	left.parent = node.parent
	node.parent = left

	// update hights
	node.height = largerOf(height(node.left), height(node.right)) + 1
	left.height = largerOf(height(left.left), height(left.right)) + 1

	// new root of the subtree
	return left
}

// Rotate left operation
func (t *AvlTree) rotateL(node *AvlNode) *AvlNode {

	right := node.right
	// rotate
	node.right = right.left
	right.left = node
	right.parent = node.parent
	node.parent = right

	// update heights
	node.height = largerOf(height(node.left), height(node.right)) + 1
	right.height = largerOf(height(right.left), height(right.right)) + 1

	// new root of the subtree
	return right
}

// After rotation, parent node (usually) points  to wrong node, fix this.
func (t *AvlTree) updateParent(node, prev, parent *AvlNode) *AvlNode {

	switch {

	case parent == nil: // do nothing here, just return

	case prev == parent.right:
		parent.right = node

	case prev == parent.left:
		parent.left = node
	}
	return parent
}

// Rebalance the tree.
func (t *AvlTree) balance(node *AvlNode) *AvlTree {

	var parent *AvlNode = nil
	cur := node.parent // node's parents and up might be unbalanced...

	for cur != nil {

		cur.height = 1 + largerOf(height(cur.left), height(cur.right)) // we must update height of the current node up to root
		bal := cur.Balance()                                           // get balance factor for the node

		// remember node's parent and to which direction it is pointing... parent, I mean...
		parent = cur.parent
		if bal > 1 {
			if cur.left.Balance() > 0 { // left-left case
				cur = t.rotateR(cur)
				parent = t.updateParent(cur, cur.right, parent)
			} else { // left-right case
				cur.left = t.rotateL(cur.left)
				cur = t.rotateR(cur)
				parent = t.updateParent(cur, cur.right, parent)
			}
		} else if bal < -1 {
			if cur.right.Balance() < 0 { // right-right case
				cur = t.rotateL(cur)
				parent = t.updateParent(cur, cur.left, parent)
			} else { // right-left case
				cur.right = t.rotateR(cur.right)
				cur = t.rotateL(cur)
				parent = t.updateParent(cur, cur.left, parent)
			}
		}
		if parent == nil {
			t.Root = cur
		} // we might have a new root...
		cur = cur.parent
	}
	return t
}

// Insert a new node into the tree. Basically the same as the generic BST tree insert. Balancing is done in its own method.
func (t *AvlTree) insert(node *AvlNode) *AvlNode {

	// if the tree is empty, just define the root..
	if t.Root == nil {
		t.Root = node
		t.Len += 1
		return node
	}

	// traverse the tree to find the right spot
	cur := t.Root
	prev := t.Root
	for cur != nil {

		// if data already exists, just return without hesitation...
		if cur.Data == node.Data {
			return node
		}

		prev = cur
		if node.Data < cur.Data {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}

	// insert the new element
	switch {

	case node.Data < prev.Data:
		prev.left = node
		t.Len += 1
		node.parent = prev

	case node.Data > prev.Data:
		prev.right = node
		t.Len += 1
		node.parent = prev

	default: // do nothing when element already exists
	}
	return node
}

// Insert a new node into the tree. Insertion is done iteratively, not using recursion (as usual).
// We try to avoid the problem with recursion depth, when tree grows really large.
func (t *AvlTree) Insert(node *AvlNode) {

	// insert new node operation, the same as generic BST
	node = t.insert(node)
	// now we rebalance the tree
	t = t.balance(node)
}

// Add a new element to the tree.
func (t *AvlTree) Add(data int) {
	n := NewAvlNode(data)
	t.Insert(n)
}

// Post-order traversal of the tree.
func (t *AvlTree) TraversePostOrder() {
	fmt.Print("Post-order: ")
	t.postorder(t.Root)
	fmt.Println()
}

// Post-order traversal
func (bt *AvlTree) postorder(node *AvlNode) {

	if node == nil {
		return
	}
	bt.postorder(node.left)
	bt.postorder(node.right)
	fmt.Printf("%d ", node.Data)
}

// Pre-order traversal of the tree.
func (t *AvlTree) TraversePreOrder() {
	fmt.Print(" Pre-order: ")
	t.preorder(t.Root)
	fmt.Println()
}

// Pre-order traversal
func (bt *AvlTree) preorder(node *AvlNode) {

	if node == nil {
		return
	}
	fmt.Printf("%d ", node.Data)
	bt.preorder(node.left)
	bt.preorder(node.right)
}

// In-order traversal of the tree.
func (t *AvlTree) Traverse() {
	fmt.Print("  In-order: ")
	t.inorder(t.Root)
	fmt.Println()
}

// traversing left iteratively
func (bt *AvlTree) inorder(node *AvlNode) {

	if node == nil {
		return
	}
	bt.inorder(node.left)
	fmt.Printf("%d ", node.Data)
	bt.inorder(node.right)
}

/*
// iterative implementation of in-order traversal --- XXX: somewhere, there's a bug hiding...
func (bt *AvlTree) InorderIter(node *AvlNode) {
	var cur *AvlNode = bt.Root
	var prev *AvlNode = nil
	var next *AvlNode = nil

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
*/

// Searches for the MIN element of the tree; it's the deepest, far left element.
func (t *AvlTree) findMinElem(node *AvlNode) (*AvlNode, error) {

	var cur *AvlNode = node // we start at node

	// if tree is still empty, just return an error
	if cur == nil {
		return cur, fmt.Errorf("Empty tree")
	}

	for cur.left != nil {
		cur = cur.left
	}
	return cur, nil
}

// Return the Min() element in the tree. This element is deepest and far left.
func (t *AvlTree) Min() (int, error) {

	if cur, err := t.findMinElem(t.Root); err != nil {
		return 0, nil
	} else {
		return cur.Data, err
	}
}

// Searches for the MAX element of the tree; it's the deepest, far right element.
func (t *AvlTree) findMaxElem(node *AvlNode) (*AvlNode, error) {

	var cur *AvlNode = node // we start at node (usually root)

	// if tree is still empty, just return an error
	if cur == nil {
		return nil, fmt.Errorf("Empty tree")
	}

	for cur.right != nil {
		cur = cur.right
	}
	return cur, nil

}

// Return the Max() element in the tree. This element is far right.
func (t *AvlTree) Max() (int, error) {

	if cur, err := t.findMaxElem(t.Root); err != nil {
		return 0, err
	} else {
		return cur.Data, nil
	}
}

// find the larger child of the node.
func (t *AvlTree) findLarger(node *AvlNode) *AvlNode {

	if node == nil {
		return nil
	} // safeguard

	switch {

	case node.right == nil && node.left == nil: // node is a leaf
		return nil

	case node.right == nil: // node has only left child
		return node.left

	case node.left == nil: // node has only left child
		return node.right

	default: // node has both children
		if node.right.Balance() > node.left.Balance() {
			return node.right
		} else {
			return node.left
		}
	}
}

// Rebalance the tree after deleting a node.
func (t *AvlTree) balanceD(node *AvlNode) *AvlTree {

	// if we had only one node, the tree is now empty, return
	if node == nil {
		return t
	}

	var parent *AvlNode = nil
	var child *AvlNode = nil
	cur := node.parent // node's parents and up might be unbalanced...

	for cur != nil {

		cur.height = 1 + largerOf(height(cur.left), height(cur.right)) // we must update height of the current node up to root
		bal := cur.Balance()                                           // get balance factor for the node

		// remember node's parent and to which direction it is pointing... parent, I mean...
		parent = cur.parent
		child = t.findLarger(cur)
		if bal > 1 {
			if child.Balance() > 0 { // left-left case
				cur = t.rotateR(cur)
				parent = t.updateParent(cur, cur.right, parent)
			} else { // left-right case
				cur.left = t.rotateL(cur.left)
				cur = t.rotateR(cur)
				parent = t.updateParent(cur, cur.right, parent)
			}
		} else if bal < -1 {
			if child.Balance() < 0 { // right-right case
				cur = t.rotateL(cur)
				parent = t.updateParent(cur, cur.left, parent)
			} else { // right-left case
				cur.right = t.rotateR(cur.right)
				cur = t.rotateL(cur)
				parent = t.updateParent(cur, cur.left, parent)
			}
		}
		if parent == nil { // we might have a new root...
			t.Root = cur
		}
		cur = cur.parent
	}
	return t
}

// Delete an element from the tree.
func (t *AvlTree) Delete(node *AvlNode) *AvlTree {

	var elem *AvlNode
	// find the element
	if elem = t.find(t.Root, node.Data); elem != nil {

		parent := elem.parent

		switch {

		case elem.left == nil && elem.right == nil: // found element is a leaf
			if parent.right == elem {
				parent.right = nil
			} else {
				parent.left = nil
			}

		case elem.left == nil: // found element has only right child
			if parent.left == elem {
				parent.left = elem.right
			} else {
				parent.right = elem.right
			}
			elem.right.parent = parent

		case elem.right == nil: // found element has only left child
			if parent.left == elem {
				parent.left = elem.left
			} else {
				parent.right = elem.left
			}
			elem.left.parent = parent

		default: // in general, found element has both children
			min, _ := t.findMinElem(elem.right)       // find MIN value in right subtree
			elem.Data, min.Data = min.Data, elem.Data // exchange the element value and the MIN value of right subtree
			if min.right != nil {
				min.Data = min.right.Data //
				min.right = nil
			} else {
				min.parent.right = nil
			}
		}
		t.Len -= 1 // we have one element less...
	}

	// Now we have rebalance the tree... But only if the element was found in the tree.
	if elem != nil {
		t = t.balanceD(elem)
	}
	// NOTE: if no element is found (elem == nil), just return; works for empty root, too
	return t
}

// Remove the node with the given value.
func (t *AvlTree) Remove(data int) *AvlTree {
	n := NewAvlNode(data)
	return t.Delete(n)
}

/* some additional aux functions... */

// Return height of the node, (nil-)safe from empty nodes.
func height(node *AvlNode) int {
	if node == nil {
		return 0
	} else {
		return node.height
	}
}

// Just  return the higher integer of the two.
func largerOf(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}
