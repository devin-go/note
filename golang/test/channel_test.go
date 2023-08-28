package test

import (
	"testing"
)


var nilCh chan string

func readNil(t *testing.T)  {
	select {
	case <-nilCh:// 永久阻塞
	}
	t.Log("readNil exit")
}
func writeNil()  {
	nilCh<-""// 永久阻塞
}

func TestChannel(t *testing.T) {
	//readNil(t)
	//writeNil()
	var p *int
	var i interface{}
	i = p
	t.Log(i==nil)
}
