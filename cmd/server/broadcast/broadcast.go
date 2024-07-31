package broadcast

import (
	"log"
	"net/http"
)



type Broadcast struct{
	Hub *Hub
}

func (b *Broadcast) handleWebsocket(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Println(err)
		return
	}
	client := &Client{hub: b.Hub, conn: conn}
	client.hub.register <- client
	go client.writePump()
}

func NewBroadcast() *Broadcast{
	return &Broadcast{
		Hub: NewHub(),
	}
}

