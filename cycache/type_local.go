package cycache

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type TypeLocalCacheItem[T any] struct {
	Time int64
	Data T
}

// TypeLocalCache 本地缓存
type TypeLocalCache[P any, T any] struct {
	timeoutSeconds int64 // 数据失效时间
	clearSeconds   int64 // 清理数据的时间，大于timeoutSeconds时，会尽量避免缓存镂空
	mutex          sync.RWMutex
	sg             singleflight.Group // 用于防止单用户缓存穿透
	quit           chan int

	dataMap  *SafeMap    // map[string]*TypeLocalCacheItem
	keyQueue stringQueue // key队列，用于过期检查

	loader func(P) (T, bool)
	zero   T
}

// NewTypeLocalCache 新建缓存
func NewTypeLocalCache[P any, T any](cacheTimeoutSeconds, cacheClearSeconds int64, loader func(P) (T, bool)) *TypeLocalCache[P, T] {
	cache := &TypeLocalCache[P, T]{
		timeoutSeconds: cacheTimeoutSeconds,
		clearSeconds:   cacheClearSeconds,
		dataMap:        NewSafeMap(), //map[string]*TypeLocalCacheItem{},
		keyQueue:       stringQueue{},
		quit:           make(chan int),
		loader:         loader,
	}
	go cache.checkExpired() // 开启过期检查goroutine
	return cache
}

// Load 加载数据
func (c *TypeLocalCache[P, T]) Load(key string, param P) (T, bool, error) {
	if key == "" {
		return c.zero, false, errors.New("key must be greater than 0")
	}

	now := time.Now().UnixNano() / int64(time.Second)
	data, ok, isTimeout := c.fromCache(key, now)
	if ok { // 命中返回
		if isTimeout { // 缓存超时，刷新
			go c.loadRealData(key, now, param)
		}
		return data, true, nil
	}

	data, isHit := c.loadRealData(key, now, param)
	return data, isHit, nil
}

func (c *TypeLocalCache[P, T]) loadRealData(key string, now int64, param P) (T, bool) {
	isHit := false
	dataFromLoader, _, _ := c.sg.Do(key, func() (interface{}, error) { // singleflight.Group保护loader的调用
		_data, ok, isTimeout := c.fromCache(key, now) // 二次校验，保障在有缓存时不调用loader
		if !ok || isTimeout {                         // 二次校验，没有命中
			_data, ok = c.loader(param)
			if ok {
				c.toCache(key, now, _data)
			}
		} else {
			isHit = true
		}
		return _data, nil
	})
	return dataFromLoader.(T), isHit
}

// Close 关闭缓存,关闭过期检查goroutine
func (c *TypeLocalCache[P, T]) Close() {
	c.quit <- 1
}

// 获取本地缓存数据数据
func (c *TypeLocalCache[P, T]) fromCache(key string, now int64) (T, bool, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cacheData, ok := c.dataMap.Get(key)
	if !ok || cacheData == nil { // 不存在的缓存
		return c.zero, false, true
	}
	item := cacheData.(*TypeLocalCacheItem[T])
	timeSub := now - item.Time
	if timeSub > c.clearSeconds { // 已经过期
		return c.zero, false, true
	}

	return item.Data, true, timeSub > c.timeoutSeconds
}

// 保存数据到本地缓存
func (c *TypeLocalCache[P, T]) toCache(key string, now int64, data T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.dataMap.Set(key, &TypeLocalCacheItem[T]{
		Data: data,
		Time: now,
	})
	c.keyQueue.push(key)
}

// 过期检查
func (c *TypeLocalCache[P, T]) checkExpired() {
	for {
		select {
		case <-c.quit:
			return
		case <-time.After(time.Second * 1): // 所有过期数据都检查完毕，休息1s继续
			for c.check10000() {
			}
		}
	}
}

// 10000个一次检查,防止锁定时间太长，返回是否需要继续检查
func (c *TypeLocalCache[P, T]) check10000() bool {
	now := time.Now().UnixNano() / int64(time.Second)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.keyQueue.clearDuplicateData()
	for i := 0; i < 10000; i++ {
		key := c.keyQueue.popNoDelete() // 取值
		if key != "" {
			cacheData, ok := c.dataMap.Get(key)
			if !ok || cacheData == nil { // 已经过期和不存在的缓存
				return false
			}
			item := cacheData.(*TypeLocalCacheItem[T])
			if item != nil && now-item.Time > c.clearSeconds { // 已经过期的缓存
				c.dataMap.Del(key)
				c.keyQueue.popOnlyDelete() // 确认删除
			} else {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
