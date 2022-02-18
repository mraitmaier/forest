package forest

// avltree.go - this is an example implementation of AVL tree; currently uses ints as data.
// Basically this is an exercise in style that can be quite useful in the near future.

import (
	"fmt"
)

// AvlNode defines the single AVL Tree node
type AvlNode struct {

	// Data
	Data int

	// left, right and parent node
	left, right, parent *AvlNode

	// height of the node (in th tree)
	height int
}

// NewAvlNode creates new AVL tree node.
func NewAvlNode(data int) *AvlNode { return &AvlNode{data, nil, nil, nil, 1} }

// String  retunrs a human-readable representation of the node.
func (n *AvlNode) String() string {

	s := fmt.Sprintf("Node: %d, h=%d ", n.Data, n.height)
	if n.left != nil {
		s = fmt.Sprintf("%s left=%v ", s, n.left)
	} else {
		s = fmt.Sprintf("%s left=Empty ", s)
	}
	if n.right != nil {
		s = fmt.Sprintf("%s right=%v ", s, n.right)
	} else {
		s = fmt.Sprintf("%s right=Empty ", s)
	}
	if n.parent != nil {
		s = fmt.Sprintf("%s parent=%v ", s, n.right)
	} else {
		s = fmt.Sprintf("%s parent=Empty ", s)
	}
	return s
}

// Balance calculates the balance factor.
func (n *AvlNode) Balance() int { return height(n.left) - height(n.right) }

// AvlTree defines the AVL Tree structure.
type AvlTree struct {

	// Root represents the root of the tree.
	Root *AvlNode

	// Len conviently holds the number of nodes in the tree.
	Len int
}

// NewAvlTree creates a new empty tree.
func NewAvlTree() *AvlTree { return &AvlTree{nil, 0} }

// Returns a found node, otherwise returns nil.
func (t *AvlTree) find(node *AvlNode, data int) *AvlNode {

	var cur = node
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

// SearchFor searches for a given data in the tree. Returns nil if not found.
func (t *AvlTree) SearchFor(data int) *AvlNode { return t.find(t.Root, data) }

// In checks if given element is a member of a tree.
func (t *AvlTree) In(data int) bool {

	if t.SearchFor(data) != nil {
		return true
	}
	return false
}

// This is a rotate-right operation method.
func (t *AvlTree) rotateR(node *AvlNode) *AvlNode {

	left := node.left
	// rotate
	node.left = left.right
	left.right = node
	left.parent = node.parent
	node.parent = left

	// update heights
	node.height = largerOf(height(node.left), height(node.right)) + 1
	left.height = largerOf(height(left.left), height(left.right)) + 1

	// new root of the subtree
	return left
}

// This is a rotate-left operation method.
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

// After rotation, parent node (usually) points to wrong node, fix this.
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

// Rebalances the tree.
// XXX: there's probably a bug hiding in this method; some inserted elements are lost
// XXX: since ordinary BST insert is probably working OK, the rebalance procedure must be broken - updateParents()???
func (t *AvlTree) balance(node *AvlNode) *AvlTree {

	var parent *AvlNode
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

// Inserts a new node into the tree.
// Basically the same as the generic BST tree insert. Balancing is done in its own method.
func (t *AvlTree) insert(node *AvlNode) *AvlNode {

	// if the tree is empty, just define the root..
	if t.Root == nil {
		t.Root = node
		t.Len++
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
		t.Len++
		node.parent = prev
		fmt.Printf("DEBUG element %d inserted...\n", node.Data) // DEBUG

	case node.Data > prev.Data:
		prev.right = node
		t.Len++
		node.parent = prev
		fmt.Printf("DEBUG element %d inserted...\n", node.Data) // DEBUG

	default: // do nothing when element already exists
		fmt.Printf("DEBUG element %d already exists...\n", node.Data) // DEBUG
	}
	return node
}

// Insert inserts a new node into the tree. Insertion is done iteratively, not using recursion (as usual).
// We try to avoid the problem with recursion depth, when tree grows really large.
func (t *AvlTree) Insert(node *AvlNode) {

	// insert new node operation, the same as generic BST
	node = t.insert(node)
	// now we rebalance the tree
	t = t.balance(node)
}

// Add adds a new element to the tree.
func (t *AvlTree) Add(data int) {
	n := NewAvlNode(data)
	t.Insert(n)
}

// TraversePostOrder executes the post-order traversal of the tree.
func (t *AvlTree) TraversePostOrder() {
	fmt.Print("Post-order: ")
	t.postorder(t.Root)
	fmt.Println()
}

// Post-order traversal
func (t *AvlTree) postorder(node *AvlNode) {

	if node == nil {
		return
	}
	t.postorder(node.left)
	t.postorder(node.right)
	fmt.Printf("%d ", node.Data)
}

// TraversePreOrder executes the pre-order traversal of the tree.
func (t *AvlTree) TraversePreOrder() {
	fmt.Print(" Pre-order: ")
	t.preorder(t.Root)
	fmt.Println()
}

// Pre-order traversal
func (t *AvlTree) preorder(node *AvlNode) {

	if node == nil {
		return
	}
	fmt.Printf("%d ", node.Data)
	t.preorder(node.left)
	t.preorder(node.right)
}

// Traverse executes the in-order traversal of the tree.
func (t *AvlTree) Traverse() {
	fmt.Print("  In-order: ")
	t.inorder(t.Root)
	fmt.Println()
}

// traversing left iteratively
func (t *AvlTree) inorder(node *AvlNode) {

	if node == nil {
		return
	}
	t.inorder(node.left)
	fmt.Printf("%d ", node.Data)
	t.inorder(node.right)
}

// Recursive traverse from Min to Max value (ascending sort).
func (t *AvlTree) sortA(node *AvlNode, sorted *[]int) {

	if node == nil {
		return
	}
	t.sortA(node.left, sorted) // traverse left subtree
	*sorted = append(*sorted, node.Data)
	t.sortA(node.right, sorted) // and  right subtree
}

// Recursive traverse from Max to Min value (descending sort).
func (t *AvlTree) sortD(node *AvlNode, sorted *[]int) {

	if node == nil {
		return
	}
	t.sortD(node.right, sorted)
	*sorted = append(*sorted, node.Data)
	t.sortD(node.left, sorted)
}

func (t *AvlTree) Sorted() []int {
	var s []int
	t.sortA(t.Root, &s)
	return s
}

/*
// iterative implementation of in-order traversal --- XXX: somewhere, there's a bug hiding...
func (t *AvlTree) Sorted() []int {

	var cur *AvlNode = t.Root
	var prev *AvlNode = nil
	var next *AvlNode = nil
    var m []int

    // go far left to the MIN member
	for cur.left != nil {
		cur = cur.left
	}

    // now let's traverse the tree and build a sorted slice
	for cur != nil {
		switch {
		case prev == cur.parent: // we are in parent node, we try to go left, then right, then back to parent
			if cur.left != nil {
				next = cur.left
			} else if cur.right != nil {
                m = append(m, cur.Data)
				next = cur.right
			} else {
                m = append(m, cur.Data)
				next = cur.parent
			}
		case prev == cur.left: // we are in left element: we try to go right, then back to parent
			if cur.right != nil {
                m = append(m, cur.Data)
				next = cur.right
			} else {
                m = append(m, cur.Data)
				next = cur.parent
			}
		default: // we are in right element, go back to parent
			next = cur.parent
		}
		prev = cur
		cur = next
	}
    return m
}
*/

// Searches for the MIN element of the tree; it's the deepest, far left element.
func (t *AvlTree) findMinElem(node *AvlNode) (*AvlNode, bool) {

	var cur = node // we start at node

	// if tree is still empty, just return an error
	if cur == nil {
		return cur,false 
	}

	for cur.left != nil {
		cur = cur.left
	}
	return cur, true
}

// Min returns the MIN element in the tree. This element is located far left in the tree.
func (t *AvlTree) Min() (int, error) {

    if cur, ok := t.findMinElem(t.Root); ok {
		return cur.Data, nil
	}
	return 0, fmt.Errorf("Cannot find Min element")
}

// Searches for the MAX element of the tree; it's the deepest, far right element.
func (t *AvlTree) findMaxElem(node *AvlNode) (*AvlNode, error) {

	var cur = node // we start at node (usually root)

	// if tree is still empty, just return an error
	if cur == nil {
		return nil, fmt.Errorf("Empty tree")
	}

	for cur.right != nil {
		cur = cur.right
	}
	return cur, nil

}

// Max returns the MAX element in the tree. This element is located far right in the tree.
func (t *AvlTree) Max() (int, error) {

	var cur *AvlNode
	var err error
	if cur, err = t.findMaxElem(t.Root); err != nil {
		return 0, err
	}
	return cur.Data, nil
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
		}
		return node.left
	}
}

// Rebalance the tree after deleting a node.
func (t *AvlTree) balanceD(node *AvlNode) *AvlTree {

	// if we had only one node, the tree is now empty, return
	if node == nil {
		return t
	}

	var parent *AvlNode
	var child *AvlNode
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

// Delete removes an element from the tree.
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
			if min, found := t.findMinElem(elem.right); found { // find MIN value in right subtree
				elem.Data = min.Data // exchange the element value and the MIN value of right subtree
				if min.parent.left == min {
					min.parent.left = nil //
				} else {
					min.parent.right = nil
				}
			}
		}
		t.Len-- // we have one element less...
	}

	// Now we have rebalance the tree... But only if the element was found in the tree.
	if elem != nil {
		t = t.balanceD(elem)
	}
	// NOTE: if no element is found (elem == nil), just return; works for empty root, too
	return t
}

// Remove removes the node with the given value.
func (t *AvlTree) Remove(data int) *AvlTree {
	n := NewAvlNode(data)
	return t.Delete(n)
}

/* some additional aux functions... */

// Return height of the node, (nil-)safe from empty nodes.
func height(node *AvlNode) int {
	if node == nil {
		return 0
	}
	return node.height
}

// Just  return the higher integer of the two.
func largerOf(x, y int) int {
	if x > y {
		return x
	}
	return y
}
