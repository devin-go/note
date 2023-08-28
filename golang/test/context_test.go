package test

import (
	"context"
	"testing"
)



func TestContextWithValue(t *testing.T) {
	ctx1:=context.WithValue(context.TODO(), "aa",1)
	ctx2,cc:=context.WithCancel(ctx1)
	defer cc()
	ctx3:=context.WithValue(ctx2,"bb",2)
	ctx4:=context.WithValue(ctx3,"aa",11)
	//根据key首先判断当前有没有，如果没有，则一直往上找，找到就停止
	//相同key不会覆盖
	t.Log(ctx1.Value("aa"))
	t.Log(ctx2.Value("aa"))
	t.Log(ctx3.Value("aa"))
	t.Log(ctx4.Value("aa"))
	t.Log(ctx1.Value("aa"))
}
