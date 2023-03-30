package tree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testMap = []*struct {
		treeList, dlrList, ldrList, lrdList []int
	}{
		{
			treeList: []int{1, 2, 3, 4, 5, 6, 7},
			dlrList:  []int{1, 2, 4, 5, 3, 6, 7},
			ldrList:  []int{4, 2, 5, 1, 6, 3, 7},
			lrdList:  []int{4, 5, 2, 6, 7, 3, 1},
		},
	}
)

func TestBinary(t *testing.T) {
	for _, v := range testMap {
		result := make([]int, 0, len(v.treeList))
		fn := func(i int) {
			result = append(result, i)
		}

		bt := NewBinaryTree[int](v.treeList...)
		bt.DLR(fn)
		assert.ElementsMatch(t, v.dlrList, result, "dlr")

		result = make([]int, 0, len(v.treeList))
		bt.LDR(fn)
		assert.ElementsMatch(t, v.ldrList, result, "ldr")

		result = make([]int, 0, len(v.treeList))
		bt.LRD(fn)
		assert.ElementsMatch(t, v.lrdList, result, "lrd")
	}

	bt := NewBinaryTree[int](1)
	bt.AddNode(2).AddNode(3).AddNode(4).AddNode(5).AddNode(6).AddNode(7).AddNode(8)
	bt.DLR(func(i int) {
		fmt.Printf("%d ", i)
	})
}
