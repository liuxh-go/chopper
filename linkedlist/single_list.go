package linkedlist

// NewSingleList 新建单向链表
func NewSingleList[T any]() *SingleNode[T] {
	return &SingleNode[T]{}
}

// IsNil 是否为空
func (sn *SingleNode[T]) IsNil() bool {
	return sn == nil || sn.Node == nil
}

// Push 追加元素
func (sn *SingleNode[T]) Push(list ...T) *SingleNode[T] {
	cur := sn
	for _, t := range list {
		node := &Node[T]{
			data: t,
		}

		if cur.IsNil() {
			cur.Node = node
		} else {
			cur.next = &SingleNode[T]{
				Node: node,
			}
			cur = cur.next
		}
	}

	return sn
}

// InsertHead 插入头部
func (sn *SingleNode[T]) InsertHead(t T) *SingleNode[T] {
	var next *SingleNode[T]
	if !sn.IsNil() {
		cur := *sn
		next = &cur
	}

	*sn = SingleNode[T]{
		Node: &Node[T]{
			data: t,
		},
		next: next,
	}

	return sn
}

// Pop 从尾部弹出元素
func (sn *SingleNode[T]) Pop() (t T, exists bool) {
	cur := sn
	for ; !cur.IsNil() && !cur.next.IsNil() && !cur.next.next.IsNil(); cur = cur.next {
	}
	if !cur.IsNil() {
		if cur.next.IsNil() {
			t = cur.data
			sn.Node = nil
			exists = true
			return
		}

		t = cur.next.data
		cur.next = nil
		exists = true
	}

	return
}

// Next 迭代
func (sn *SingleNode[T]) Next() *SingleNode[T] {
	if !sn.IsNil() {
		return sn.next
	}

	return nil
}

// Len 获取长度
func (sn *SingleNode[T]) Len() int {
	l := 0
	for cur := sn; !cur.IsNil(); cur = cur.next {
		l++
	}

	return l
}

// Delete 删除元素
func (sn *SingleNode[T]) Delete(deleteFn func(T) bool) *SingleNode[T] {
	p := &SingleNode[T]{
		Node: &Node[T]{},
		next: sn,
	}

	for cur := p; !cur.IsNil() && !cur.next.IsNil(); {
		if deleteFn(cur.next.data) {
			// 删除一起的元素
			c := cur.next
			for ; !c.next.IsNil(); c = c.next {
				if !deleteFn(c.next.data) {
					break
				}
			}
			cur.next = c.next
		}
		cur = cur.next
	}

	if p.next.IsNil() {
		return sn.Reset()
	}

	sn.Node = p.next.Node
	sn.next = p.next.next

	return sn
}

// List 转换为切片
func (sn *SingleNode[T]) List() []T {
	var result []T
	for cur := sn; !cur.IsNil(); cur = cur.next {
		result = append(result, cur.data)
	}

	return result
}

// Each 遍历执行方法
func (sn *SingleNode[T]) Each(fn func(T)) {
	for cur := sn; !cur.IsNil(); cur = cur.next {
		fn(cur.data)
	}
}

// Edit 修改元素
func (sn *SingleNode[T]) Edit(fn func(T) T) *SingleNode[T] {
	for cur := sn; !cur.IsNil(); cur = cur.next {
		cur.data = fn(cur.data)
	}

	return sn
}

// Search 查找元素
func (sn *SingleNode[T]) Search(fn func(T) bool) (t T, exists bool) {
	for cur := sn; !cur.IsNil(); cur = cur.next {
		if !fn(cur.data) {
			continue
		}

		t = cur.data
		exists = true
		return
	}

	return
}

// Reset 重置链表
func (sn *SingleNode[T]) Reset() *SingleNode[T] {
	sn.Node = nil
	sn.next = nil
	return sn
}
