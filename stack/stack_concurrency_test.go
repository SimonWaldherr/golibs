package stack

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestSnapshotAndConversionPreserveReceiver(t *testing.T) {
	tests := []struct {
		name        string
		constructor func() *Stack
		popOrder    []interface{}
		insertion   []interface{}
	}{
		{"lifo", Lifo, []interface{}{3, 2, 1}, []interface{}{1, 2, 3}},
		{"fifo", Fifo, []interface{}{2, 3, 4}, []interface{}{2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.constructor()
			for _, value := range []int{1, 2, 3, 4} {
				s.Push(value)
			}
			s.Pop()
			values := s.Val()
			if !reflect.DeepEqual(values, tt.popOrder) {
				t.Fatalf("Val() = %v, want %v", values, tt.popOrder)
			}
			values[0] = "changed snapshot"
			fifo, lifo := s.ToFifo(), s.ToLifo()
			if got := fifo.Val(); !reflect.DeepEqual(got, tt.insertion) {
				t.Errorf("FIFO copy = %v, want %v", got, tt.insertion)
			}
			wantLifo := []interface{}{tt.insertion[2], tt.insertion[1], tt.insertion[0]}
			if got := lifo.Val(); !reflect.DeepEqual(got, wantLifo) {
				t.Errorf("LIFO copy = %v, want %v", got, wantLifo)
			}
			if got := s.Val(); !reflect.DeepEqual(got, tt.popOrder) {
				t.Errorf("receiver changed after snapshots/conversions: %v", got)
			}
			fifo.Pop()
			fifo.Push("changed copy")
			lifo.Pop()
			lifo.Push("changed copy")
			if got := s.Val(); !reflect.DeepEqual(got, tt.popOrder) {
				t.Errorf("copy mutations changed receiver: %v", got)
			}
			s.Push("changed receiver")
			if got := fifo.Val(); got[len(got)-1] != "changed copy" {
				t.Errorf("receiver mutation changed FIFO copy: %v", got)
			}
			if got := lifo.Val(); got[0] != "changed copy" {
				t.Errorf("receiver mutation changed LIFO copy: %v", got)
			}
		})
	}
}

func TestTryPopAndPeekDistinguishEmptyValues(t *testing.T) {
	for _, constructor := range []func() *Stack{Fifo, Lifo} {
		for _, value := range []interface{}{"", nil, false, 0} {
			s := constructor()
			s.Push(value)
			if got, ok := s.Peek(); !ok || got != value || s.Len() != 1 {
				t.Errorf("Peek() = %v, %v with Len %d; want %v, true and Len 1", got, ok, s.Len(), value)
			}
			if got, ok := s.TryPop(); !ok || got != value || s.Len() != 0 {
				t.Errorf("TryPop() = %v, %v with Len %d; want %v, true and Len 0", got, ok, s.Len(), value)
			}
			if got, ok := s.Peek(); ok || got != nil {
				t.Errorf("empty Peek() = %v, %v", got, ok)
			}
			if got, ok := s.TryPop(); ok || got != nil {
				t.Errorf("empty TryPop() = %v, %v", got, ok)
			}
			if got := s.Pop(); got != "" {
				t.Errorf("empty Pop() = %v, want empty string", got)
			}
			s.Push("reused")
			if got := s.Pop(); got != "reused" {
				t.Errorf("Pop() after reusing empty stack = %v", got)
			}
		}
	}
}

func TestIterSnapshotAllowsAbandonedIteration(t *testing.T) {
	s := Fifo()
	for _, value := range []int{1, 2, 3, 4} {
		s.Push(value)
	}
	s.Pop()
	iterator := s.Iter()
	if first := <-iterator; first != 2 {
		t.Fatalf("first iterated value = %v, want 2", first)
	}
	// Leave the iterator partly unread while mutating its original stack.
	done := make(chan struct{})
	go func() {
		s.Pop()
		s.Push(5)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("partly unread iterator blocked mutation")
	}
	var remaining []interface{}
	for value := range iterator {
		remaining = append(remaining, value)
	}
	if want := []interface{}{3, 4}; !reflect.DeepEqual(remaining, want) {
		t.Errorf("iterator changed after stack mutation: %v, want %v", remaining, want)
	}
	if want := []interface{}{3, 4, 5}; !reflect.DeepEqual(s.Val(), want) {
		t.Errorf("stack values = %v, want %v", s.Val(), want)
	}
}

func TestMarshalRemainingValuesAndRestoreOrdering(t *testing.T) {
	for name, constructor := range map[string]func() *Stack{"lifo": Lifo, "fifo": Fifo} {
		t.Run(name, func(t *testing.T) {
			s := constructor()
			for _, value := range []int{1, 2, 3} {
				s.Push(value)
			}
			s.Pop()
			data, err := s.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			var nodes []interface{}
			if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&nodes); err != nil {
				t.Fatalf("gob []interface{} compatibility: %v", err)
			}
			want := []interface{}{1, 2}
			if name == "fifo" {
				want = []interface{}{2, 3}
			}
			if !reflect.DeepEqual(nodes, want) {
				t.Errorf("encoded values = %v, want %v", nodes, want)
			}
			for targetName, targetConstructor := range map[string]func() *Stack{"lifo": Lifo, "fifo": Fifo} {
				target := targetConstructor()
				target.Push("old")
				target.Pop()
				target.Push("discarded")
				if err := target.Unmarshal(data); err != nil {
					t.Fatal(err)
				}
				if target.Len() != 2 {
					t.Errorf("%s restored length = %d, want 2", targetName, target.Len())
				}
				wantPop := []interface{}{want[0], want[1]}
				if targetName == "lifo" {
					wantPop[0], wantPop[1] = wantPop[1], wantPop[0]
				}
				for _, value := range wantPop {
					if got := target.Pop(); got != value {
						t.Errorf("%s restored Pop() = %v, want %v", targetName, got, value)
					}
				}
			}
		})
	}
}

func TestUnmarshalErrorPreservesContents(t *testing.T) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode([]interface{}{"replacement", "values"}); err != nil {
		t.Fatal(err)
	}
	for _, constructor := range []func() *Stack{Fifo, Lifo} {
		s := constructor()
		s.Push(1)
		s.Push(2)
		want := s.Val()
		for _, invalid := range [][]byte{nil, []byte("invalid gob"), buf.Bytes()[:buf.Len()-1]} {
			if err := s.Unmarshal(invalid); err == nil {
				t.Error("invalid Unmarshal succeeded")
			}
			if got := s.Val(); !reflect.DeepEqual(got, want) || s.Len() != len(want) {
				t.Errorf("failed Unmarshal changed contents: %v, want %v", got, want)
			}
		}
	}
}

func TestConcurrentStackOperations(t *testing.T) {
	const workers, perWorker = 4, 100
	for name, constructor := range map[string]func() *Stack{"lifo": Lifo, "fifo": Fifo} {
		t.Run(name, func(t *testing.T) {
			s := constructor()
			var writers sync.WaitGroup
			for worker := 0; worker < workers; worker++ {
				writers.Add(1)
				go func(worker int) {
					defer writers.Done()
					for i := 0; i < perWorker; i++ {
						s.Push(worker*perWorker + i)
						s.Len()
						s.Peek()
						s.Val()
						s.ToFifo()
						s.ToLifo()
						for range s.Iter() {
						}
						if _, err := s.Marshal(); err != nil {
							t.Errorf("Marshal: %v", err)
						}
					}
				}(worker)
			}
			writers.Wait()
			if got := s.Len(); got != workers*perWorker {
				t.Fatalf("concurrent pushes left %d values, want %d", got, workers*perWorker)
			}
			results := make(chan interface{}, workers*perWorker)
			var readers sync.WaitGroup
			for worker := 0; worker < workers; worker++ {
				readers.Add(1)
				go func() {
					defer readers.Done()
					for {
						value, ok := s.TryPop()
						if !ok {
							return
						}
						results <- value
					}
				}()
			}
			readers.Wait()
			close(results)
			seen := make(map[interface{}]bool)
			for value := range results {
				if seen[value] {
					t.Errorf("concurrent pops repeated value %v", value)
				}
				seen[value] = true
			}
			if len(seen) != workers*perWorker || s.Len() != 0 {
				t.Errorf("concurrent pops returned %d unique values and left %d", len(seen), s.Len())
			}
		})
	}
}

func TestUnsetConcurrentWithReadsAndWrites(t *testing.T) {
	s := Fifo()
	var workers sync.WaitGroup
	workers.Add(2)
	go func() {
		defer workers.Done()
		for i := 0; i < 1000; i++ {
			s.Unset()
		}
	}()
	go func() {
		defer workers.Done()
		for i := 0; i < 1000; i++ {
			s.Push(i)
			s.Pop()
			s.Val()
			s.ToFifo()
			s.ToLifo()
			s.Peek()
			s.Len()
		}
	}()
	workers.Wait()
	if s.Len() != -1 || s.Pop() != "" || s.Val() != nil {
		t.Error("Unset did not preserve the uninitialized state")
	}
}

func TestFifoSustainedQueue(t *testing.T) {
	queue := Fifo()
	const pending = 32
	for i := 0; i < pending; i++ {
		queue.Push(i)
	}
	for i := 0; i < 5000; i++ {
		value, ok := queue.TryPop()
		if !ok || value != i {
			t.Fatalf("TryPop() = %v, %v; want %d", value, ok, i)
		}
		queue.Push(i + pending)
	}
	if queue.Len() != pending {
		t.Fatalf("Len() = %d, want %d", queue.Len(), pending)
	}
	if len(queue.nodes) > 2048 {
		t.Fatalf("queue retained %d slots for %d pending values", len(queue.nodes), pending)
	}
	for i, value := range queue.Val() {
		if value != 5000+i {
			t.Fatalf("remaining value %d = %v, want %d", i, value, 5000+i)
		}
	}
}
