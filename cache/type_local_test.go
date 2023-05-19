package cache

import (
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTypeLocalCache(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	start := time.Now()
	cache := NewTypeLocalCache(10, 20, func(customerID string) (string, bool) {
		return "C" + customerID, true
	})

	var wg sync.WaitGroup
	var count = 10 // 10*100000次缓存请求,10
	for count > 0 {
		count--

		for i := 1; i <= 100000; i++ {
			wg.Add(1)

			go func(customerID string) {
				conf, _, err := cache.Load(customerID, customerID)
				if err != nil {
					log.Println(err)
				}
				if "C"+customerID != conf {
					log.Println("error")
				}
				wg.Done()
			}(strconv.Itoa(i))
		}
	}
	wg.Wait()

	log.Println("used: ", time.Since(start))
}

func TestTypeLocalCacheCheckExpired(t *testing.T) {
	start := time.Now()
	cache := NewTypeLocalCache(5, 6, func(customerID string) (string, bool) {
		return "C" + customerID, true
	}) // 5秒过期
	var wg sync.WaitGroup
	for j := 0; j < 2; j++ {
		for i := 1; i <= 10000; i++ {
			wg.Add(1)

			go func(customerID string) {
				conf, _, err := cache.Load(customerID, customerID)
				if err != nil {
					log.Println(err)
				}
				if "C"+customerID != conf {
					log.Println("error")
				}
				wg.Done()
			}(strconv.Itoa(i))
		}
	}
	wg.Wait()

	if cache.dataMap.Size() != 10000 { // 缓存没有过期，应该是10000个数据
		t.Fatal("TestCheckExpired error")
	}

	time.Sleep(8 * time.Second) //8秒都会过期

	if cache.dataMap.Size() != 0 { // 缓存都过期了，应该是0个数据
		t.Fatal("TestCheckExpired error")
	}

	log.Println(time.Since(start))
}

func TestTypeLocalCacheClose(t *testing.T) {
	start := time.Now()
	cache := NewTypeLocalCache(5, 6, func(customerID string) (string, bool) {
		return "C" + customerID, true
	}) // 5秒过期
	var wg sync.WaitGroup
	for i := 1; i <= 10000; i++ {
		wg.Add(1)

		go func(customerID string) {
			_, _, err := cache.Load(customerID, customerID)
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}(strconv.Itoa(i))
	}
	wg.Wait()

	cache.Close() //关闭

	if cache.dataMap.Size() != 10000 { // 缓存没有过期，应该是10000个数据
		t.Fatal("TestCheckExpired error")
	}

	time.Sleep(8 * time.Second) //8秒都会过期，但是关闭了检查goroutine

	if cache.dataMap.Size() != 10000 { //10000
		t.Fatal("TestCheckExpired error")
	}

	log.Println(time.Since(start))
}
