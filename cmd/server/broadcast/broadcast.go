package broadcast

import (
	"log"
	"net/http"
)

type BroadcastEngine struct {
	Hub *Hub
}

func (b *BroadcastEngine) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling websocket connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Connection established")
	client := &Client{hub: b.Hub, conn: conn}
	client.hub.register <- client
}

func NewBroadcastEngine() *BroadcastEngine {
	log.Println("Creating new broadcast struct")
	return &BroadcastEngine{
		Hub: NewHub(),
	}
}
