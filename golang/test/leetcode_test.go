package test

import (
	"testing"
)

func searchInsert(nums []int, target int) int {
	mid:=len(nums)/2
	left:=0
	right:=len(nums)
	for right > left {
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			right = mid-1
		}else {
			left = mid+1
		}

		mid = (right+left)/2
	}

	return mid
}

func TestLeetCode(t *testing.T) {
	nums:=[]int{1,3}
	index:=searchInsert(nums, 2)
	t.Log(index)
}
