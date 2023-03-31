package tree

import (
	"fmt"
	"testing"
)

func TestPath(t *testing.T) {
	type data struct {
		fn func(any)
	}

	pn := &PathNode[data]{}

	pn.AddNode("/t", &data{fn: func(a any) {
		fmt.Println("t", a)
	}})
	pn.AddNode("/t1", &data{fn: func(a any) {
		fmt.Println("t1", a)
	}})
	pn.AddNode("/t2", &data{fn: func(a any) {
		fmt.Println("t2", a)
	}})
	pn.AddNode("/t12", &data{fn: func(a any) {
		fmt.Println("t12", a)
	}})
	pn.AddNode("/a", &data{fn: func(a any) {
		fmt.Println("a", a)
	}})
	pn.AddNode("/a1", &data{fn: func(a any) {
		fmt.Println("a1", a)
	}})
	pn.AddNode("/a2", &data{fn: func(a any) {
		fmt.Println("a2", a)
	}})

	v := pn.GetValue("/a1", &[]skippedNode[data]{})
	if v.Data != nil {
		v.Data.fn("test")
	}
}
