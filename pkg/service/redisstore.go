package service

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const roomKeyPrefix = "agenticsfu:room:"

// RedisStore provides distributed persistent room metadata via Redis.
type RedisStore struct {
	client redis.UniversalClient
}

func NewRedisStore(client redis.UniversalClient) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (r *RedisStore) SetRoomMeta(ctx context.Context, roomName string, meta string) error {
	key := fmt.Sprintf("%s%s", roomKeyPrefix, roomName)
	return r.client.Set(ctx, key, meta, 0).Err()
}

func (r *RedisStore) GetRoomMeta(ctx context.Context, roomName string) (string, error) {
	key := fmt.Sprintf("%s%s", roomKeyPrefix, roomName)
	return r.client.Get(ctx, key).Result()
}

func (r *RedisStore) DeleteRoomMeta(ctx context.Context, roomName string) error {
	key := fmt.Sprintf("%s%s", roomKeyPrefix, roomName)
	return r.client.Del(ctx, key).Err()
}
