package sync

import "sync"

// Map 支持泛型的线程安全map
// K : comparable类型的键
// V : any类型的数据
type Map[K comparable, V any] struct {
	data  map[K]V
	mutex sync.RWMutex
}

// NewMap 构造函数
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		data: make(map[K]V),
	}
}

// Set 设置键值对
func (m *Map[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = value
}

// Get 通过键获取值
func (m *Map[K, V]) Get(key K) (value V, exists bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, exists = m.data[key]
	return
}

// Delete 从map删除键值对
func (m *Map[K, V]) Delete(key K) {
	if func() bool {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
		_, exists := m.data[key]

		return !exists
	}() {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, key)
}

// Len 获取map长度
func (m *Map[K, V]) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.data)
}

// Keys 获取所有的key
func (m *Map[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]K, 0, len(m.data))
	for k := range m.data {
		result = append(result, k)
	}

	return result
}

// Values 获取所有的value
func (m *Map[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]V, 0, len(m.data))
	for _, v := range m.data {
		result = append(result, v)
	}

	return result
}

// Filter 过滤数据
func (m *Map[K, V]) Filter(fn func(K, V) bool) map[K]V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[K]V)
	for k, v := range m.data {
		if !fn(k, v) {
			continue
		}

		result[k] = v
	}

	return result
}

// Range 遍历执行
func (m *Map[K, V]) Range(fn func(K, V) V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.data {
		m.data[k] = fn(k, v)
	}
}

// Reset 重置map
func (m *Map[K, V]) Reset() {
	m.data = make(map[K]V)
}
