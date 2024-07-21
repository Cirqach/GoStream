package server

import (
	"log"

)

type Hub struct {
	clients map[*Client]bool

	register chan *Client

	unregister chan *Client
}

func NewHub() *Hub {
	log.Println("Creating new hub")
	return &Hub{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
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