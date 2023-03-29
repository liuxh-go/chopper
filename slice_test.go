package chopper

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSlice(t *testing.T) {
	s := NewSlice[int]()
	s.FromSlice([]int{1, 2, 3}).Append(3, 4, 5).Filter(func(i int) bool {
		return true
	}).Unique(func(i int) string {
		return strconv.Itoa(i)
	}).Range(func(i, v int) {
		fmt.Println("one", i, v)
	}).Combo(NewSlice[int]().Append(10, 3, 5, 8, 2, 0)).Range(func(i, v int) {
		fmt.Println("two", i, v)
	}).Sort(func(i, i2 int) bool {
		return i < i2
	}).Range(func(i, v int) {
		fmt.Println("three", i, v)
	}).Unique(func(i int) string {
		return strconv.Itoa(i)
	}).Range(func(i, v int) {
		fmt.Println("four", i, v)
	})
}
