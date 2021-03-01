package cache

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
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
		key, value = fmt.Sprint(i), "test"
		c.Set(key, value)
	}

	for i := 0; i < 50; i++ {
		key = fmt.Sprint(i)
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
		key, value = fmt.Sprint(i), "test"
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
		key, value = fmt.Sprint(i), "test"
		c.Set(key, value)
	}

	size1 := c.Size()
	time.Sleep(330 * time.Millisecond)
	size2 := c.Size()

	if size1 != 100 || size2 != 0 {
		t.Fatalf("Cache_DeleteAllWithFunc2 Test failed")
	}
}

func Test_Cache_Overwrite(t *testing.T) {
	c := New(5*time.Second, 1*time.Second)

	key, value = "foo", "bar"
	c.Add(key, value)

	duration, _ := time.ParseDuration("2h30m")
	c.SetWithDuration(key, value, time.Now(), duration)
	ti := c.Time(key).Unix()
	if time.Now().Unix() != ti {
		t.Fatalf("Cache_Time Test failed")
	}

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

func Test_Cache_Clear2(t *testing.T) {
	var count int
	c := New2(0, 1*time.Second, func(key string, value interface{}) {
		count++
	})

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

func Test_Cache_String(t *testing.T) {
	c := New(5*time.Second, 1*time.Second)

	c.Add("0", 23)
	c.Add("1", 1*time.Second)
	c.Add("2", "foobar")
	c.Add("3", false)
	c.Add("4", []byte("test"))

	if c.String() != "0\t23\n1\t1s\n2\tfoobar\n3\tfalse\n4\t[116 101 115 116]\n" {
		t.Fatalf("Cache_String Test failed")
	}
}

func Test_Concurrent(t *testing.T) {
	var wg sync.WaitGroup
	var count int

	rand.Seed(time.Now().UnixNano())

	wg.Add(3100)

	c := New2(10*time.Second, 1*time.Second, func(key string, value interface{}) {
		count++
	})

	go func() {
		for i := 0; i < 100; i++ {
			c.Add(fmt.Sprintf("i%v", i), fmt.Sprintf("v%v", i*10))
			wg.Done()
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			r := rand.Intn(63)
			c.Set(fmt.Sprintf("i%v", r), fmt.Sprintf("v%v", i*10))
			wg.Done()
		}
	}()
	go func() {
		for i := 0; i < 2000; i++ {
			r := rand.Intn(63)
			c.Get(fmt.Sprintf("i%v", r))
			wg.Done()
		}
	}()
	wg.Wait()
	c.Clear()
}

func Test_Export(t *testing.T) {
	c := New(5*time.Second, 1*time.Second)

	c.Add("0", 23)
	c.Add("1", 1*time.Second)
	c.Add("2", "foobar")
	c.Add("3", false)
	c.Add("4", []byte("test"))

	var buf bytes.Buffer

	fmt.Println(buf.Len())

	//fmt.Printf("Pre-Exported: %#v\n", buf)
	fmt.Printf("Buffer-Len: %v\n", buf.Len())
	c.Export(&buf)
	fmt.Printf("Buffer-Len: %v\n", buf.Len())
	//fmt.Printf("Exported: %#v\n", buf)

	c2 := New(15*time.Second, 1*time.Second)
	c2.Import(&buf)
	fmt.Printf("Buffer-Len: %v\n", buf.Len())
	//fmt.Println(&buf)
	fmt.Println(c2.String())
}

func BenchmarkAdd(b *testing.B) {
	var count int
	c := New2(0, 1*time.Second, func(key string, value interface{}) {
		count++
	})
	for i := 0; i < b.N; i++ {
		c.Add(fmt.Sprintf("i%v", i), fmt.Sprintf("v%v", i*10))
	}

	c.Clear()
}

func BenchmarkAddGet(b *testing.B) {
	var count int
	c := New2(0, 1*time.Second, func(key string, value interface{}) {
		count++
	})
	for i := 0; i < b.N; i++ {
		c.Add(fmt.Sprintf("i%v", i), fmt.Sprintf("v%v", i*10))
		c.Get(fmt.Sprintf("i%v", i))
	}

	c.Clear()
}

func BenchmarkGet(b *testing.B) {
	var count int
	c := New2(0, 1*time.Second, func(key string, value interface{}) {
		count++
	})
	c.Add("foo", "bar")
	for i := 0; i < b.N; i++ {
		c.Get("foo")
	}

	c.Clear()
}
