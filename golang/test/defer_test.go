package test

import (
	"testing"
)

func a() (result int) {
	i:=1
	defer func() {
		result++
	}()
	return i
	/*
	实际执行代码
	result = i
	result++
	return
	 */
}
func b() int {
	var i int
	defer func() {
		i++
	}()
	return i
}

type bbObj struct {
	i int
}

func bb() *bbObj {
	out:=&bbObj{i:1}
	defer func() {
		out.i++
	}()
	return out
}

func TestDefer(t *testing.T) {
	t.Log(a())
	t.Log(b(), bb())
	/**
	b()与bb()解析：
		相当于函数有一个匿名变量
		b()的defer()之所以修改不了返回值，是因为匿名变量不是指针
		而bb()之所以能修改是在为 匿名变量是指针
	 */
}

func TestPanic(t *testing.T) {
	defer t.Log("in main")
	defer func() {
		defer func() {
			panic("panic again and again")
		}()
		t.Log("aaaa")
		panic("panic again")
	}()

	panic("once panic")
}

func p(t *testing.T)  {

	defer t.Log("p func defer 1")

	panic("p func panic")
	defer t.Log("p func defer 2")//不会被执行
}

func TestDeferPanic(t *testing.T) {
	defer func() {
		t.Log("main defer")
	}()
	p(t)
	t.Log("main exit")//不会被执行
}



