package test

import (
	"testing"
)

type NumArray struct {
	sums []int
}


func Constructor(nums []int) NumArray {
	sums:=make([]int, len(nums))
	sums[0] = nums[0]
	for i:=1;i<len(nums);i++ {
		sums[i] = sums[i-1] + nums[i]
	}
	return NumArray{sums}
}


func (this *NumArray) SumRange(left int, right int) int {
	if left == 0 {
		return this.sums[right]
	}
	return this.sums[right]-this.sums[left-1]
	//return this.sums[right+1] - this.sums[left]
}


/**
 * Your NumArray object will be instantiated and called as such:
 * obj := Constructor(nums);
 * param_1 := obj.SumRange(left,right);
 */
func TestPreSum(t *testing.T) {
	nums:=[]int{-2, 0, 3, -5, 2, -1}
	c:=Constructor(nums)
	sum:=c.SumRange(2,5)
	t.Log(sum)
}
