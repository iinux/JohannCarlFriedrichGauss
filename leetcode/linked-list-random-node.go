package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	// 初始化一个单链表 [1,2,3].
	head := ListNode{Val:1};
	head.Next = &ListNode{Val:2};
	head.Next.Next = &ListNode{Val:3};
	solution := Constructor(&head);

	// getRandom()方法应随机返回1,2,3中的一个，保证每个元素被返回的概率相等。
	fmt.Println(solution.GetRandom());
}

/**
 * Definition for singly-linked list.
 */

type ListNode struct {
    Val int
    Next *ListNode
}

type Solution struct {
	list *ListNode

}


/** @param head The linked list's head.
        Note that the head is guaranteed to be not null, so it contains at least one node. */
func Constructor(head *ListNode) Solution {
	return Solution{list:head}
}


/** Returns a random node's value. */
func (this *Solution) GetRandom() int {
	var k int64
	// int8 is right?
	k = math.MaxInt8
	var r int64
	var list *ListNode
	for true {
		if list == nil {
			list = this.list
		}
		r = rand.Int63n(k)
		if r == 0 {
			break
		} else {
			k--
			list = list.Next
		}
	}
	return list.Val
}


/**
 * Your Solution object will be instantiated and called as such:
 * obj := Constructor(head);
 * param_1 := obj.GetRandom();
 */
