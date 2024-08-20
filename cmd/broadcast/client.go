package broadcast

import (
	"time"

	"github.com/gorilla/websocket"
)

// TODO: add client token
// Client struct    allow access to control websockets connection
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)
