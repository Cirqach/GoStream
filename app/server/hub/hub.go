package hub

import (
	"log"

	"github.com/gorilla/websocket"

	"github.com/Cirqach/GoStream/server/client"
)

type Hub struct {
	clients map[*client.Client]bool

	register chan *client.Client

	unregister chan *client.Client
}

func NewHub() *Hub {
	log.Println("Creating new hub")
	return &Hub{
		clients: make(map[*client.Client]bool),
		register: make(chan *client.Client),
		unregister: make(chan *client.Client),
	}
}

func (h *Hub) Run() {
	log.Println("Starting hub")
	for{
		select{
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok{
				delete(h.clients,client)
				close(client.send)
			}
		}
		
	}
}