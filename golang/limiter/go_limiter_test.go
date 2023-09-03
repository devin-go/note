package limiter

import (
	"testing"
	"time"
	"golang.org/x/time/rate"
)

func TestGoLimiter(t *testing.T) {
	limiter := rate.NewLimiter(rate.Every(time.Millisecond*31), 2)
limiter.AllowN()
}
