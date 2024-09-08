package cache

import (
	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Set(key string, value interface{}) error {
	return c.client.Set(key, value, 0).Err()
}

func (c *Cache) Get(key string, value interface{}) error {
	return c.client.Get(key).Scan(value)
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(key).Err()
}

func (c *Cache) LPush(key string, value interface{}) error {
	return c.client.LPush(key, value).Err()
}

func (c *Cache) LRange(key string, start, stop int64, data interface{}) error {
	return c.client.LRange(key, start, stop).ScanSlice(data)
}

func (c *Cache) LPop(key string) (string, error) {
	return c.client.LPop(key).Result()
}

func (c *Cache) LIndex(key string, index int, data interface{}) error {
	return c.client.LIndex(key, int64(index)).Scan(data)
}
