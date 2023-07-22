package log

import (
	"testing"
	"time"
)

var _t *testing.T

func Init(t *testing.T)  {
	_t = t
}
func Infof(format string, a ...interface{})  {
	aa:=[]interface{}{}
	aa = append(aa, time.Now())
	aa = append(aa, a...)
	_t.Logf(" %v " + format,aa...)
}
