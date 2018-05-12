package main

func main() {

}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func isSubtree(s *TreeNode, t *TreeNode) bool {
	if s == nil  && t == nil {
		return true
	}
	if s == nil && t != nil {
		return false
	}

	if s.Val == t.Val {
		c := compareTree(s, t)
		if c {
			return c
		}
	}

	return isSubtree(s.Left, t) || isSubtree(s.Right, t)
}

func compareTree(s *TreeNode, t *TreeNode) bool {
	if s == nil && t == nil {
		return true
	}
	if s == nil || t == nil {
		return false
	}
	return s.Val == t.Val && compareTree(s.Left, t.Left) && compareTree(s.Right, t.Right)
}
