package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Locker interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (token string, acquired bool, err error)
	Release(ctx context.Context, key string, token string) error
}

type locker struct {
	client *redis.Client
}

var releaseScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
end
return 0
`)

func NewLocker(client *redis.Client) Locker {
	return &locker{client: client}
}

func (l *locker) Acquire(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
	token := uuid.NewString()
	acquired, err := l.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return "", false, err
	}
	if !acquired {
		return "", false, nil
	}
	return token, true, nil
}

func (l *locker) Release(ctx context.Context, key string, token string) error {
	_, err := releaseScript.Run(ctx, l.client, []string{key}, token).Result()
	return err
}
