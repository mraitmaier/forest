package forest

// This is implementation of the scapegoat tree in Go.
//
// FIXME: not finished yet!

import (
	"fmt"
	"math"
)

/* NOTE: Node is defined in <bintree.go>; no need to duplicate... */

// The default value of Alpha factor
const defAlpha float64 = 0.667

// helper function that calculates node size: this is number of nodes in subtree rooted by node
func size(node *Node) int {
	if node == nil {
		return 0
	}
	return size(node.left) + size(node.right) + 1
}

// SgTree defines the Scapegoat Tree structure.
type SgTree struct {

	// Root pointd to the root of the tree
	Root *Node

	// Size represents the number of nodes in the tree
	Size int

	// alpha factor for the tree: must be between 0.5 and 1.0; this is user-definable parameter;
	// we set this factor to value of 'defAlpha' constant (see above) by default
	alpha float64

	// Maximum size of the tree before it is rebuild
	maxSize int

	// channels defined
	//insch, rmch, fch, foundch chan *Node
	//quitch                    chan int
}

// NewSgTree creates a new empty scape-goat tree.
func NewSgTree() *SgTree { return &SgTree{nil, 0, defAlpha, 0} }

/*
func NewSgTree() *SgTree {
	//    return &SgTree{nil, 0, defAlpha, 0,
	//                         make(chan *Node), make(chan *Node), make(chan *Node), make(chan *Node), make(chan int) }
	n := &SgTree{
		Root:    nil,
		Size:    0,
		alpha:   defAlpha,
		maxSize: 0,
		insch:   make(chan *Node),
		rmch:    make(chan *Node),
		fch:     make(chan *Node),
		foundch: make(chan *Node),
		quitch:  make(chan int),
	}
	go n.do() // start a dispather goroutine
	return n
}
*/

// SetAlpha changes the value of Alpha factor.
// For value of 0.5, scapegoat tree is equivalent to perfectly balanced BST; for value of 1.0,
// a linked list is considered as balanced BST.
func (t *SgTree) SetAlpha(val float64) error {
	// check Alpha value, must be between 0.5 and 1.0
	if val < 0.5 || val > 1.0 {
		return fmt.Errorf("Alpha value must be between 0.5 and 1.0")
	}
	t.alpha = val
	return nil
}

// GetAlpha returns the value of Alpha factor used by algorithm.
func (t *SgTree) GetAlpha() float64 { return t.alpha }

// Calculate tree alpha-height value.
func (t *SgTree) heightFactor(count int) float64 {
	return math.Log10(float64(count)) / math.Log10(1/t.alpha)
}

// IsEmpty checks if tree is empty.
func (t *SgTree) IsEmpty() bool {
	if t.Root == nil {
		return true
	}
	return false
}

// Clear clears (destroys) the complete tree. The result of operation is empty tree.
func (t *SgTree) Clear() {
	t.Root = nil
	t.Size = 0
	t.maxSize = 0
}

// Return a node to be found. If not found, return nil.
func (t *SgTree) find(root *Node, data int) *Node {

	cur := root
	for cur != nil {
		if cur.Data == data { // found!
			break
		} else if cur.Data > data {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return cur
}

// SearchFor searches for a given data in the tree.
func (t *SgTree) SearchFor(data int) *Node { return t.find(t.Root, data) }

// In checks if given element is a member of a tree.
func (t *SgTree) In(data int) bool {

	if t.SearchFor(data) != nil {
		return true
	}
	return false
}

// Insert a new node into the tree. Basically the same as the generic BST tree insert. Returns (pointer to) the
// inserted node and the depth of the node. Tree rebuild is done in its own method.
func (t *SgTree) insert(node *Node) (*Node, int) {

	// if the tree is empty, just define the root..
	if t.Root == nil {
		t.Root = node
		t.Size++
		t.maxSize++
		return nil, 0
	}

	// traverse the tree to find the right spot
	depth := 0
	cur := t.Root
	prev := t.Root
	for cur != nil {

		// if data already exists, just return without hesitation...
		if cur.Data == node.Data {
			return cur, depth
		}

		prev = cur
		if node.Data < cur.Data {
			cur = cur.left
			//	depth++
		} else {
			cur = cur.right
			//	depth++
		}
		depth++
	}

	// insert the new element
	switch {

	case node.Data < prev.Data:
		prev.left = node
		t.Size++
		t.maxSize++
		node.parent = prev

	case node.Data > prev.Data:
		prev.right = node
		t.Size++
		t.maxSize++
		node.parent = prev
	}

	if t.Size > t.maxSize {
		t.maxSize = t.Size
	}
	return node, depth
}

// Add adds a new element to the tree.
func (t *SgTree) Add(data int) { t.Insert(NewNode(data)) }

// Insert a new node into the tree. Insertion is done iteratively, not using recursion (as usual).
// We try to avoid the problem with recursion depth when tree grows really large.
func (t *SgTree) Insert(node *Node) {

	// insert new node operation, the same as generic BST; also get depth for the node
	node, depth := t.insert(node)
	// if node is empty, no need to do anything...
	// but if not, we check the alpha-height-balance value of the tree; if height of the tree exceeds
	// the alpha-height-balance value, find a scapegoat node and rebuild a subtree rooted by scapegoat
	if node != nil {

		if float64(depth) > t.heightFactor(t.Size) {

			fmt.Printf("DEBUG alpha-height-factor=%f\n", t.heightFactor(t.Size))
			// let's find the scapegoat and rebuild the tree around it
			sg := t.findScapegoat(node.parent)
			fmt.Printf("DEBUG found scapegoat=%d\n", sg)
			if sg != nil {
				t.rebuild(sg)
			}
		}
	}
}

// Find the scapegoat node.
func (t *SgTree) findScapegoat(node *Node) *Node {

	var sibling *Node
	h := 0 // current height
	csize, totalsize := 1, 0
	cur := node

	for cur.parent != nil {

		h++
		// we define the current node's sibling
		if cur == cur.parent.left {
			sibling = cur.parent.right
		} else {
			sibling = cur.parent.left
		}
		totalsize = 1 + csize + size(sibling)

		if float64(h) > (math.Log10(float64(totalsize)) / math.Log10(1/t.alpha)) {
			return cur.parent
		}
		cur = cur.parent
		csize = totalsize
	}
	return cur
}

// Rebuild the tree after scape-goat has been found (root is here the root of the given subtree!).
func (t *SgTree) rebuild(root *Node) *Node {

	// flatten the subtree into linked list
	t.tree2List(root)
	// rebuild the linked list into balanced BST
	//root = t.list2Tree(root)
	t.list2Tree(root)
	// since we have parent ptrs, we need to do some book-keeping
	t.updateParents(root)

	return root
}

// Tree2Vine converts a BST into a vine (sorted linked list) using left pointers.
func (t *SgTree) tree2List(root *Node) {

	var cur, prev = root, root
	var temp *Node

	for cur != nil {

		if cur.right == nil {
			// if there's no right child, we don't need to do anything
			prev = cur
			cur = cur.left
		} else {
			// otherwise we need to make left rotation: right child is inserted between current and previous (parent) node
			temp = cur.right
			cur.right = temp.left
			temp.left = cur
			cur = temp
			prev.left = temp
		}
	}
}

// The list2Tree method creates the balanced tree from the linked list...
//func (t *SgTree) list2Tree(root *Node) *Node {
func (t *SgTree) list2Tree(root *Node) {

	s := size(root)
	// calculate the number of leaves in the bottom level of the balanced tree
	// note: function is defined in <dslwalgo.go>
	leaves := numOfLeaves(s)

	// now do the compression
	// the first compression iteration is to reduce the compression to general case; only when leaves > 0 (!)
	if leaves > 0 {
		t.compress(root, leaves)
	}
	fmt.Printf("DEBUG list2Tree() special to general case finished. \n")
	t.Traverse() // DEBUG

	s = s - leaves // number of nodes in main vine
	cur := root
	for s > 1 {
		fmt.Printf("DEBUG list2Tree() vine=%d. \n", s)
		t.compress(cur, (s / 2))
		s /= 2
	}
	fmt.Printf("DEBUG list2Tree() finished. \n")
	t.Traverse() // DEBUG
	//	return root
}

/*
// This is a rotate-right operation method.
func (t *SgTree) rotateR(node *Node) *Node {

	left := node.left
	// rotate
	node.left = left.right
	left.right = node
	left.parent = node.parent
	node.parent = left

	// new root of the subtree
	return left
}
*/

// Vine-to-balanced-tree compress helper function.
func (t *SgTree) compress(root *Node, count int) {

	if root == nil {
		return
	}

	var cur, child *Node

	cur = root
	for ; count != 0; count-- {
		child = cur.left
		cur.left = child.left
		cur = cur.right
		child.left = cur.left
		cur.right = child
	}
}

// Update the parent pointers after the tree's been rebalanced.
func (t *SgTree) updateParents(root *Node) {

	var prev, next *Node
	cur := root

	root.parent = nil // make sure root's parent does not point anywhere...
	for cur != nil {

		switch {
		case prev == cur.parent: // we are in parent node, we try to go left, then right, then back to parent
			if cur.left != nil {
				next = cur.left
				next.parent = cur
			} else if cur.right != nil {
				next = cur.right
				next.parent = cur
			} else {
				next = cur.parent
			}
		case prev == cur.left: // we are in left element: we try to go right, then back to parent
			if cur.right != nil {
				next = cur.right
				next.parent = cur
			} else {
				next = cur.parent
			}
		default: // we are in right element, go back to parent
			next = cur.parent
		}
		prev = cur
		cur = next
	}
}

/*
// Post-order traversal of the tree.
func (t *SgTree) PostOrder() {
	fmt.Print("Post-order: ")
	t.postorder(t.Root)
	fmt.Println()
}

// Post-order traversal
func (t *SgTree) postorder(node *Node) {

	if node == nil {
		return
	}
	t.postorder(node.left)
	t.postorder(node.right)
	fmt.Printf("%d ", node.Data)
}
*/

// PreOrder traverses the tree in pre-order fashion.
func (t *SgTree) PreOrder() {
	fmt.Print(" Pre-order: ")
	t.preorder(t.Root)
	fmt.Println()
}

// Pre-order traversal
func (t *SgTree) preorder(node *Node) {

	if node == nil {
		return
	}
	fmt.Printf("%d ", node.Data)
	t.preorder(node.left)
	t.preorder(node.right)
}

// InOrder traverses the tree in in-order fashion.
func (t *SgTree) InOrder() {
	fmt.Print("  In-order: ")
	t.inorder(t.Root)
	fmt.Println()
}

// traversing left
func (t *SgTree) inorder(node *Node) {

	if node == nil {
		return
	}
	t.inorder(node.left)
	fmt.Printf("%d ", node.Data)
	t.inorder(node.right)
}

// Searches for the MIN element of the tree; it's the far left element.
func (t *SgTree) findMinElem(node *Node) (*Node, error) {

	cur := node // we start at node
	// if tree is still empty, just return an error
	if cur == nil {
		return cur, fmt.Errorf("Empty tree")
	}

	for cur.left != nil {
		cur = cur.left
	}
	return cur, nil
}

// Min returns the Min() element in the tree. This element is the far left.
func (t *SgTree) Min() (int, error) {

	cur, err := t.findMinElem(t.Root)
	if err != nil {
		return 0, nil
	}
	return cur.Data, err
}

// Searches for the MAX element of the (sub)tree; it's the far right element.
func (t *SgTree) findMaxElem(node *Node) (*Node, error) {

	cur := node
	// if tree is still empty, just return an error
	if cur == nil {
		return nil, fmt.Errorf("Empty tree")
	}

	for cur.right != nil {
		cur = cur.right
	}
	return cur, nil

}

// Max returns the Max() element in the tree. This element is far right.
func (t *SgTree) Max() (int, error) {

	cur, err := t.findMaxElem(t.Root)
	if err != nil {
		return 0, err
	}
	return cur.Data, nil
}

// Delete deletes an element from the tree.
func (t *SgTree) Delete(node *Node) {

	var elem *Node
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
		t.Size-- // we have one element less...
	}

	// Check if rebuild of tree is needed
	if elem != nil {
		if float64(t.Size) < (float64(t.maxSize) * t.alpha) {
			t.rebuild(t.Root)
			t.maxSize = t.Size
		}
	}
	// NOTE: if no element is found (elem == nil), just return; works for empty root, too
	return
}

// Remove removes the node with the given value from the tree.
func (t *SgTree) Remove(data int) {
	n := NewNode(data)
	t.Delete(n)
}

/*
// goroutine used as a dispatcher for operations on tree.
func (t *SgTree) do() {

	var node *Node

	select {

	case node = <-t.insch:
		//node = <-insch
		t.Insert(node)

	case node = <-t.rmch:
		//node = <-rmch
		t.Delete(node)

		// case node = <-t.fch:
		//     foundch <- t.find(t.Root,  )

	case <-t.quitch:
		close(t.insch)
		close(t.rmch)
		close(t.fch)
		close(t.foundch)
		close(t.quitch)
		return
	}
}

// Stop sends a signal on the 'quit' channel that we're done.
func (t *SgTree) Stop() { t.quitch <- 1 }
*/
// Traverse implements the tree traversing left iteratively
// FIXME: there's obviously a bug somewhere in this method...
func (t *SgTree) Traverse() {

	var cur = t.Root
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
