package cache

import (
	"fmt"
	"testing"
	"time"
)

var key, value string

func Test_Cache_Expire(t *testing.T) {
	c := New(1*time.Second, 100*time.Millisecond)

	key, value = "foo", "bar"
	c.Set(key, value)

	key, value = "lorem", "ipsum"
	c.Set(key, value)

	key = "foo"
	val := c.Get(key)

	if val.(string) != "bar" {
		t.Fatalf("Cache Test failed")
	}

	time.Sleep(1200 * time.Millisecond)

	key = "lorem"
	val = c.Get(key)

	if val != nil {
		t.Fatalf("Cache_Expire Test failed")
	}
}

func Test_Cache_DeleteExpiredWithFunc(t *testing.T) {
	c := New(100*time.Millisecond, 2*time.Second)
	var count int

	for i := 0; i < 100; i++ {
		key, value = string(i), "test"
		c.Set(key, value)
	}

	for i := 0; i < 50; i++ {
		key = string(i)
		c.Get(key)
		time.Sleep(10 * time.Millisecond)
	}
	size := c.Size()
	c.DeleteExpiredWithFunc(func(key string, value interface{}) {
		count++
	})

	if size != count {
		t.Fatalf("Cache_DeleteExpiredWithFunc Test failed")
	}
}

func Test_Cache_DeleteAllWithFunc(t *testing.T) {
	c := New(100*time.Second, 100*time.Second)
	var count int

	for i := 0; i < 100; i++ {
		key, value = string(i), "test"
		c.Set(key, value)
	}

	size := c.Size()
	c.DeleteAllWithFunc(func(key string, value interface{}) {
		count++
	})

	if size != count {
		t.Fatalf("Cache_DeleteAllWithFunc Test failed")
	}
}

func Test_Cache_DeleteAllWithFunc2(t *testing.T) {
	var count int
	c := New2(100*time.Millisecond, 60*time.Millisecond, func(key string, value interface{}) {
		count++
	})

	for i := 0; i < 100; i++ {
		key, value = string(i), "test"
		c.Set(key, value)
	}

	size1 := c.Size()
	time.Sleep(130 * time.Millisecond)
	size2 := c.Size()

	if size1 != 100 || size2 != 0 {
		t.Fatalf("Cache_DeleteAllWithFunc2 Test failed")
	}
}

func Test_Cache_Overwrite(t *testing.T) {
	c := New(5*time.Second, 1*time.Second)

	key, value = "foo", "bar"
	c.Add(key, value)

	value = "ipsum"
	c.Add(key, value)

	val := c.Get(key)

	if val.(string) != "bar" {
		t.Fatalf("Cache_Overwrite Test failed")
	}

	c.Update(key, value)

	val = c.Get(key)

	if val.(string) != "ipsum" {
		t.Fatalf("Cache_Overwrite Test failed")
	}

	key, value = "dolor", "sit"
	b := c.Update(key, value)
	if b != false || c.Get(key) != nil {
		t.Fatalf("Cache_Overwrite Test failed")
	}
}

func Test_Cache_Clear(t *testing.T) {
	c := New(0, 1*time.Second)

	key, value = "foo", "bar"
	c.Add(key, value)

	c.Clear()

	val := c.Get(key)

	if val != nil {
		t.Fatalf("Cache_Clear Test failed")
	}
}

func Test_Cache_Size(t *testing.T) {
	c := New(5*time.Second, 1*time.Second)

	value = "foobar"
	for i := 0; i < 10; i++ {
		c.Add(fmt.Sprintf("Value %v", i), value)
	}

	if c.Size() != 10 {
		t.Fatalf("Cache_Size Test failed")
	}

	c.Delete("Value 5")

	if c.Size() != 9 {
		t.Fatalf("Cache_Size Test failed")
	}

	c.Clear()

	if c.Size() != 0 {
		t.Fatalf("Cache_Size Test failed")
	}
}
