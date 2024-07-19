package client

import (
	"log"
	"net/http"
	"time"

	"github.com/Cirqach/GoStream/server/hub"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 1024
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

)

type Client struct{
	hub *hub.Hub

	conn *websocket.Conn

	send chan []byte
}

func serveWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w,r,nil); if err != nil{
		log.Println(err)
		return
	}
	log.Println("Upgrading connection to websocket")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
}