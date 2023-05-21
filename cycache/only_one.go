package cycache

import (
	"sync"
	"time"
)

// OnlyOneCache 本地缓存
type OnlyOneCache[T any] struct {
	timeoutSeconds int64
	mutex          sync.RWMutex
	loader         func() (T, bool)

	data T
}

// NewOnlyOneCache 新建缓存
func NewOnlyOneCache[T any](cacheTimeoutSeconds int64, loader func() (T, bool)) *OnlyOneCache[T] {
	cache := &OnlyOneCache[T]{
		timeoutSeconds: cacheTimeoutSeconds,
		loader:         loader,
	}
	cache.RefreshCache()
	go cache.startRefresh() // 开启新线程
	return cache

}

// Load 加载数据
func (c *OnlyOneCache[T]) Load() T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.data
}

// 循环刷新
func (c *OnlyOneCache[T]) startRefresh() {
	for {
		time.Sleep(time.Second * time.Duration(c.timeoutSeconds))
		c.RefreshCache()
	}
}

// 刷新缓存
func (c *OnlyOneCache[T]) RefreshCache() {
	newData, ok := c.loader()
	if !ok { // 是否成功
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = newData
}
