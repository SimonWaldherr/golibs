package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := Client(3 * time.Second)
	if client.Transport != Transporter || client.Timeout != 3*time.Second {
		t.Fatalf("unexpected client configuration: %+v", client)
	}
}

func TestGetString(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != stdhttp.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		w.WriteHeader(stdhttp.StatusNotFound)
		_, _ = io.WriteString(w, "missing 😃")
	}))
	defer server.Close()

	body, err := GetString(server.URL)
	if err != nil || body != "missing 😃" {
		t.Fatalf("GetString() = %q, %v", body, err)
	}
}

func TestGetStringContextCancellation(t *testing.T) {
	for _, streamResponse := range []bool{false, true} {
		name := "before headers"
		if streamResponse {
			name = "during body read"
		}
		t.Run(name, func(t *testing.T) {
			testGetStringContextCancellation(t, streamResponse)
		})
	}
}

func testGetStringContextCancellation(t *testing.T, streamResponse bool) {
	t.Helper()
	started := make(chan struct{})
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if streamResponse {
			_, _ = io.WriteString(w, "partial")
			w.(stdhttp.Flusher).Flush()
		}
		close(started)
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := GetStringContext(ctx, server.URL)
		result <- err
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("request did not reach the server")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("request did not stop after cancellation")
	}
}

func TestGetStringReadError(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Length", "100")
		_, _ = io.WriteString(w, "short")
	}))
	defer server.Close()

	body, err := GetString(server.URL)
	if body != "" || !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("GetString() = %q, %v; want an empty body and unexpected EOF", body, err)
	}
}

func TestGetStringInvalidURL(t *testing.T) {
	if _, err := GetString(":invalid"); err == nil {
		t.Fatal("invalid URL did not return an error")
	}
}

func TestGetStringConcurrent(t *testing.T) {
	var requests int32
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		atomic.AddInt32(&requests, 1)
		_, _ = io.WriteString(w, "hello")
	}))
	defer server.Close()

	const count = 32
	var wg sync.WaitGroup
	failures := make(chan error, count)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body, err := GetString(server.URL)
			if err != nil || body != "hello" {
				failures <- fmt.Errorf("GetString() = %q, %v", body, err)
			}
		}()
	}
	wg.Wait()
	close(failures)
	for err := range failures {
		t.Error(err)
	}
	if got := atomic.LoadInt32(&requests); got != count {
		t.Errorf("requests = %d, want %d", got, count)
	}
}
