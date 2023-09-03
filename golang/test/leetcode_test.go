package test

import (
	"log"
	"strings"
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

func containsDuplicate(nums []int) bool {
	k:=0
	for _,v:=range nums {
		if v == 0 {
			continue
		}
		if k ==0 {
			k = v
			continue
		}
		k = k^v
		if k == 0 {
			return true
		}
	}
	return false
}

func romanToInt(s string) int {
	m:=map[rune]int{
		'I':1,
		'V':5,
		'X':10,
		'L':50,
		'C':100,
		'D':500,
		'M':1000,
	}
	sum:=0
	var pre int
	for _,r:=range s {
		v:=m[r]
		sum = sum+v
		if pre < v {
			sum = sum - 2*pre
		}
		pre = v
		// switch(pre){
		//     case 'I':
		//         sum = sum - 1
		//     case 'X':
		//         sum = sum - 10
		//     case 'C':
		//         sum = sum - 100
		// }

	}
	return sum
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 1 {
		return strs[0]
	}
	pre:=strs[0]
	strs = strs[1:]

	for _,str:=range strs {
		if pre == "" {
			return ""
		}

		i,j:=0,0
		n1,n2:=len(pre), len(str)
		for  {
			if i >=n1 || j >=n2 {
				pre = pre[0:i]
				break
			}
			if pre[i] != str[j] {
				pre = pre[0:i]
				break
			}
			i++
			j++
		}
	}
	return pre
}

func lengthOfLastWord(s string) int {
	strs:=strings.Split(s," ")
	i:=len(strs)-1

	for i>=0 {
		if strs[i] != "" {
			return len(strs[i])
		}
		i--
	}
	return 0
}

func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	size:=len(s)
	if size <= 1 {
		return true
	}
	if size == 2 {
		if s[0] == s[1] {
			return true
		}

		return !isEnByte(s[0]) || !isEnByte(s[1])
	}
	i,j:=0,size-1
	for i<j{
		if s[i]!=s[j] {
			if !isEnByte(s[i]) {
				i++
			} else if !isEnByte(s[j]) {
				j--
			} else {
				return false
			}
		} else {
			i++
			j--
		}
	}
	return true
}
func isEnByte(by byte) bool {
	return by>='a' && by<='z' || by >= '0' && by <='9'
}

func longestPalindrome(s string) string {
	start,end:=0,0
	l,r:=0,0
	n:=len(s)
	for i:=1;i<n-1;i++ {
		l = i-1
		r = i+1
		for {
			if l < 0 || r >= n{
				l = l+1
				r = r-1
				break
			}
			if s[l] != s[r] {
				l = l+1
				r = r-1
				break
			}
			l--
			r++
		}

		if r-l == 0 && end-start == 0 && s[i] == s[i+1] {
			start,end = i, i+1
			continue
		}
		if (r-l) > (end-start) {
			start,end = l,r
		}
	}
	return s[start:end+1]
}
func aa(s string ) [26]int  {
	var cnt [26]int
	for _,ch:=range s  {
		cnt[ch-'a']++
	}
	log.Print(cnt)
	return cnt
}

//-2,1,-3,4,-1,2,1,-5,4
func maxSubArray(nums []int) int {
	max:=nums[0]
	for i:=1;i<len(nums);i++{
		if nums[i] + nums[i-1] > nums[i] {
			nums[i] = nums[i] + nums[i-1]
		}
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}


func step(n int) int {
	sum := 0
	for n > 0 {
		sum += (n%10) * (n%10)
		n = n/10
	}
	return sum
}

func isHappy(n int) bool {
	slow,fast:=step(n),step(step(n))
	for {
		if slow == fast {
			return false
		}
		if fast == 1 {
			return true
		}
		slow, fast = step(slow), step(step(fast))
	}
}

func TestLeetCode(t *testing.T) {
	//nums:=[]int{1,3}
	//index:=searchInsert(nums, 2)
	//t.Log(index)

	//nums:=[]int{1,2,3,4}
	//b:=containsDuplicate(nums)
	//t.Log(b)

	//s:=romanToInt("MCMXCIV")
	//t.Log(s)

	//b:=isValid("()")
	//t.Log(b)

	//s1:="   fly me   to   the moon  "
	//size:=lengthOfLastWord(s1)
	//t.Log(size)

	//s:="0P"
	////s:="a,,a"
	//b:=isPalindrome(s)
	//t.Log(b)

	//str:="bb"
	//s:=longestPalindrome(str)
	//t.Log(s)


	//cnt1:=aa("aaabcz")
	//cnt2:=aa("aaabcz")
	//t.Log(cnt1==cnt2)

	//nums:=[]int{-2,1,-3,4,-1,2,1,-5,4}
	//max:=maxSubArray(nums)
	//t.Log(max)

	s:=isHappy(2)
	t.Log(s)

	sum := step(1)
	t.Log(sum)




}
