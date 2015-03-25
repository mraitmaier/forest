package forest

//
// dswalgo.go - an implementation of the Day-Stout-Warren algorithm (DSW), a method of rebalancing the ordinary Binary Search Trees
//
// see more at: http://en.wikipedia.org/wiki/Day-Stout-Warren_algorithm

import (
//    "fmt"
)

// Balance rebalances the BST.
func Balance(tree *BinTree) {

	if tree.Root == nil {
		return
	}

	root := Tree2Vine(tree.Root)
	//    traverseVine(root) // DEBUG
	tree.Root = Vine2Tree(root, tree.Len)
	updateParents(tree)
}

// Tree2Vine converts a BST into a vine (sorted linked list) using left pointers.
func Tree2Vine(root *Node) *Node {

	var prev *Node
	var temp *Node
	var cur = root

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
			if prev != nil { // prev can be empty, when prev is root...
				prev.left = temp
			}

			if cur.Data > root.Data { // define new vine root
				root = cur
			}
		}
	}
	return root
}

// Calculates the number of leaves in the bottom level of the balanced tree
func numOfLeaves(size int) int {

	leaves, next := size+1, 0
	for {
		next = leaves & (leaves - 1)
		if next == 0 {
			break
		}
		leaves = next
	}
	return size + 1 - leaves
}

// Vine2Tree transforms the given vine back into a balanced tree.
func Vine2Tree(root *Node, size int) *Node {

	// calculate the number of leaves in the bottom level of the balanced tree
	leaves := numOfLeaves(size)

	// now do the compression
	// the first compression iteration is to reduce the compression to general case; only when leaves > 0 (!)
	if leaves > 0 {
		root = compress(root, leaves)
	}
	vine := size - leaves // number of nodes in main vine

	for vine > 1 {
		vine /= 2
		root = compress(root, int(vine))
	}
	return root
}

// Vine-to-balanced-tree compress helper function.
func compress(root *Node, count int) *Node {

	red := root
	black := red.left

	root = black // new root
	root.parent = nil

	for ; count != 0; count-- {
		red.left = black.right
		black.right = red
		red = black.left
		if count != 1 { // the last count this step must be omitted; otherwise we lose an element...
			black.left = red.left
		}
		black = red.left
	}

	return root
}

// Update the parent pointers after the tree's been rebalanced.
func updateParents(bt *BinTree) {

	var cur = bt.Root
	var prev *Node
	var next *Node

	cur.parent = nil // make sure root's parent does not point anywhere...
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
func traverseVine(root *Node) {
    fmt.Print("Traversing vine: ")
    for cur := root; cur != nil; cur = cur.left {
        fmt.Printf("%d ", cur.Data)
    }
    fmt.Println()
}
*/
