// Package cache simplifies caching with GC
package cache

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

// Item as a single object with an interface as value and two time.Time values for creation and expiration dates
type Item struct {
	Object     interface{}
	Creation   time.Time
	Expiration time.Time
}

func (item *Item) isExpired() bool {
	if item.Expiration.IsZero() {
		return false
	}
	return item.Expiration.Before(time.Now())
}

// Cache is a storage for all Item items
type Cache struct {
	Expiration time.Duration
	items      map[string]*Item
	lock       sync.RWMutex
}

func (cache *Cache) String() string {
	var str string
	var keys []string

	cache.lock.RLock()
	defer cache.lock.RUnlock()

	for k := range cache.items {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		str += k + "\t" + fmt.Sprintf("%v", cache.items[k].Object) + "\n"
	}

	return str
}

// Set creates an Item in the cache, if there is already an item with that name it get overwritten
func (cache *Cache) Set(key string, value interface{}) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	cache.items[key] = &Item{
		Object:     value,
		Creation:   time.Now(),
		Expiration: time.Now().Add(cache.Expiration),
	}
}

// SetWithDuration does the same as Set but with an specific expiration date
func (cache *Cache) SetWithDuration(key string, value interface{}, creation time.Time, duration time.Duration) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	cache.items[key] = &Item{
		Object:     value,
		Creation:   creation,
		Expiration: time.Now().Add(duration),
	}
}

// Time returns the creation date of a cached item
func (cache *Cache) Time(key string) time.Time {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	return cache.items[key].Creation
}

// Get returns the value of a cached item or nil if expired
func (cache *Cache) Get(key string) interface{} {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	item, ok := cache.items[key]
	if !ok || item.isExpired() {
		return nil
	}
	return item.Object
}

// Delete deletes a cached item
func (cache *Cache) Delete(key string) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	delete(cache.items, key)
}

// Add creates an cached item
func (cache *Cache) Add(key string, value interface{}) bool {
	item := cache.Get(key)
	if item != nil {
		return false
	}
	cache.Set(key, value)
	return true
}

func (cache *Cache) Update(key string, value interface{}) bool {
	item := cache.Get(key)
	if item == nil {
		return false
	}
	cache.Set(key, value)
	return true
}

func (cache *Cache) DeleteExpired() {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	for k, v := range cache.items {
		if v.isExpired() {
			delete(cache.items, k)
		}
	}
}

func (cache *Cache) DeleteExpiredWithFunc(fn func(key string, value interface{})) {
	for k, v := range cache.items {
		if v.isExpired() {
			fn(k, cache.items[k].Object)
			cache.lock.Lock()
			delete(cache.items, k)
			cache.lock.Unlock()
		}
	}
}

func (cache *Cache) DeleteAllWithFunc(fn func(key string, value interface{})) {
	for k := range cache.items {
		fn(k, cache.items[k].Object)
		cache.lock.Lock()
		delete(cache.items, k)
		cache.lock.Unlock()
	}
}

func (cache *Cache) Size() int {
	n := len(cache.items)
	return n
}

func (cache *Cache) Clear() {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	cache.items = map[string]*Item{}
}

func cleaner(cache *Cache, interval time.Duration) {
	ticker := time.Tick(interval)
	for {
		select {
		case <-ticker:
			cache.DeleteExpired()
		}
	}
}

func cleanerWithFunc(cache *Cache, interval time.Duration, fn func(key string, value interface{})) {
	defer cache.DeleteAllWithFunc(fn)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ticker := time.Tick(interval)
	for {
		select {
		case <-ticker:
			cache.DeleteExpiredWithFunc(fn)
		case <-c:
			cache.DeleteAllWithFunc(fn)
			os.Exit(1)
		}
	}
}

func New(expirationTime, cleanupInterval time.Duration) *Cache {
	items := make(map[string]*Item)
	if expirationTime == 0 {
		expirationTime = -1
	}
	cache := &Cache{
		Expiration: expirationTime,
		items:      items,
	}
	go cleaner(cache, cleanupInterval)

	return cache
}

func New2(expirationTime, cleanupInterval time.Duration, fn func(key string, value interface{})) *Cache {
	items := make(map[string]*Item)
	if expirationTime == 0 {
		expirationTime = -1
	}
	cache := &Cache{
		Expiration: expirationTime,
		items:      items,
	}
	go cleanerWithFunc(cache, cleanupInterval, fn)

	return cache
}
