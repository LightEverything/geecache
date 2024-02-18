package Geecache

import (
	"Geecache/byteview"
	"Geecache/lru"
	"sync"
)

type GeeCache struct {
	cache      *lru.Cache
	cacheGuard sync.RWMutex
	maxBytes   int
}

func (c *GeeCache) Add(key string, value byteview.Byteview) {
	c.cacheGuard.Lock()
	defer c.cacheGuard.Unlock()

	// 延迟初始化
	if c.cache != nil {
		c.cache = (&lru.Cache{}).New(c.maxBytes, nil)
	}
	c.cache.Add(key, value)
}

func (c *GeeCache) Get(key string) (b byteview.Byteview, ok bool) {
	c.cacheGuard.RLock()
	defer c.cacheGuard.RUnlock()

	// 延迟初始化
	if c.cache == nil {
		c.cache = (&lru.Cache{}).New(c.maxBytes, nil)
		return
	}

	if v, ok := c.cache.Get(key); ok {
		return v.(byteview.Byteview), ok
	}
	return b, false
}

// 设置最大值（返回链式调用）
func (c *GeeCache) SetMaxByte(maxByte int) *GeeCache {
	c.maxBytes = maxByte
	return c
}
