// Package cache provides a concurrency-safe in-memory cache with optional expiry.
package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

// Item holds a cached value and its creation and expiration times.
// A zero Expiration means that the item does not expire.
type Item struct {
	Object     interface{}
	Creation   time.Time
	Expiration time.Time
}

func (item *Item) isExpired() bool {
	return item.expiredAt(time.Now())
}

func (item *Item) expiredAt(now time.Time) bool {
	return !item.Expiration.IsZero() && !item.Expiration.After(now)
}

// Cache stores items safely across goroutines. A zero Cache is ready to use.
// Configure Expiration before sharing the cache between goroutines; nonpositive
// durations keep items indefinitely. A Cache must not be copied after first
// use. Stored values remain caller-owned and are not deeply copied.
type Cache struct {
	Expiration time.Duration
	items      map[string]*Item
	lock       sync.RWMutex
	stop       chan struct{}
	stopped    chan struct{}
	stopOnce   sync.Once
}

// snapshot copies item values while holding the map lock. User-provided writers
// and formatters must run after releasing that lock.
func (cache *Cache) snapshot() map[string]interface{} {
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	values := make(map[string]interface{}, len(cache.items))
	for key, item := range cache.items {
		values[key] = item.Object
	}
	return values
}

// Export writes all stored values, including expired items, in gob format.
// Register concrete interface value types with gob when necessary. Expiration
// metadata is not exported; Import applies the receiving cache's default TTL.
func (cache *Cache) Export(w io.Writer) error {
	return gob.NewEncoder(w).Encode(cache.snapshot())
}

// Export exports all items in the given cache to a writer.
// This is a package-level convenience wrapper around (*Cache).Export.
func Export(c *Cache, w io.Writer) error {
	return c.Export(w)
}

// Import merges values from a gob stream, applying the cache's default TTL.
func (cache *Cache) Import(r io.Reader) error {
	values := make(map[string]interface{})
	if err := gob.NewDecoder(r).Decode(&values); err != nil {
		return err
	}

	cache.lock.Lock()
	defer cache.lock.Unlock()
	for key, value := range values {
		cache.setLocked(key, value, time.Now(), cache.Expiration)
	}
	return nil
}

// String returns stored values in key order, including expired items.
func (cache *Cache) String() string {
	values := cache.snapshot()
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var str strings.Builder
	for _, key := range keys {
		fmt.Fprintf(&str, "%s\t%v\n", key, values[key])
	}
	return str.String()
}

func (cache *Cache) setLocked(key string, value interface{}, creation time.Time, duration time.Duration) {
	if cache.items == nil {
		cache.items = make(map[string]*Item)
	}
	var expiration time.Time
	if duration > 0 {
		expiration = time.Now().Add(duration)
	}
	cache.items[key] = &Item{Object: value, Creation: creation, Expiration: expiration}
}

// Set stores a value with the default TTL, replacing any existing item.
func (cache *Cache) Set(key string, value interface{}) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.setLocked(key, value, time.Now(), cache.Expiration)
}

// SetWithDuration stores a value with a supplied creation time and TTL measured
// from this call. Nonpositive durations keep the item indefinitely.
func (cache *Cache) SetWithDuration(key string, value interface{}, creation time.Time, duration time.Duration) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.setLocked(key, value, creation, duration)
}

// Time returns the stored creation time, or zero if the key is absent. Expired
// items retain their metadata until removed by a cleanup operation.
func (cache *Cache) Time(key string) time.Time {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	if item, ok := cache.items[key]; ok {
		return item.Creation
	}
	return time.Time{}
}

// Get returns the cached value, or nil if the key is absent or expired.
// Use GetWithOK to distinguish a stored nil value from an absent item.
func (cache *Cache) Get(key string) interface{} {
	value, _ := cache.GetWithOK(key)
	return value
}

// GetWithOK returns the value and whether an unexpired item exists.
func (cache *Cache) GetWithOK(key string) (interface{}, bool) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	item, ok := cache.items[key]
	if !ok || item.isExpired() {
		return nil, false
	}
	return item.Object, true
}

// Delete removes a key from the cache.
func (cache *Cache) Delete(key string) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	delete(cache.items, key)
}

// Add atomically stores a value only if the key is absent or expired.
func (cache *Cache) Add(key string, value interface{}) bool {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	if item, ok := cache.items[key]; ok && !item.isExpired() {
		return false
	}
	cache.setLocked(key, value, time.Now(), cache.Expiration)
	return true
}

// Update atomically replaces an unexpired item and refreshes its default TTL.
// It returns false if the key is absent or expired, including after deletion.
func (cache *Cache) Update(key string, value interface{}) bool {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	if item, ok := cache.items[key]; !ok || item.isExpired() {
		return false
	}
	cache.setLocked(key, value, time.Now(), cache.Expiration)
	return true
}

// remove collects deleted items under the lock so callbacks may use the cache.
func (cache *Cache) remove(expiredOnly bool) map[string]interface{} {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	now := time.Now()
	removed := make(map[string]interface{})
	for key, item := range cache.items {
		if !expiredOnly || item.expiredAt(now) {
			removed[key] = item.Object
			delete(cache.items, key)
		}
	}
	return removed
}

// DeleteExpired removes expired items.
func (cache *Cache) DeleteExpired() {
	cache.remove(true)
}

// DeleteExpiredWithFunc removes expired items, then calls fn for each value
// outside the cache lock. A nil callback simply removes the items. Callbacks
// may read or modify the cache; newly added items survive the cleanup.
func (cache *Cache) DeleteExpiredWithFunc(fn func(key string, value interface{})) {
	removed := cache.remove(true)
	if fn != nil {
		for key, value := range removed {
			fn(key, value)
		}
	}
}

// DeleteAllWithFunc removes all stored items, then calls fn for each value
// outside the cache lock. A nil callback simply removes the items.
func (cache *Cache) DeleteAllWithFunc(fn func(key string, value interface{})) {
	removed := cache.remove(false)
	if fn != nil {
		for key, value := range removed {
			fn(key, value)
		}
	}
}

// Size returns the number of stored items, including expired items.
func (cache *Cache) Size() int {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	return len(cache.items)
}

// Clear removes all items without invoking cleanup callbacks.
func (cache *Cache) Clear() {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.items = make(map[string]*Item)
}

func newCache(expirationTime, cleanupInterval time.Duration, fn func(string, interface{})) *Cache {
	cache := &Cache{Expiration: expirationTime, items: make(map[string]*Item)}
	if cleanupInterval > 0 {
		cache.stop = make(chan struct{})
		cache.stopped = make(chan struct{})
		go cache.cleaner(cleanupInterval, fn)
	}
	return cache
}

func (cache *Cache) cleaner(interval time.Duration, fn func(string, interface{})) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(cache.stopped)
	for {
		select {
		case <-cache.stop:
			return
		case <-ticker.C:
			// Prefer shutdown when a tick and Stop become ready together.
			select {
			case <-cache.stop:
				return
			default:
			}
			cache.DeleteExpiredWithFunc(fn)
		}
	}
}

// Stop disables automatic cleanup. It is safe to call repeatedly or from a
// cleanup callback. An in-progress cleanup may finish after Stop returns.
// Stored items and manual cache operations remain available after shutdown.
func (cache *Cache) Stop() {
	cache.stopOnce.Do(func() {
		if cache.stop != nil {
			close(cache.stop)
		}
	})
}

// New creates a cache. Nonpositive expiration durations disable expiry;
// nonpositive cleanup intervals disable automatic cleanup. Call Stop when a
// cache with automatic cleanup is no longer needed.
func New(expirationTime, cleanupInterval time.Duration) *Cache {
	return newCache(expirationTime, cleanupInterval, nil)
}

// New2 creates a cache with a callback for automatically removed expired items.
// Callbacks run outside the cache lock. The cache does not handle process
// signals; callers may explicitly flush values with DeleteAllWithFunc.
func New2(expirationTime, cleanupInterval time.Duration, fn func(key string, value interface{})) *Cache {
	return newCache(expirationTime, cleanupInterval, fn)
}
