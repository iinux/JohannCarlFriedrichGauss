package main

import (
	"fmt"
	"strconv"
)

func main2()  {
	nums := []int{1,1,2}
	println(removeDuplicates(nums))
	printNums(nums)
}

func main() {
	str := "1"
	fmt.Println(strconv.Atoi(str))	// int error
	fmt.Println(strconv.ParseBool("false"))
	fmt.Println(strconv.Itoa(1))
	fmt.Println([]byte("test"))// 字符串转byte数组
	// 2,byte转为string
	byte1 := []byte{116,101,115,116}
	fmt.Println(string(byte1[:]))
}

func main1() {
	var numbers []int
	printSlice(numbers)

	/* 允许追加空切片 */
	numbers = append(numbers, 0)
	printSlice(numbers)

	/* 向切片添加一个元素 */
	numbers = append(numbers, 1)
	printSlice(numbers)

	/* 同时添加多个元素 */
	numbers = append(numbers, 2,3,4)
	printSlice(numbers)

	/* 创建切片 numbers1 是之前切片的两倍容量*/
	numbers1 := make([]int, len(numbers), (cap(numbers))*2)

	/* 拷贝 numbers 的内容到 numbers1 */
	copy(numbers1,numbers)
	printSlice(numbers1)
}

func printSlice(x []int){
	fmt.Printf("len=%d cap=%d slice=%v\n",len(x),cap(x),x)
}

func printNums(nums []int)  {
	l := len(nums)
	for i := 0; i < l; i++ {
		print(nums[i], " ")
	}
	println()
}

func removeDuplicates(nums []int) int {
	l := len(nums)

	if l == 0 {
		return 0
	}

	p1 := 1
	p2 := 1
	for ; p1 < l; p1++ {
		if nums[p1] != nums[p1-1] {
			nums[p2] = nums[p1]
			p2++
		}
	}

	nums = nums[:p2]
	return p2
}
