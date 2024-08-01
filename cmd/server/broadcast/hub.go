package broadcast

import (
	"log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}


type Hub struct {
	clients map[*Client]bool
	Stream chan []byte
	register chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	log.Println("Creating new hub")
	return &Hub{
		Stream: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (h *Hub) Run(){
	log.Println("Starting hub")
	for {
		select{
			case client := <-h.register:
				log.Println("Registering client")
				h.clients[client] = true
			case client := <-h.unregister:
				log.Println("Unregistering client")
					delete(h.clients, client)
			case message := <-h.Stream:
				log.Println("Broadcasting message")
				h.Broadcast(message)
		}
	}
}

func (h *Hub) Broadcast(message []byte){
	for client := range h.clients{
		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}

