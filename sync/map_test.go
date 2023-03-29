package sync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := NewMap[int, int]()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Set(i, 2*i)
		}()
	}
	wg.Wait()
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Delete(i)
		}()
	}
	wg.Wait()

	assert.Equal(t, 15, m.Len())

	for i := 0; i < 20; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, exists := m.Get(i)
			if i < 5 {
				assert.Equal(t, false, exists, i)
			} else {
				assert.Equal(t, true, exists, i)
				assert.Equal(t, v, 2*i)
			}
		}()
	}
	wg.Wait()

	m.Range(func(k, v int) int {
		if k > 5 && k <= 10 {
			return v - 1
		}

		return v
	})

	for i := 0; i < 20; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, _ := m.Get(i)
			if i > 5 && i <= 10 {
				assert.Equal(t, v, i*2-1)
			}
		}()
	}
	wg.Wait()

	fData := m.Filter(func(k, v int) bool {
		return v%k != 0
	})
	assert.Equal(t, 5, len(fData))
	assert.Equal(t, 15, len(m.Keys()))
	assert.Equal(t, 15, len(m.Values()))

	m.Reset()
	assert.Equal(t, 0, m.Len())
}
