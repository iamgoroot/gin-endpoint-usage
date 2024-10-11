package ginUsageStats

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisBackend struct {
	RedisClient redis.UniversalClient
}

func (m *RedisBackend) Collect(ctx context.Context, method, endpoint string) error {
	key := fmt.Sprintf("%s::%s", method, endpoint)
	return m.RedisClient.HIncrBy(ctx, "gin-endpoint-usage-stats", key, 1).Err()
}
func (m *RedisBackend) GetStats(ctx context.Context) ([]Stat, error) {
	stats, err := m.RedisClient.HGetAll(ctx, "gin-endpoint-usage-stats").Result()
	if err != nil {
		return nil, err
	}

	var statsSlice []Stat
	for k, v := range stats {
		endpointData := strings.Split(k, "::")
		count, _ := strconv.Atoi(v)
		statsSlice = append(statsSlice, Stat{
			Method:   endpointData[0],
			Endpoint: endpointData[1],
			Count:    count,
		})
	}
	return statsSlice, nil
}
