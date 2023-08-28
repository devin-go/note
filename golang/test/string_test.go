package test

import (
	"bytes"
	"strings"
	"testing"
)

func TestStringBuilder(t *testing.T) {
	//string Builder 内部使用[]byte实现
	sb:=strings.Builder{}
	sb.String()//转换为string时，直接指针转换，不用发生拷贝
	//也是用[]byte实现
	bf:=bytes.NewBuffer(nil)
	bf.String()//强制转换，会发生拷贝
	var aa chan string
	t.Log(aa)
}
