package channel

import (
	"sync"
	"testing"
	"time"
)

func receive(t *testing.T, receiver <-chan interface{}) (interface{}, bool) {
	t.Helper()
	select {
	case value, ok := <-receiver:
		return value, ok
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for receiver")
		return nil, false
	}
}

func TestBroadcastPreservesOrder(t *testing.T) {
	hub := Init()
	defer hub.Close()
	first, second := hub.AddReceiver(), hub.AddReceiver()
	for i := 0; i < 100; i++ {
		hub.Send(i)
	}
	for i := 0; i < 100; i++ {
		for _, receiver := range []chan interface{}{first, second} {
			if value, ok := receive(t, receiver); !ok || value != i {
				t.Fatalf("received %v, %v; want %d", value, ok, i)
			}
		}
	}
}

func TestReceiverRemovalAndRegistration(t *testing.T) {
	hub := Init()
	defer hub.Close()
	old := hub.AddReceiver()
	if count := hub.CloseReceiver(old); count != 0 {
		t.Fatalf("receiver count = %d, want 0", count)
	}
	if _, ok := receive(t, old); ok {
		t.Fatal("removed receiver remains open")
	}
	hub.CloseReceiver(old)
	hub.SendTo("discarded", old)
	newReceiver := hub.AddReceiver()
	hub.Send("new receiver")
	if value, ok := receive(t, newReceiver); !ok || value != "new receiver" {
		t.Fatalf("received %v, %v", value, ok)
	}
}

func TestRemoveUnreadReceiver(t *testing.T) {
	hub := Init()
	defer hub.Close()
	unread, active := hub.AddReceiver(), hub.AddReceiver()
	for i := 0; i < 100; i++ {
		hub.Send(i)
	}
	if count := hub.CloseReceiver(unread); count != 1 {
		t.Fatalf("receiver count = %d, want 1", count)
	}
	for i := 0; i < 100; i++ {
		if value, ok := receive(t, active); !ok || value != i {
			t.Fatalf("received %v, %v; want %d", value, ok, i)
		}
	}
	// Removal may let a send already in progress complete; the channel must
	// nevertheless close without needing to consume the whole pending queue.
	for range unread {
	}
}

func TestCloseFromReceiver(t *testing.T) {
	hub := Init()
	defer hub.Close()
	receiver := hub.AddReceiver()
	done := make(chan struct{})
	go func() {
		<-receiver
		hub.CloseReceiver(receiver)
		close(done)
	}()
	hub.Send("message")
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("receiver removal blocked")
	}
	if hub.CountReceiver() != 0 {
		t.Fatal("receiver was not removed")
	}
}

func TestShutdown(t *testing.T) {
	for _, closeHub := range []struct {
		name string
		call func(*Communication)
	}{
		{"close", func(h *Communication) { h.Close() }},
		{"all", func(h *Communication) { h.CloseAll() }},
		{"transmitter", func(h *Communication) { h.CloseTransmitter(h.AddTransmitter()) }},
		{"raw transmitter", func(h *Communication) { close(h.AddTransmitter()); h.Close() }},
	} {
		t.Run(closeHub.name, func(t *testing.T) {
			hub := Init()
			receiver := hub.AddReceiver()
			for i := 0; i < 100; i++ {
				hub.Send(i)
			}
			closeHub.call(hub)
			hub.Close()
			hub.Send("after close")
			hub.SendTo("after close", receiver)
			if _, ok := receive(t, receiver); ok {
				t.Fatal("receiver remains open after shutdown")
			}
			if _, ok := receive(t, hub.AddReceiver()); ok {
				t.Fatal("receiver added after shutdown remains open")
			}
			if hub.CountReceiver() != 0 {
				t.Fatal("shutdown left registered receivers")
			}
		})
	}
}

func TestSendTo(t *testing.T) {
	hub := Init()
	defer hub.Close()
	receiver := hub.AddReceiver()
	hub.Send("broadcast")
	hub.SendTo("direct", receiver)
	if value, ok := receive(t, receiver); !ok || value != "broadcast" {
		t.Fatalf("received %v, %v before direct message", value, ok)
	}
	if value, ok := receive(t, receiver); !ok || value != "direct" {
		t.Fatalf("received %v, %v", value, ok)
	}
}

func TestConcurrentLifecycle(t *testing.T) {
	hub := Init()
	var wg sync.WaitGroup
	ready := make(chan struct{}, 8)
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				receiver := hub.AddReceiver()
				hub.Send(j)
				hub.SendTo(j, receiver)
				_ = hub.CountReceiver()
				hub.CloseReceiver(receiver)
				if j == 0 {
					ready <- struct{}{}
					<-start
				}
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 8; i++ {
			<-ready
		}
		close(start)
		hub.CloseAll()
	}()
	wg.Wait()
	hub.Close()
}
