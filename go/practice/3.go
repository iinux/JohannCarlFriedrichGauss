package main

import (
	"fmt"
)

type user struct{ name string }

func main() {
	m := map[int]user{ // 当 map 因扩张而重新哈希时，各键值项存储位置都会发生改变。 因此，map
		1: {"user1"}, // 被设计成 not addressable。 类似 m[1].name 这种期望透过原 value
	} // 指针修改成员的行为自然会被禁止。

	fmt.Println(m[1].name)
	//m[1].name = "Tom" // Error: cannot assign to m[1].name
}
