package cache

import "time"

type Item struct {
	Object     interface{}
	Expiration time.Time
}

func (item *Item) isExpired() bool {
	if item.Expiration.IsZero() == true {
		return false
	}
	return item.Expiration.Before(time.Now())
}

type Cache struct {
	Expiration time.Duration
	items      map[string]*Item
}

func (cache *Cache) Set(key string, value interface{}) {
	cache.items[key] = &Item{
		Object:     value,
		Expiration: time.Now().Add(cache.Expiration),
	}
}

func (cache *Cache) Get(key string) interface{} {
	item, ok := cache.items[key]
	if !ok || item.isExpired() {
		return nil
	}
	return item.Object
}

func (cache *Cache) Delete(key string) {
	delete(cache.items, key)
}

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
	for k, v := range cache.items {
		if v.isExpired() {
			cache.Delete(k)
		}
	}
}

func (cache *Cache) Size() int {
	n := len(cache.items)
	return n
}

func (cache *Cache) Clear() {
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
