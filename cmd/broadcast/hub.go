package broadcast

import (
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub struct    allow access to control websockets clients
type Hub struct {
	clients    map[*Client]bool
	Stream     chan []byte
	Register   chan *Client
	unregister chan *Client
}

// NewHub function    create new Hub
func NewHub() *Hub {
	logger.LogMessage(logger.GetFuncName(0), "Creating new hub")
	return &Hub{
		Stream:     make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run method    run hub, for registering clients and sending messages to them
func (h *Hub) Run() {
	logger.LogMessage(logger.GetFuncName(0), "Starting hub")
	for {
		select {
		case client := <-h.Register:
			logger.LogMessage(logger.GetFuncName(0), "Registering client")
			h.clients[client] = true
		case client := <-h.unregister:
			logger.LogMessage(logger.GetFuncName(0), "Deleting client from hub")
			delete(h.clients, client)
		case message := <-h.Stream:
			logger.LogMessage(logger.GetFuncName(0), "Broadcasting message to clients")
			h.Broadcast(message)
		}
	}
}

// Broadcast method    send message to all clients
func (h *Hub) Broadcast(message []byte) {
	logger.LogMessage(logger.GetFuncName(0), "Broadcasting message to clients: "+string(message))
	for client := range h.clients {
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

// TODO: rewrite with tokens
// client  
// func (h *Hub) SendToClient(client *Client, message []byte) {
// 	logger.LogMessage(logger.GetFuncName(0), "Sending message to client: "+client.Id)
// 	err := client.Conn.WriteMessage(websocket.TextMessage, message)
// 	if err != nil {
// 		logger.LogError(logger.GetFuncName(0), err.Error())
// 	}
// }

// TODO: rewrite with tokens
// FindClient method    find client by token and return pointer to it
// func (h *Hub) FindClient(clientID string) *Client {
// 	for client := range h.clients {
// 		if client.Id == clientID {
// 			logger.LogMessage(logger.GetFuncName(0), "Found client with ID: "+clientID)
// 			return client
// 		}
// 	}
// 	logger.LogMessage(logger.GetFuncName(0), "Could not find client with ID: "+clientID)
// 	return nil
// }
