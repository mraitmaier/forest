package forest

// This is implementation of the scapegoat tree in Go.
// 
// FIXME: not finished yet!

import (
	"fmt"
	"math"
)

// Define the ScapeGoat Tree node
type SgNode struct {

	// Data
	Data int

	// left, right and parent node
	left, right, parent *SgNode
}

// Create a new tree node.
func NewSgNode(data int) *SgNode { return &SgNode{data, nil, nil, nil} }

// A string representation of the node.
func (n *SgNode) String() string {

	s := fmt.Sprintf("Node: %d ", n.Data)
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
		s = fmt.Sprintf("%s parent=&v ", s, n.parent)
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

// Define the Scape Goat Tree
type SgTree struct {

	// The root of the tree
	Root *SgNode

	// Number of nodes in the tree
	Size int

	// alpha factor for the tree: must be between 0.5 and 1.0; this is user-definable parameter;
	// we set this factor to value of 'defAlpha' constant (see above) by default
	alpha float64

	// Maximum size of the tree before it is rebuild
	maxSize int

    // channels defined 
    insch, rmch, fch, foundch chan *SgNode
    quitch chan int
}

// Create a new empty tree.
//func NewSgTree() *SgTree { return &SgTree{nil, 0, defAlpha, 0} }
func NewSgTree() *SgTree {
//    return &SgTree{nil, 0, defAlpha, 0,
//                                  make(chan *SgNode), make(chan *SgNode), make(chan *SgNode), make(chan *SgNode), make(chan int) }
    n :=  &SgTree{nil, 0, defAlpha, 0,
                                  make(chan *SgNode), make(chan *SgNode), make(chan *SgNode), make(chan *SgNode), make(chan int) }
    go n.do() // start a dispather goroutine
    return n
}

// Set the value of Alpha factor.
func (t *SgTree) SetAlpha(val float64) error {

	// check Alpha value, must be between 0.5 and 1.0
	if val < 0.5 || val > 1.0 {
		return fmt.Errorf("Alpha value must be between 0.5 and 1.0")
	}
	t.alpha = val
	return nil
}

// Return the value of Alpha factor.
func (t *SgTree) GetAlpha() float64 { return t.alpha }

// Calculate tree height.
func (t *SgTree) height(count int) float64 {
	return math.Log10(float64(count)) / math.Log10(1/t.alpha)
}

// Check if tree is empty.
func (t *SgTree) IsEmpty() bool {

	if t.Root == nil {
		return true
	}
	return false
}

// Clear (destroy) the complete tree. The result is empty tree.
func (t *SgTree) Clear() {
	t.Root = nil
	t.Size = 0
    t.maxSize = 0
}

// Return a node to be found. If not found, return nil.
func (t *SgTree) find(root *SgNode, data int) *SgNode {

	var cur *SgNode = root
	for cur != nil  {
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

// Search for a given data in the tree
func (t *SgTree) SearchFor(data int) *SgNode { return t.find(t.Root, data) }

// Check if given element is a member of a tree.
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
		t.Size += 1
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
			depth += 1
		} else {
			cur = cur.right
			depth += 1
		}
	}

	// insert the new element
	switch {

	case node.Data < prev.Data:
		prev.left = node
		t.Size += 1
		node.parent = prev

	case node.Data > prev.Data:
		prev.right = node
		t.Size += 1
		node.parent = prev

	default: // do nothing when element already exists
	}

	if t.Size > t.maxSize {
		t.maxSize = t.Size
	}
	return node, depth
}

// Insert a new node into the tree. Insertion is done iteratively, not using recursion (as usual).
// We try to avoid the problem with recursion depth, when tree grows really large.
func (t *SgTree) Insert(node *SgNode) {

	var d int
	// insert new node operation, the same as generic BST
	node, d = t.insert(node)

	// if node is empty, nothing happened during insert(), no need to rebuild;
	// but if not, we check the alpha-height-balance value and rebuild the tree, if needed;
	if node != nil {
		fmt.Printf("DEBUG Insert(): depth=%d\n", d)

        //if float64(d) > (t.height(size(node)) + 1) {
        if float64(d) > t.height(size(node)) {
            fmt.Printf("DEBUG alpha-height-factor=%f\n", t.height(t.Size))

            // let's find the scapegoat and rebuild the tree around it
            sg := t.findScapegoat(node)
            parent := sg.parent
            sg = t.rebuild(sg)
            sg.parent = parent
        }
	}
}

// Add a new element to the tree.
func (t *SgTree) Add(data int) {
	n := NewSgNode(data)
	t.Insert(n)
}

// Find the scapegoat node.
func (t *SgTree) findScapegoat(node *SgNode) *SgNode {

    fmt.Printf("DEBUG findScapegoat(): start...\n") // DEBUG
    var sibling *SgNode

    csize := 1
    var totalsize, sibsize int
    cur := node.parent
    for cur != nil  {

        //if cur.parent == nil { // we are at root...
        //    return nil
        //}

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

// 
func (t *SgTree) rebuild(root *SgNode) *SgNode{

    // empty node as scapegoat? return...
    if root == nil {
        return nil
    }

    // flatten the subtree
//    size := 0

    // now rebuild around new root

    // TODO
    return nil
}

//
func (t *SgTree) flatten(root *SgNode) *SgNode {

    if root == nil {
        return root
    }

    return nil
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

// Pre-order traversal of the tree.
func (t *SgTree) PreOrder() {
	fmt.Print(" Pre-order: ")
	t.preorder(t.Root)
	fmt.Println()
}

// Pre-order traversal
func (bt *SgTree) preorder(node *SgNode) {

	if node == nil {
		return
	}
	fmt.Printf("%d ", node.Data)
	bt.preorder(node.left)
	bt.preorder(node.right)
}

// In-order traversal of the tree.
func (t *SgTree) InOrder() {
	fmt.Print("  In-order: ")
	t.inorder(t.Root)
	fmt.Println()
}

// traversing left iteratively
func (bt *SgTree) inorder(node *SgNode) {

	if node == nil {
		return
	}
	bt.inorder(node.left)
	fmt.Printf("%d ", node.Data)
	bt.inorder(node.right)
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

	var cur *SgNode = node // we start at node

	// if tree is still empty, just return an error
	if cur == nil {
		return cur, fmt.Errorf("Empty tree")
	}

	for cur.left != nil {
		cur = cur.left
	}
	return cur, nil
}

// Return the Min() element in the tree. This element is the far left.
func (t *SgTree) Min() (int, error) {

	if cur, err := t.findMinElem(t.Root); err != nil {
		return 0, nil
	} else {
		return cur.Data, err
	}
}

// Searches for the MAX element of the (sub)tree; it's the far right element.
func (t *SgTree) findMaxElem(node *SgNode) (*SgNode, error) {

	var cur *SgNode = node // we start at node (usually root)

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
func (t *SgTree) Max() (int, error) {

	if cur, err := t.findMaxElem(t.Root); err != nil {
		return 0, err
	} else {
		return cur.Data, nil
	}
}

// Delete an element from the tree.
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
			min, _ := t.findMinElem(elem.right)         // find MIN value in right subtree
			elem.Data, min.Data = min.Data, elem.Data   // exchange the element value and the MIN value of right subtree
            if min.right != nil {
                min.Data = min.right.Data               //
                min.right = nil
            } else {
                min.parent.right = nil
            }
		}
		t.Size -= 1 // we have one element less...
	}

    // 
    if elem != nil {
        if float64(t.Size) < float64(t.maxSize) * t.alpha {
            t.Root = t.rebuild(t.Root)
            t.maxSize = t.Size
        }
    }
	// NOTE: if no element is found (elem == nil), just return; works for empty root, too
	return
}

// Remove the node with the given value.
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

func(t *SgTree) Stop() { t.quitch <- 1 }

