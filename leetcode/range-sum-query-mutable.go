package main

type NumArray struct {
	Data []int
}

func Constructor(nums []int) NumArray {
	return NumArray{Data:nums}
}

func (this *NumArray) Update(i int, val int) {
	this.Data[i] = val
}

func (this *NumArray) SumRange(i int, j int) int {
	sum := 0
	for ; i <= j; i++ {
		sum += this.Data[i]
	}
	return sum
}


/**
 * Your NumArray object will be instantiated and called as such:
 * obj := Constructor(nums);
 * obj.Update(i,val);
 * param_2 := obj.SumRange(i,j);
 */
