package cache_test

import (
	"sync"

	"github.com/grafana/gomemcache/memcache"
)

type mockMemcache struct {
	sync.RWMutex
	contents map[string][]byte

	keysFetchedCount int
}

func newMockMemcache() *mockMemcache {
	return &mockMemcache{
		contents: map[string][]byte{},
	}
}

func (m *mockMemcache) GetMulti(keys []string, _ ...memcache.Option) (map[string]*memcache.Item, error) {
	m.RLock()
	defer m.RUnlock()
	m.keysFetchedCount += len(keys)
	result := map[string]*memcache.Item{}
	for _, k := range keys {
		if c, ok := m.contents[k]; ok {
			result[k] = &memcache.Item{
				Value: c,
			}
		}
	}
	return result, nil
}

func (m *mockMemcache) Set(item *memcache.Item) error {
	m.Lock()
	defer m.Unlock()
	m.contents[item.Key] = item.Value
	return nil
}
