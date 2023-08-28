package test

import (
	"log"
	"testing"
)

func A()  {
	defer A1()
	defer A2()
	panic("panic a")
}

func A1()  {
	e:=recover()
	log.Printf("A1 func recover:%v",e)
	log.Print("A1 func")
}
func A2()  {
	defer B()
	panic("panic a2")
}

func B()  {
	e:=recover()
	log.Printf("B func recover:%v",e)
	log.Print("B func")
}



func TestRecover(t *testing.T) {
	A()
}
