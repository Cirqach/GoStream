package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


type Client struct{
	hub *Hub

	conn *websocket.Conn

	send chan []byte
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request){
	fmt.Println("Upgrading connection to websocket")
	conn, err := upgrader.Upgrade(w,r,nil); if err != nil{
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
}