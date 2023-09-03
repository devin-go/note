# map
```
一般将数据存储起来，降低时间复杂度
```
## 相关题
[数组里是否存在两个元素之和等于target](https://leetcode.cn/problems/two-sum/)

# 按位异域运算
```
按位异域运算的特点：
(1)a^a=0
(2)a^0=a
(3)a^b^c = a^c^b
```
## 相关题
[丢失的数字](https://leetcode.cn/problems/missing-number/description/)
[只出现一次的数字](https://leetcode.cn/problems/single-number/)

# 快慢指针
```
用于判断环的情况
```
## 相关题 
[快乐数](https://leetcode.cn/problems/happy-number/)
[环形链表](https://leetcode.cn/problems/linked-list-cycle/)

# 双指针
```
和快慢指针不同，每次都只走一步，一般用于同时遍历两个数组或链表(当有两个数组或链表时，可以考虑用它)
然后跟条件进行移动
```
## 相关题
[相交链表](https://leetcode.cn/problems/intersection-of-two-linked-lists/)
[删除有序数组中的重复项](https://leetcode.cn/problems/remove-duplicates-from-sorted-array/)

# 投票算法
```
数组里有一个元素占比大于1/2，则可以用投票法
```
## 相关题
[多数元素](https://leetcode.cn/problems/majority-element/description/)

# 回文符
```
正反读字符串都一样
思路：
    (1)能构成回文串的，字符数量是偶数的
    (2)处理aabaa或aaaaa这种特殊情况
```
## 相关题
[最长回文串](https://leetcode.cn/problems/longest-palindrome/description/)

# 二分查找
```
(1)二分查找，有左区间，右区间，闭区间写法
(2)闭区间，即包括左右元素，使用它比较好理解些
```
## 相关题
[在有序数组里，找到插入的位置](https://leetcode.cn/problems/search-insert-position/description/)

# 动态规划
```
掌握规律，写出方程组
```
## 相关题
[使用最小花费爬楼梯](https://leetcode.cn/problems/min-cost-climbing-stairs/)



# 一些经典题
[最长公共前辍](https://leetcode.cn/problems/longest-common-prefix/description/)
```
(1)假设数组的第一个元素就是最长公共前辍
(2)然后用公共前辍去检查剩下的元素，在检查中不断缩小公共前辍
```
[判断子序列](https://leetcode.cn/problems/is-subsequence/description/)
```
子序列，例如：ace是abcde的子序列，但aec不是(因为顺序不一样)
思路：
    既然要求两个有序的，双指针安排起
```
[最大子数组和](https://leetcode.cn/problems/maximum-subarray/description/)
```
(1)假设f(i)是第i元素之和，那么f(i) = max(f(i-1)+nums[i], nums[i])
(2)再从f(i)中找到最大值即可
```