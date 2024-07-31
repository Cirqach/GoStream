package broadcast

import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}


type Hub struct {
	clients map[*Client]bool
	stream chan []byte
	register chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		stream: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) Run(){
	for {
		select{
			case client := <-h.register:
				h.clients[client] = true
			case client := <-h.unregister:
					delete(h.clients, client)
			case message := <-h.stream:
				h.broadcast(message)
		}
	}
}

func (h *Hub) broadcast(m []byte){
	for client := range h.clients{
		select{
			case client.send <- m:
			default:
				close(client.send)
				delete(h.clients, client)
		}
	}
}