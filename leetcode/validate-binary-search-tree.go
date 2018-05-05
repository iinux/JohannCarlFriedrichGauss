package main

import "fmt"

func main() {
	root := &TreeNode{}
	root.Val = 1
	root.Left = &TreeNode{}
	root.Left.Val = 1
	fmt.Println(isValidBST(root))
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

var nums []int

func isValidBST(root *TreeNode) bool {
	nums = []int{}

	traversal(root)

	l := len(nums)
	for i := 1; i < l; i++ {
		if nums[i-1] >= nums[i] {
			return false
		}
	}
	return true
}

func traversal(root *TreeNode) {
	if root == nil {
		return
	}
	traversal(root.Left)
	nums = append(nums, root.Val)
	traversal(root.Right)
}

