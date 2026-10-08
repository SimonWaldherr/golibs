// Package channel provides a broadcast hub for communicating between goroutines.
package channel

import "sync"

// Communication broadcasts each message to its registered receivers. Messages
// are delivered in order to each receiver. Receivers should keep consuming
// messages to avoid accumulating a pending queue.
//
// Create a hub with Init and call Close when it is no longer needed. Methods
// are safe for concurrent use. Callers must not close receiver channels or send
// through raw transmitters after closing the hub. A Communication must not be
// copied after first use.
type Communication struct {
	lock         sync.RWMutex
	receiver     map[chan interface{}]*subscription
	messages     chan interface{}
	done         chan struct{}
	dispatchDone chan struct{}
	workers      sync.WaitGroup
	closed       bool
}

type subscription struct {
	output chan interface{}
	done   chan struct{}
	wake   chan struct{}
	lock   sync.Mutex
	queue  []interface{}
	head   int
	closed bool
}

type directedMessage struct {
	receiver chan interface{}
	value    interface{}
}

func (s *subscription) enqueue(message interface{}) {
	s.lock.Lock()
	if !s.closed {
		s.queue = append(s.queue, message)
		select {
		case s.wake <- struct{}{}:
		default:
		}
	}
	s.lock.Unlock()
}

func (s *subscription) stop() {
	s.lock.Lock()
	if !s.closed {
		s.closed = true
		s.queue = nil
		s.head = 0
		close(s.done)
	}
	s.lock.Unlock()
}

func (s *subscription) run() {
	defer close(s.output)
	for {
		s.lock.Lock()
		if s.closed {
			s.lock.Unlock()
			return
		}
		if s.head == len(s.queue) {
			s.lock.Unlock()
			select {
			case <-s.done:
				return
			case <-s.wake:
			}
			continue
		}
		message := s.queue[s.head]
		s.queue[s.head] = nil
		s.head++
		if s.head == len(s.queue) {
			s.queue = nil
			s.head = 0
		} else if s.head >= 1024 && s.head >= len(s.queue)/2 {
			s.queue = append([]interface{}(nil), s.queue[s.head:]...)
			s.head = 0
		}
		s.lock.Unlock()
		select {
		case <-s.done:
			return
		case s.output <- message:
		}
	}
}

// Init creates a new Communication.
func Init() *Communication {
	hub := &Communication{
		receiver:     make(map[chan interface{}]*subscription),
		messages:     make(chan interface{}, 31),
		done:         make(chan struct{}),
		dispatchDone: make(chan struct{}),
	}
	go hub.run()
	return hub
}

func (hub *Communication) run() {
	defer close(hub.dispatchDone)
	for {
		select {
		case <-hub.done:
			return
		case message, ok := <-hub.messages:
			if !ok {
				hub.shutdown()
				return
			}
			hub.lock.RLock()
			if directed, ok := message.(directedMessage); ok {
				if receiver, exists := hub.receiver[directed.receiver]; exists {
					receiver.enqueue(directed.value)
				}
			} else {
				for _, receiver := range hub.receiver {
					receiver.enqueue(message)
				}
			}
			hub.lock.RUnlock()
		}
	}
}

// AddReceiver registers a receiver. It returns a closed channel if the hub is
// closed. CloseReceiver or Close closes the returned channel.
func (hub *Communication) AddReceiver() chan interface{} {
	s := &subscription{
		output: make(chan interface{}),
		done:   make(chan struct{}),
		wake:   make(chan struct{}, 1),
	}
	hub.lock.Lock()
	defer hub.lock.Unlock()
	if hub.closed {
		close(s.output)
		return s.output
	}
	hub.receiver[s.output] = s
	hub.workers.Add(1)
	go func() {
		defer hub.workers.Done()
		s.run()
	}()
	return s.output
}

// CloseReceiver removes a receiver and discards its pending messages. The
// receiver channel closes when its worker exits. It returns the remaining
// receiver count, and is safe to call repeatedly or from a receiving goroutine.
func (hub *Communication) CloseReceiver(ch chan interface{}) int {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	if receiver, ok := hub.receiver[ch]; ok {
		delete(hub.receiver, ch)
		receiver.stop()
	}
	return len(hub.receiver)
}

// CountReceiver returns the current number of registered receivers.
func (hub *Communication) CountReceiver() int {
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	return len(hub.receiver)
}

// AddTransmitter returns the hub's shared input channel. Stop sending through
// this channel before calling Close or CloseTransmitter. Use Send for
// cancellation during shutdown. Prefer CloseTransmitter to closing the raw
// channel directly; a direct channel close requires all producers to stop first.
func (hub *Communication) AddTransmitter() chan<- interface{} {
	return hub.messages
}

func (hub *Communication) shutdown() {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	if hub.closed {
		return
	}
	hub.closed = true
	close(hub.done)
	for ch, receiver := range hub.receiver {
		delete(hub.receiver, ch)
		receiver.stop()
	}
}

// Close stops the hub and all receivers, discarding pending messages. It waits
// for internal workers to exit and is safe to call repeatedly.
func (hub *Communication) Close() {
	hub.shutdown()
	<-hub.dispatchDone
	hub.workers.Wait()
}

// CloseTransmitter closes the hub when given a transmitter returned by
// AddTransmitter. All transmitters share the same input channel.
func (hub *Communication) CloseTransmitter(ch chan<- interface{}) {
	if ch == hub.messages {
		hub.Close()
	} else {
		close(ch)
	}
}

// CloseAll is an alias for Close.
func (hub *Communication) CloseAll() {
	hub.Close()
}

// Send queues a message for all receivers, or returns if the hub is closed.
func (hub *Communication) Send(message interface{}) {
	select {
	case <-hub.done:
		return
	default:
	}
	select {
	case <-hub.done:
	case hub.messages <- message:
	}
}

// SendTo queues a message for a registered receiver. Messages for removed or
// unknown receivers are discarded.
func (hub *Communication) SendTo(message interface{}, receiver chan interface{}) {
	hub.lock.RLock()
	_, ok := hub.receiver[receiver]
	closed := hub.closed
	hub.lock.RUnlock()
	if closed || !ok {
		return
	}
	hub.Send(directedMessage{receiver: receiver, value: message})
}
