package main

import (
	"strings"
	"fmt"
)

func main() {
	fmt.Println(isValidSerialization("9,3,4,#,#,1,#,#,2,#,6,#,#"))
	fmt.Println(isValidSerialization("1,#"))
	fmt.Println(isValidSerialization("9,#,#,1"))
	fmt.Println(isValidSerialization("#,#"))
	fmt.Println(isValidSerialization("1"))
	fmt.Println(isValidSerialization("1,#,#,#,#"))
}

type node struct {
	v       string
	isLeaf  bool
	deleted bool
}

func isValidSerialization(preorder string) bool {
	ps := strings.Split(preorder, ",")
	var ns []*node
	for _, s := range ps {
		if s == "#" {
			ns = append(ns, &node{v:s, isLeaf:true, deleted:false})
		} else {
			ns = append(ns, &node{v:s, deleted:false})
		}
	}
	for true {
		do := 0
		var collection []*node
		for _, v := range ns {
			if !v.deleted {
				collection = append(collection, v)
			}
		}

		l := len(collection)
		for i := 1; i < l - 1; i++ {
			if !collection[i].deleted && !collection[i + 1].deleted &&
				collection[i].isLeaf && collection[i + 1].isLeaf &&
				!collection[i - 1].isLeaf {
				collection[i - 1].isLeaf = true
				collection[i].deleted = true
				collection[i + 1].deleted = true
				do++
			}
		}
		if do == 0 {
			if l == 1 && collection[0].isLeaf {
				return true
			} else {
				return false
			}
		}
	}

	return false
}
