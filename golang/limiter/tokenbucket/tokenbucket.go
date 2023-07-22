package tokenbucket

import (
	"math"
	"note/limiter"
	"sync"
	"time"
)

type tokenBucket struct {
	capacity int
	getRate int
	mutex sync.Mutex

	tokens int
	lastTimeSec int64
}

func New(capacity int, genRate int)  limiter.Limiter{
	out:=tokenBucket{capacity: capacity,getRate: genRate,tokens: capacity,lastTimeSec: time.Now().Unix()}
	return &out
}

func (this *tokenBucket) Allow() bool{
	this.mutex.Lock()
	defer this.mutex.Unlock()
	curSec:=time.Now().Unix()
	needGenTokens:=int(curSec-this.lastTimeSec)*this.getRate
	this.tokens = int(math.Min(float64(needGenTokens+this.tokens),float64(this.capacity)))
	this.lastTimeSec = curSec
	if this.tokens > 0 {
		this.tokens--
		return true
	}
	return false
}
