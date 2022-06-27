package gcurses

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func (w *Writer) clear() {
	for i := 0; i < w.lines; i++ {
		fmt.Fprint(w.Writer, "\033[0A")
		fmt.Fprint(w.Writer, "\033[2K\r")
	}
}

type Writer struct {
	Writer  io.Writer
	running bool
	kill    bool
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
	w.clear()

	for _, byte := range w.buffer.Bytes() {
		if byte == '\n' {
			lines++
		}
	}
	w.lines = lines
	_, err := w.Writer.Write(w.buffer.Bytes())
	w.buffer.Reset()
	return err
}

func (w *Writer) Start() {
	w.kill = false
	go w.Listen()
}

func (w *Writer) Stop() {
	w.Flush()
	w.kill = true
}

func (w *Writer) Listen() {
	if w.running {
		return
	}
	for !w.kill {
		select {
		default:
			w.Wait()
		}
	}
	w.running = false
	return
}

func (w *Writer) Wait() {
	time.Sleep(time.Millisecond * 20)
	w.Flush()
}

func (w *Writer) Write(b []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.buffer.Write(b)
}
