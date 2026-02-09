package cache_redis

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/bamgoo/bamgoo"
	"github.com/bamgoo/cache"
	"github.com/redis/go-redis/v9"
)

type redisDriver struct{}

type redisConnection struct {
	client *redis.Client
}

func init() {
	bamgoo.Register("redis", &redisDriver{})
}

func (d *redisDriver) Connect(inst *cache.Instance) (cache.Connect, error) {
	addr, _ := inst.Config.Setting["server"].(string)
	if addr == "" {
		addr, _ = inst.Config.Setting["addr"].(string)
	}
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	username, _ := inst.Config.Setting["username"].(string)
	password, _ := inst.Config.Setting["password"].(string)

	db := 0
	if v, ok := inst.Config.Setting["database"].(int); ok {
		db = v
	} else if v, ok := inst.Config.Setting["database"].(int64); ok {
		db = int(v)
	} else if v, ok := inst.Config.Setting["database"].(string); ok {
		if n, err := strconv.Atoi(v); err == nil {
			db = n
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	return &redisConnection{client: client}, nil
}

func (c *redisConnection) Open() error  { return nil }
func (c *redisConnection) Close() error { return c.client.Close() }

func (c *redisConnection) Read(key string) ([]byte, error) {
	val, err := c.client.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (c *redisConnection) Write(key string, val []byte, expire time.Duration) error {
	return c.client.Set(context.Background(), key, val, expire).Err()
}

func (c *redisConnection) Exists(key string) (bool, error) {
	cnt, err := c.client.Exists(context.Background(), key).Result()
	return cnt > 0, err
}

func (c *redisConnection) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *redisConnection) Sequence(key string, start, step int64, expire time.Duration) (int64, error) {
	val, err := c.client.IncrBy(context.Background(), key, step).Result()
	if err != nil {
		return -1, err
	}
	if val == step {
		if start != 0 {
			val = start
			_ = c.client.Set(context.Background(), key, val, expire).Err()
		} else if expire > 0 {
			_ = c.client.Expire(context.Background(), key, expire).Err()
		}
	}
	if expire > 0 {
		_ = c.client.Expire(context.Background(), key, expire).Err()
	}
	return val, nil
}

func (c *redisConnection) Keys(prefix string) ([]string, error) {
	if prefix == "" {
		prefix = "*"
	} else if !strings.HasSuffix(prefix, "*") {
		prefix = prefix + "*"
	}
	return c.client.Keys(context.Background(), prefix).Result()
}

func (c *redisConnection) Clear(prefix string) error {
	keys, err := c.Keys(prefix)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(context.Background(), keys...).Err()
}
