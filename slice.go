package chopper

import "sort"

// Slice 泛型切片扩展
type Slice[T any] struct {
	data []T
}

// NewSlice 构造函数
func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{}
}

// FromSlice 从切片构造
func (s *Slice[T]) FromSlice(slice []T) *Slice[T] {
	s.data = make([]T, len(slice))
	copy(s.data, slice)
	return s
}

// Append 添加元素
func (s *Slice[T]) Append(list ...T) *Slice[T] {
	s.data = append(s.data, list...)
	return s
}

// Combo 组合其他泛型切片
func (s *Slice[T]) Combo(anotherList ...*Slice[T]) *Slice[T] {
	for _, list := range anotherList {
		newSlice := make([]T, s.Len()+list.Len())
		copy(newSlice[:s.Len()], s.data)
		copy(newSlice[s.Len():], list.data)
		s.data = make([]T, len(newSlice))
		copy(s.data, newSlice)
	}

	return s
}

// Data 获取原始数据
func (s *Slice[T]) Data() []T {
	r := make([]T, len(s.data))
	copy(r, s.data)
	return r
}

// Len 获取切片长度
func (s *Slice[T]) Len() int {
	return len(s.data)
}

// Sort 排序
func (s *Slice[T]) Sort(fn func(T, T) bool) *Slice[T] {
	sort.Slice(s.data, func(i, j int) bool {
		return fn(s.data[i], s.data[j])
	})

	return s
}

// Unique 去重
func (s *Slice[T]) Unique(keyFn func(T) string) *Slice[T] {
	var r []T
	checkMap := make(map[string]struct{})
	for _, v := range s.data {
		k := keyFn(v)
		if _, exists := checkMap[k]; exists {
			continue
		}

		r = append(r, v)
		checkMap[k] = struct{}{}
	}
	s.data = s.data[:len(r)]
	copy(s.data, r)
	return s
}

// Filter 过滤元素
func (s *Slice[T]) Filter(fn func(T) bool) *Slice[T] {
	var filterSlice []T
	for _, v := range s.data {
		if !fn(v) {
			continue
		}

		filterSlice = append(filterSlice, v)
	}

	s.data = s.data[:len(filterSlice)]
	copy(s.data, filterSlice)
	return s
}

// Range 遍历执行函数
func (s *Slice[T]) Range(fn func(int, T)) *Slice[T] {
	for i, v := range s.data {
		fn(i, v)
	}

	return s
}

// MaxOrMin 获取最大或最小元素
func (s *Slice[T]) MaxOrMin(compareFn func(T, T) bool) (t T) {
	if len(s.data) == 0 {
		return
	}

	t = s.data[0]
	for _, v := range s.data {
		if compareFn(v, t) {
			t = v
		}
	}

	return
}

// First 获取第一个元素
func (s *Slice[T]) First() (t T) {
	if len(s.data) == 0 {
		return
	}

	t = s.data[0]
	return
}

// Last 获取最后一个元素
func (s *Slice[T]) Last() (t T) {
	if len(s.data) == 0 {
		return
	}

	t = s.data[len(s.data)-1]
	return
}

// Reset 重置切片
func (s *Slice[T]) Reset() *Slice[T] {
	s.data = make([]T, 0)
	return s
}
