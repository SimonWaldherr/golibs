// channel is a simple communication hub for go routines
package channel

type Communication struct {
	receiver    map[chan interface{}]bool
	addReceiver chan chan interface{}
	rmReceiver  chan chan interface{}
	messages    chan interface{}
}

// Init creates a new Communication
func Init() *Communication {
	var hub = &Communication{
		receiver:    make(map[chan interface{}]bool),
		addReceiver: make(chan (chan interface{})),
		rmReceiver:  make(chan (chan interface{})),
		messages:    make(chan interface{}, 31),
	}

	go func() {
		for {
			select {
			case s := <-hub.addReceiver:
				hub.receiver[s] = true
			case s := <-hub.rmReceiver:
				delete(hub.receiver, s)
				if len(hub.receiver) == 0 {
					return
				}
			case msg := <-hub.messages:
				for rec := range hub.receiver {
					go func(message interface{}, receiver chan interface{}) {
						receiver <- message
					}(msg, rec)
				}
			}
		}
	}()

	return hub
}

// AddReceiver adds a new receiver to the hub
func (hub *Communication) AddReceiver() chan interface{} {
	messageChannel := make(chan interface{})
	hub.addReceiver <- messageChannel
	return messageChannel
}

// CloseReceiver closes a receiver
func (hub *Communication) CloseReceiver(ch chan interface{}) int {
	hub.rmReceiver <- ch
	return hub.CountReceiver()
}

// CountReceiver returns the number of receivers
func (hub *Communication) CountReceiver() int {
	return len(hub.receiver)
}

// AddTransmitter adds a new transmitter to the hub
func (hub *Communication) AddTransmitter() chan<- interface{} {
	return hub.messages
}

// Close closes the hub
func (hub *Communication) Close() {
	close(hub.messages)
}

// CloseTransmitter closes a transmitter
func (hub *Communication) CloseTransmitter(ch chan<- interface{}) {
	close(ch)
}

// CloseAll closes all receivers and transmitters
func (hub *Communication) CloseAll() {
	for rec := range hub.receiver {
		hub.CloseReceiver(rec)
	}
	hub.Close()
}

// Send sends a message to all receivers
func (hub *Communication) Send(message interface{}) {
	hub.messages <- message
}

// SendTo sends a message to a specific receiver
func (hub *Communication) SendTo(message interface{}, receiver chan interface{}) {
	receiver <- message
}
