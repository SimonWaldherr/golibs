package cache

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNonpositiveDurationsKeepItems(t *testing.T) {
	for _, duration := range []time.Duration{0, -time.Second} {
		t.Run(duration.String(), func(t *testing.T) {
			c := New(duration, 0)
			t.Cleanup(c.Stop)
			c.Set("default", "value")
			c.SetWithDuration("explicit", "value", time.Now(), duration)
			c.DeleteExpired()
			for _, key := range []string{"default", "explicit"} {
				if value, ok := c.GetWithOK(key); !ok || value != "value" {
					t.Fatalf("GetWithOK(%q) = %v, %v", key, value, ok)
				}
			}
		})
	}
}

func TestZeroCacheAndNilValues(t *testing.T) {
	var c Cache
	c.Stop()
	c.Stop()
	if !c.Time("absent").IsZero() {
		t.Fatal("absent key has a creation time")
	}
	if !c.Add("nil", nil) {
		t.Fatal("first Add failed")
	}
	if value, ok := c.GetWithOK("nil"); value != nil || !ok {
		t.Fatalf("stored nil = %v, %v", value, ok)
	}
	if c.Add("nil", "overwrite") {
		t.Fatal("Add replaced an existing nil value")
	}
	if !c.Update("nil", "updated") || c.Get("nil") != "updated" {
		t.Fatal("Update did not replace an existing nil value")
	}
	c.Delete("nil")
	if c.Update("nil", "resurrected") {
		t.Fatal("Update recreated a deleted key")
	}
	if _, ok := c.GetWithOK("nil"); ok || !c.Time("nil").IsZero() {
		t.Fatal("deleted item remains visible")
	}
}

func TestConcurrentAddHasOneWinner(t *testing.T) {
	c := New(0, 0)
	const writers = 128
	var wg sync.WaitGroup
	var winners int32
	start := make(chan struct{})
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if c.Add("shared", nil) {
				atomic.AddInt32(&winners, 1)
			}
		}()
	}
	close(start)
	wg.Wait()
	if winners != 1 {
		t.Fatalf("Add succeeded %d times, want 1", winners)
	}
}

func expireItem(c *Cache, key string) {
	c.lock.Lock()
	c.items[key].Expiration = time.Now().Add(-time.Second)
	c.lock.Unlock()
}

func TestExpiredItemCanBeAddedButNotUpdated(t *testing.T) {
	c := New(time.Minute, 0)
	c.Set("key", "old")
	expireItem(c, "key")
	if value, ok := c.GetWithOK("key"); value != nil || ok {
		t.Fatalf("expired value = %v, %v", value, ok)
	}
	if c.Update("key", "updated") {
		t.Fatal("Update accepted an expired item")
	}
	if !c.Add("key", "replacement") || c.Get("key") != "replacement" {
		t.Fatal("Add did not replace an expired item")
	}
}

func TestCleanupCallbacksCanReenterCache(t *testing.T) {
	for _, expiredOnly := range []bool{false, true} {
		t.Run(fmt.Sprint(expiredOnly), func(t *testing.T) {
			c := New(time.Minute, 0)
			c.Set("one", 1)
			c.Set("two", 2)
			expireItem(c, "one")
			expireItem(c, "two")
			done := make(chan struct{})
			go func() {
				defer close(done)
				callback := func(key string, value interface{}) {
					c.Size()
					c.Time(key)
					c.Get(key)
					c.Set(key, value)
					c.Stop()
				}
				if expiredOnly {
					c.DeleteExpiredWithFunc(callback)
				} else {
					c.DeleteAllWithFunc(callback)
				}
			}()
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				t.Fatal("cleanup callback deadlocked")
			}
			if c.Size() != 2 || c.Get("one") != 1 || c.Get("two") != 2 {
				t.Fatal("cleanup removed values reinserted by callbacks")
			}
		})
	}
}

func TestCleanupAllowsNilCallback(t *testing.T) {
	c := New(time.Minute, 0)
	c.Set("expired", 1)
	c.Set("live", 2)
	expireItem(c, "expired")
	c.DeleteExpiredWithFunc(nil)
	if c.Size() != 1 || c.Get("live") != 2 {
		t.Fatal("expired cleanup did not retain the live item")
	}
	c.DeleteAllWithFunc(nil)
	if c.Size() != 0 {
		t.Fatal("all-item cleanup retained items")
	}
}

func TestMetadataDuringConcurrentMutation(t *testing.T) {
	c := New(0, 0)
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 300; i++ {
				key := fmt.Sprint(i % 17)
				if worker%2 == 0 {
					c.Set(key, i)
					c.Delete(key)
					if i%13 == 0 {
						c.Clear()
					}
				} else {
					c.Size()
					c.Time(key)
					c.GetWithOK(key)
				}
			}
		}(worker)
	}
	wg.Wait()
}

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (int, error) { return fn(p) }

type reentrantStringer struct{ cache *Cache }

func (s reentrantStringer) String() string {
	s.cache.Set("from-stringer", true)
	return "formatted"
}

func TestExportAndStringDoNotLockUserCode(t *testing.T) {
	c := New(0, 0)
	c.Set("key", "value")
	done := make(chan error, 1)
	go func() {
		var buf bytes.Buffer
		err := c.Export(writerFunc(func(p []byte) (int, error) {
			c.Set("from-writer", true)
			return buf.Write(p)
		}))
		if err == nil {
			c.Set("stringer", reentrantStringer{c})
			_ = c.String()
		}
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("user writer or formatter deadlocked")
	}
	if c.Get("from-writer") != true || c.Get("from-stringer") != true {
		t.Fatal("user code did not run")
	}
}

func TestStopFromJanitorCallback(t *testing.T) {
	var c *Cache
	called := make(chan struct{})
	c = New2(time.Minute, time.Millisecond, func(string, interface{}) {
		c.Stop()
		close(called)
	})
	t.Cleanup(c.Stop)
	c.Set("expired", true)
	expireItem(c, "expired")
	select {
	case <-called:
	case <-time.After(2 * time.Second):
		t.Fatal("janitor did not invoke callback")
	}
	select {
	case <-c.stopped:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop inside callback did not stop janitor")
	}
	c.Stop()
	c.Set("after-stop", "usable")
	if c.Get("after-stop") != "usable" {
		t.Fatal("Stop disabled manual cache operations")
	}
}

func TestConcurrentStop(t *testing.T) {
	c := New(time.Minute, time.Hour)
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Stop()
		}()
	}
	wg.Wait()
	select {
	case <-c.stopped:
	case <-time.After(2 * time.Second):
		t.Fatal("janitor did not exit after concurrent Stop calls")
	}
}
