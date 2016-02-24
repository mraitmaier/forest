package forest

// This is implementation of the scapegoat tree in Go.
//
// FIXME: not finished yet!

import (
	"fmt"
	"math"
)

// SgNode defines the ScapeGoat Tree node.
type SgNode struct {

	// Data
	Data int

	// left, right and parent node
	left, right, parent *SgNode
}

// NewSgNode creates a new tree node.
func NewSgNode(data int) *SgNode { return &SgNode{data, nil, nil, nil} }

// A string representation of the node.
func (n *SgNode) String() string {

	s := fmt.Sprintf("Node: %d ", n.Data)
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
		s = fmt.Sprintf("%s parent=%v ", s, n.parent)
	} else {
		s = fmt.Sprintf("%s Root", s)
	}
	return s
}

// Node size
func size(node *SgNode) int {

	if node == nil {
		return 0
	}
	return size(node.left) + size(node.right) + 1
}

// The default value of Alpha factor
const defAlpha = 0.6

// SgTree defines the Scapegoat Tree structure.
type SgTree struct {

	// Root pointd to the root of the tree
	Root *SgNode

	// Size represents the number of nodes in the tree
	Size int

	// alpha factor for the tree: must be between 0.5 and 1.0; this is user-definable parameter;
	// we set this factor to value of 'defAlpha' constant (see above) by default
	alpha float64

	// Maximum size of the tree before it is rebuild
	maxSize int

	// channels defined
	insch, rmch, fch, foundch chan *SgNode
	quitch                    chan int
}

// NewSgTree creates a new empty scape-goat tree.
//func NewSgTree() *SgTree { return &SgTree{nil, 0, defAlpha, 0} }
func NewSgTree() *SgTree {
	//    return &SgTree{nil, 0, defAlpha, 0,
	//                         make(chan *SgNode), make(chan *SgNode), make(chan *SgNode), make(chan *SgNode), make(chan int) }
	n := &SgTree{
		Root:    nil,
		Size:    0,
		alpha:   defAlpha,
		maxSize: 0,
		insch:   make(chan *SgNode),
		rmch:    make(chan *SgNode),
		fch:     make(chan *SgNode),
		foundch: make(chan *SgNode),
		quitch:  make(chan int),
	}
	go n.do() // start a dispather goroutine
	return n
}

// SetAlpha changes the value of Alpha factor.
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

// Calculate tree height.
func (t *SgTree) height(count int) float64 {
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
func (t *SgTree) find(root *SgNode, data int) *SgNode {

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
func (t *SgTree) SearchFor(data int) *SgNode { return t.find(t.Root, data) }

// In checks if given element is a member of a tree.
func (t *SgTree) In(data int) bool {

	if t.SearchFor(data) != nil {
		return true
	}
	return false
}

// Insert a new node into the tree. Basically the same as the generic BST tree insert. Returns (pointer to) the inserted node and
// the depth of the node. Tree rebuild is done in its own method.
func (t *SgTree) insert(node *SgNode) (*SgNode, int) {

	// if the tree is empty, just define the root..
	if t.Root == nil {
		t.Root = node
		t.Size++
		return node, 0
	}

	// traverse the tree to find the right spot
	depth := 0
	cur := t.Root
	prev := t.Root
	for cur != nil {

		// if data already exists, just return without hesitation...
		if cur.Data == node.Data {
			return nil, 0
		}

		prev = cur
		if node.Data < cur.Data {
			cur = cur.left
			depth++
		} else {
			cur = cur.right
			depth++
		}
	}

	// insert the new element
	switch {

	case node.Data < prev.Data:
		prev.left = node
		t.Size++
		node.parent = prev

	case node.Data > prev.Data:
		prev.right = node
		t.Size++
		node.parent = prev

	default: // do nothing when element already exists
	}

	if t.Size > t.maxSize {
		t.maxSize = t.Size
	}
	return node, depth
}

// Add adds a new element to the tree.
func (t *SgTree) Add(data int) { t.Insert(NewSgNode(data)) }

// Insert a new node into the tree. Insertion is done iteratively, not using recursion (as usual).
// We try to avoid the problem with recursion depth when tree grows really large.
func (t *SgTree) Insert(node *SgNode) {

	// insert new node operation, the same as generic BST; get depth for the node
	var d int
	node, d = t.insert(node)

	// if node is empty, nothing happened during insert(), no need to rebuild;
	// but if not, we check the alpha-height-balance value and rebuild the tree when needed;
	if node != nil {
		fmt.Printf("DEBUG Insert(): depth=%d\n", d)

		//if float64(d) > (t.height(size(node)) + 1) {
		if float64(d) > t.height(size(node)) {
			fmt.Printf("DEBUG alpha-height-factor=%f\n", t.height(t.Size))

			// let's find the scapegoat and rebuild the tree around it
			sg := t.findScapegoat(node)
			fmt.Printf("DEBUG scapegoat=%d\n", sg)
			if sg != nil {
				t.rebuild(sg)
			}
		}
	}
}

// Find the scapegoat node.
func (t *SgTree) findScapegoat(node *SgNode) *SgNode {

	fmt.Printf("DEBUG findScapegoat(): start...\n") // DEBUG
	var sibling *SgNode

	csize := 1
	var totalsize, sibsize int
	cur := node.parent
	for cur != nil {

		if cur.parent == nil { // we are at root...
			break
		}

		// we define the current node's sibling
		if cur == cur.parent.left {
			sibling = cur.parent.right
		} else {
			sibling = cur.parent.left
		}
		sibsize = size(sibling) // sibling's size
		totalsize = 1 + csize + sibsize

		// check alpha-weight-balance violation
		alphaWeight := t.alpha * float64(totalsize)
		if float64(csize) > alphaWeight || float64(sibsize) > alphaWeight {
			return cur.parent
		}

		cur = cur.parent
		csize = totalsize
	}
	return cur
}

// Rebuild the tree after scape-goat has been found.
func (t *SgTree) rebuild(root *SgNode) {

	// empty node as scapegoat? Do nothing...
	if root == nil {
		return
	}

	//
	nodesize := float64(size(root))
	// remember the root's parent
	p := root.parent

	// flatten the subtree
	nodes := make([]*SgNode, 0, size(root))
	t.flatten(root, nodes)

	// if parent is empty (this is root node...)
	if p == nil {
		r := t.buildBalanced(nodes, 0.0, nodesize)
		r.parent = nil
	} else if p.right == root {
		p.right = t.buildBalanced(nodes, 0.0, nodesize)
		p.right.parent = p
	} else {
		p.left = t.buildBalanced(nodes, 0.0, nodesize)
		p.left.parent = p
	}

}

// Rebuild the balanced subtree after flattening.
func (t *SgTree) buildBalanced(nodes []*SgNode, start, end float64) *SgNode {

	if start >= end {
		return nil
	}

	mid := int(math.Ceil(start + (end-start)/2.0))

	fmt.Printf("DEBUG buildBalanced(): nodes=%v\n", nodes)                           // DEBUG
	fmt.Printf("DEBUG buildBalanced(): start=%f, end=%f, mid=%d\n", start, end, mid) // DEBUG

	node := NewSgNode(nodes[mid].Data)
	node.left = t.buildBalanced(nodes, start, float64(mid-1))
	node.right = t.buildBalanced(nodes, float64(mid+1), end)
	return node
}

// Flatten the tree into slice (array) during the rebuilding phase.
func (t *SgTree) flatten(node *SgNode, nodes []*SgNode) {

	if node == nil {
		return
	}
	t.flatten(node.left, nodes)
	nodes = append(nodes, node)
	t.flatten(node.right, nodes)
}

/*
// Post-order traversal of the tree.
func (t *SgTree) PostOrder() {
	fmt.Print("Post-order: ")
	t.postorder(t.Root)
	fmt.Println()
}

// Post-order traversal
func (t *SgTree) postorder(node *SgNode) {

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
func (t *SgTree) preorder(node *SgNode) {

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
func (t *SgTree) inorder(node *SgNode) {

	if node == nil {
		return
	}
	t.inorder(node.left)
	fmt.Printf("%d ", node.Data)
	t.inorder(node.right)
}

/*
// iterative implementation of in-order traversal --- XXX: somewhere, there's a bug hiding...
func (bt *SgTree) InorderIter(node *SgNode) {

	var cur *SgNode = bt.Root
	var prev *SgNode = nil
	var next *SgNode = nil

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

// Searches for the MIN element of the tree; it's the far left element.
func (t *SgTree) findMinElem(node *SgNode) (*SgNode, error) {

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
func (t *SgTree) findMaxElem(node *SgNode) (*SgNode, error) {

	cur := node // we start at node (usually root)

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
func (t *SgTree) Delete(node *SgNode) {

	var elem *SgNode
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

	//
	if elem != nil {
		if float64(t.Size) < float64(t.maxSize)*t.alpha {
			//t.Root = t.rebuild(t.Root)
			t.rebuild(t.Root)
			t.maxSize = t.Size
		}
	}
	// NOTE: if no element is found (elem == nil), just return; works for empty root, too
	return
}

// Remove removes the node with the given value from the tree.
func (t *SgTree) Remove(data int) {
	n := NewSgNode(data)
	t.Delete(n)
}

// goroutine used as a dispatcher for operations on tree.
func (t *SgTree) do() {

	var node *SgNode

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
