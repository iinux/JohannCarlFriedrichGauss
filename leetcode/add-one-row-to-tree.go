package main

func main() {

}

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func addOneRow(root *TreeNode, v int, d int) *TreeNode {
	if root == nil {
		return root
	}
	if d == 1 {
		t := TreeNode{
			Val:v,
			Left:root,
		}
		return &t
	} else if d == 2 {
		t1 := TreeNode{
			Val:v,
			Left:root.Left,
		}
		t2 := TreeNode{
			Val:v,
			Right:root.Right,
		}
		root.Left = &t1
		root.Right = &t2
	} else {
		addOneRow(root.Left, v, d - 1)
		addOneRow(root.Right, v, d - 1)
	}
	return root
}
