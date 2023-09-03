package test

import (
	"math"
	"testing"
)

var dfMap map[int]float64 = map[int]float64{}
func recursion(i int, nums[]int) float64 {
	if l,ok:=dfMap[i];ok {
		return l
	}
	//7,7,7,7
	var maxLen float64 = 1
	for j:=i+1;j<len(nums);j++ {
		if nums[j] > nums[i] {
			maxLen = math.Max(maxLen,recursion(j, nums)+1)
		}
	}
	dfMap[i] = maxLen
	return maxLen
}

func lengthOfLIS(nums []int) int {
	maxLen:=0
	for i:=0;i<len(nums);i++ {
		l:=int(recursion(i,nums))
		if l > maxLen{
			maxLen= l
		}
	}
	return maxLen
}

func TestLongest(t *testing.T) {
	nums:=[]int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30}
	//nums:=[]int{1,3,6,7,9,4,10,5,6}
	//nums:=[]int{1,5,2,4,3}
	l:=lengthOfLIS(nums)
	t.Log(l)

}
