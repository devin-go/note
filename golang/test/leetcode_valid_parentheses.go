package test

type stack struct{
	bs []rune
}

func(this *stack) put(b rune) {
	this.bs = append(this.bs, b)
}
func (this * stack) pop() rune {
	size:=len(this.bs)
	if size == 0 {
		return '0'
	}
	out:=this.bs[size-1]
	this.bs = this.bs[:size-1]
	return out
}

func isValid(s string) bool {
	st:=stack{}
	for _,b:=range s {
		if b == '(' || b == '[' || b == '{' {
			st.put(b)
			continue
		}
		p:=st.pop()
		if p == '0' {
			return false
		} else if p == '(' && b !=')' {
			return false
		} else if p == '[' && b != ']' {
			return false
		} else if p == '{' && b != '}' {
			return false
		}
	}
	return true
}
