package tokenbucket

import (
	"note/log"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	log.Init(t)
	l:=New(10,2)
	done:=make(chan struct{})
	callFunc:= func(id int) {
		for  {
			select {
			case <-done:
				return
			default:
				if l.Allow() {
					log.Infof("thread:%d get",id)
				}
			}
		}
	}
	count:=10
	for count>0 {
		go callFunc(count)
		count--
	}
	time.Sleep(10*time.Second)
	close(done)
}
