package channel

type Connections struct {
	clients      map[chan string]bool
	addClient    chan chan string
	removeClient chan chan string
	messages     chan string
}

func Init() *Connections {
	var hub = &Connections{
		clients:      make(map[chan string]bool),
		addClient:    make(chan (chan string)),
		removeClient: make(chan (chan string)),
		messages:     make(chan string),
	}

	go func() {
		for {
			select {
			case s := <-hub.addClient:
				hub.clients[s] = true
			case s := <-hub.removeClient:
				delete(hub.clients, s)
				if len(hub.clients) == 0 {
					return
				}
			case msg := <-hub.messages:
				for s, _ := range hub.clients {
					s <- msg
				}
			}
		}
	}()

	return hub
}

func (hub *Connections) AddReceiver() chan string {
	messageChannel := make(chan string)
	hub.addClient <- messageChannel
	return messageChannel
}

func (hub *Connections) CloseReceiver(ch chan string) {
	hub.removeClient <- ch
}

func (hub *Connections) AddTransmitter() chan<- string {
	return hub.messages
}
