package gcurses

import (
	"bytes"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func (w *Writer) clear() error {
	for i := 0; i < w.lines; i++ {
		if _, err := w.Writer.Write([]byte("\033[0A")); err != nil {
			return err
		}
		if _, err := w.Writer.Write([]byte("\033[2K\r")); err != nil {
			return err
		}
	}
	return nil
}

type Writer struct {
	Writer  io.Writer
	running int32
	kill    int32
	buffer  bytes.Buffer
	mutex   *sync.Mutex
	lines   int
}

func New() *Writer {
	return &Writer{
		Writer: os.Stdout,
		mutex:  &sync.Mutex{},
	}
}

func (w *Writer) Flush() error {
	var lines int

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.buffer.Len() == 0 {
		return nil
	}
	if err := w.clear(); err != nil {
		return err
	}

	for _, b := range w.buffer.Bytes() {
		if b == '\n' {
			lines++
		}
	}
	w.lines = lines
	_, err := w.Writer.Write(w.buffer.Bytes())
	w.buffer.Reset()
	return err
}

func (w *Writer) Start() {
	atomic.StoreInt32(&w.kill, 0)
	if atomic.LoadInt32(&w.running) == 1 {
		return
	}
	atomic.StoreInt32(&w.running, 1)
	go w.Listen()
}

func (w *Writer) Stop() {
	_ = w.Flush()
	atomic.StoreInt32(&w.kill, 1)
}

func (w *Writer) Listen() {
	if atomic.LoadInt32(&w.running) == 0 {
		return
	}
	for atomic.LoadInt32(&w.kill) == 0 {
		w.Wait()
	}
	atomic.StoreInt32(&w.running, 0)
}

func (w *Writer) Wait() {
	time.Sleep(time.Millisecond * 20)
	_ = w.Flush()
}

func (w *Writer) Write(b []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.buffer.Write(b)
}
