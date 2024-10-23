package ginendpointusage

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type InMemoryBackend struct {
	stats map[string]int64
	mutex sync.Mutex
}

func (b *InMemoryBackend) Collect(_ context.Context, method, endpoint string, incr int64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.stats == nil {
		b.stats = make(map[string]int64)
	}
	key := fmt.Sprintf("%s::%s", method, endpoint)
	b.stats[key]++
	return nil
}

func (b *InMemoryBackend) GetStats(_ context.Context) ([]Stat, error) {
	var stats []Stat
	for k, v := range b.stats {
		endpointData := strings.Split(k, "::")
		stats = append(stats, Stat{
			Method:   endpointData[0],
			Endpoint: endpointData[1],
			Count:    v,
		})
	}
	return stats, nil
}
