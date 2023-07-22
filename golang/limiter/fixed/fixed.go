package fixed

import (
	"note/limiter"
	"sync"
	"time"
)

type fix struct {
	windowLimit   int64
	windowSizeSec int

	count int64
	lastTime time.Time
	mutex sync.Mutex
}

func New(windowLimit int, windowSizeSec int) limiter.Limiter {
	return &fix{windowLimit: int64(windowLimit), windowSizeSec: windowSizeSec, lastTime: time.Now(),count: 1}
}

func (this *fix) Allow() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.count++
	if this.count > this.windowLimit {
		cost:=int(time.Since(this.lastTime).Seconds())
		if  cost>= this.windowSizeSec {
			this.count = 1
			this.lastTime = time.Now()
			return true
		}
		return false
	}

	return true
}
