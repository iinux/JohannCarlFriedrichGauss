package main

func main() {
	printListNode(mergeTwoLists(makeListNode([]int{1,2,4}), makeListNode([]int{1,3,4})))
	printListNode(mergeTwoLists(makeListNode([]int{}), makeListNode([]int{})))
	printListNode(mergeTwoLists(makeListNode([]int{}), makeListNode([]int{0})))
}

func makeListNode(ns []int) *ListNode {
	l := len(ns)
	if l == 0 {
		return nil
	}

	list := &ListNode{Val: ns[0]}
	p := list
	for i := 1; i < l; i++ {
		p.Next = &ListNode{Val: ns[i]}
		p = p.Next
	}

	return list
}

func printListNode(list *ListNode)  {
	p := list
	for p != nil {
		print(p.Val, " ")
		p = p.Next
	}
	println()
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil && l2 == nil {
		return nil
	} else if l1 == nil && l2 != nil {
		return l2
	} else if l1 != nil && l2 == nil {
		return l1
	}

	p := l1
	p2 := l2

	if l1.Val > l2.Val {
		p = l2
		p2 = l1
	}
	r := p

	for true {
		if p == nil || p2 == nil {
			break
		}

		if p.Next == nil {
			p.Next = p2
			break
		}

		if p.Next.Val > p2.Val {
			t := p2
			p2 = p2.Next
			t.Next = p.Next
			p.Next = t
		} else {
			p = p.Next
		}
	}

	return r
}
