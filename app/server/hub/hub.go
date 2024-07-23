package client

import (
	"fmt"
)

type Hub struct {
	clients map[*Client]bool

	register chan *Client

	unregister chan *Client
}

func NewHub() *Hub {
	fmt.Println("New hub created")
	return &Hub{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	fmt.Println("Hub running")
	for{
		select{
		case client := <-h.register:
			fmt.Println("Client registered")
			h.clients[client] = true
		case client := <-h.unregister:
			fmt.Println("Client unregistered")
			if _, ok := h.clients[client]; ok{
				fmt.Println("Client exists")
				close(client.send)
				delete(h.clients,client)
			}
		}
		
	}
}