package channel

type Communication struct {
	receiver    map[chan interface{}]bool
	addReceiver chan chan interface{}
	rmReceiver  chan chan interface{}
	messages    chan interface{}
}

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

func (hub *Communication) AddReceiver() chan interface{} {
	messageChannel := make(chan interface{})
	hub.addReceiver <- messageChannel
	return messageChannel
}

func (hub *Communication) CloseReceiver(ch chan interface{}) int {
	hub.rmReceiver <- ch
	return hub.CountReceiver()
}

func (hub *Communication) CountReceiver() int {
	return len(hub.receiver)
}

func (hub *Communication) AddTransmitter() chan<- interface{} {
	return hub.messages
}
