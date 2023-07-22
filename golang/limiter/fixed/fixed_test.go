package fixed

import (
	"note/log"
	"testing"
	"time"
)


func TestNew(t *testing.T) {
	log.Init(t)
	l:= New(2, 1)
	done:=make(chan struct{})
	var callFuc = func(id int) {
		for {
			select {
			case <-done:
				return
			default:
				for l.Allow(){
					log.Infof("thread:%d get",id)
				}
			}
		}
	}
	threads:=10
	for threads>0 {
		go callFuc(threads)
		threads--
	}

	time.Sleep(10*time.Second)
	close(done)
}
