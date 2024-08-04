package broadcast

import (
	"log"
	"net/http"
)



type Broadcast struct{
	Hub *Hub
}

func (b *Broadcast) HandleWebsocket(w http.ResponseWriter, r *http.Request){
	log.Println("Handling websocket connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Println(err)
		return
	}
	log.Println("Connection established")
	client := &Client{hub: b.Hub, conn: conn}
	client.hub.register <- client
}

func NewBroadcast() *Broadcast{
	log.Println("Creating new broadcast struct")
	return &Broadcast{
		Hub: NewHub(),
	}
}

