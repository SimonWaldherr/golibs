// Package cache simplifies caching with GC
package cache

import (
	"encoding/gob"
	"fmt"
	"io"
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

// Export all items to a gob buffer
func (cache *Cache) Export(w io.Writer) error {
	enc := gob.NewEncoder(w)

	cache.lock.RLock()
	defer cache.lock.RUnlock()

	X := make(map[string]interface{})

	for k := range cache.items {
		X[k] = cache.items[k].Object
	}
	if err := enc.Encode(X); err != nil {
		return err
	}
	return nil
}

// Import all items from a gob buffer
func (cache *Cache) Import(r io.Reader) error {
	dec := gob.NewDecoder(r)
	X := make(map[string]interface{})

	if err := dec.Decode(&X); err != nil {
		return err
	}

	for k, v := range X {
		cache.Set(k, v)
	}
	return nil
}

// String returns all cached values as string (for debugging)
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

// Update changes the value of an key. If the key doesn't exist, it returns false
func (cache *Cache) Update(key string, value interface{}) bool {
	item := cache.Get(key)
	if item == nil {
		return false
	}
	cache.Set(key, value)
	return true
}

// DeleteExpired checks all cache items and deletes the expired items
func (cache *Cache) DeleteExpired() {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	for k, v := range cache.items {
		if v.isExpired() {
			delete(cache.items, k)
		}
	}
}

// DeleteExpiredWithFunc does the same like DeleteExpired
// but also calls a function for each deleted item
func (cache *Cache) DeleteExpiredWithFunc(fn func(key string, value interface{})) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	for k, v := range cache.items {
		if v.isExpired() {
			fn(k, cache.items[k].Object)
			delete(cache.items, k)
		}
	}
}

// DeleteAllWithFunc does the same like DeleteExpiredWithFunc
// but not just for the expired items, also the non expired
func (cache *Cache) DeleteAllWithFunc(fn func(key string, value interface{})) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	for k := range cache.items {
		fn(k, cache.items[k].Object)
		delete(cache.items, k)
	}
}

// Size returns the number of cached items
// it does not check for expired items, so run DeleteExpired before
func (cache *Cache) Size() int {
	n := len(cache.items)
	return n
}

// Clear removes all items in the cache
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

// New creates a new Cache
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

// New2 creates a new Cache with a cleaner function
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
