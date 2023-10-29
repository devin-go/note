package test

import (
	"sync"
	"testing"
	"time"
)

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	counter := 10
	waiter := 5
	wg.Add(counter)
	for counter > 0 {
		go func(i int) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			t.Logf("goroutine:%d wakeup", i)
		}(counter)
		counter--
	}

	for waiter > 0 {
		go func(i int) {
			wg.Wait()
			t.Logf("waitter:%d wakeup", i)
		}(waiter)
		waiter--
	}
	wg.Wait()
	t.Logf("main wakeup")
	time.Sleep(time.Millisecond)
}
