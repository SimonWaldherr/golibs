package re

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func drainDo(t *testing.T, data <-chan interface{}) []interface{} {
	t.Helper()
	var values []interface{}
	deadline := time.NewTimer(2 * time.Second)
	defer deadline.Stop()
	for {
		select {
		case value, ok := <-data:
			if !ok {
				return values
			}
			values = append(values, value)
		case <-deadline.C:
			t.Fatal("output did not close after stopping")
			return values
		}
	}
}

func TestDoStopDrainsActiveCallback(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var calls int32
	data, stop := Do(time.Hour, func(data chan<- interface{}) {
		atomic.AddInt32(&calls, 1)
		close(started)
		<-release
		for i := 0; i < 5; i++ {
			data <- i
		}
	})
	<-started
	stop <- true
	close(release)
	if got, want := drainDo(t, data), []interface{}{0, 1, 2, 3, 4}; !reflect.DeepEqual(got, want) {
		t.Fatalf("output = %v, want %v", got, want)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("callbacks = %d, want 1", got)
	}
	// The caller still owns stop after Do returns: closing it must not panic.
	close(stop)
}

func TestDoClosedStopFinishesActiveCallback(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var calls int32
	data, stop := Do(0, func(data chan<- interface{}) {
		atomic.AddInt32(&calls, 1)
		close(started)
		<-release
		data <- "last result"
	})
	<-started
	close(stop)
	close(release)
	if got, want := drainDo(t, data), []interface{}{"last result"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("output = %v, want %v", got, want)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("callbacks = %d, want 1", got)
	}
}

func TestDoStopInterruptsInterval(t *testing.T) {
	data, stop := Do(time.Hour, func(data chan<- interface{}) { data <- "first" })
	select {
	case value := <-data:
		if value != "first" {
			t.Fatalf("first result = %v", value)
		}
	case <-time.After(2 * time.Second):
		close(stop)
		t.Fatal("callback did not start")
	}
	stop <- true
	if got := drainDo(t, data); len(got) != 0 {
		t.Fatalf("unexpected callback after stop: %v", got)
	}
}

func TestDoCallbacksDoNotOverlap(t *testing.T) {
	var active, calls, overlap int32
	data, stop := Do(time.Millisecond, func(data chan<- interface{}) {
		if atomic.AddInt32(&active, 1) != 1 {
			atomic.StoreInt32(&overlap, 1)
		}
		call := atomic.AddInt32(&calls, 1)
		time.Sleep(5 * time.Millisecond)
		data <- call
		atomic.AddInt32(&active, -1)
	})
	for i := int32(1); i <= 3; i++ {
		select {
		case value := <-data:
			if value != i {
				t.Fatalf("result = %v, want %d", value, i)
			}
		case <-time.After(2 * time.Second):
			close(stop)
			t.Fatal("callback did not produce output")
		}
	}
	stop <- true
	drainDo(t, data)
	if atomic.LoadInt32(&overlap) != 0 || atomic.LoadInt32(&active) != 0 {
		t.Fatal("callbacks overlapped or remained active after shutdown")
	}
}
