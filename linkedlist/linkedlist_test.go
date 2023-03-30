package linkedlist

import (
	"fmt"
	"testing"
)

func TestSingleList(t *testing.T) {
	printFn := func(i int) {
		fmt.Println(i)
	}

	sl := NewSingleList[int]()
	sl.Push(1, 2, 3)

	l := sl.Len()
	for i := 0; i < l; i++ {
		t1, exists := sl.Pop()
		fmt.Println("pop", t1, exists)
	}
	//sl.Each(printFn)

	sl.InsertHead(2)
	sl.Push(8, 5, 2, 2, 4, 2, 2)
	sl.InsertHead(2).Delete(func(i int) bool {
		return i == 2
	}).Each(printFn)

	sl.Edit(func(i int) int {
		if i%2 == 0 {
			return i - 1
		}

		return i
	}).Each(printFn)
	v, exists := sl.Search(func(i int) bool {
		return i == 5
	})
	fmt.Println(v, exists)
}
