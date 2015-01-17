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
}

func Test_Cache_Clear(t *testing.T) {
	c := New(5*time.Second, 1*time.Second)

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
