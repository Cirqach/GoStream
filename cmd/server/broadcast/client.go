package broadcast

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (c *Client) writePump(){
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select{
			case message, ok := <-c.send:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok{
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil{
					return
				}
				w.Write(message)
				if err := w.Close(); err != nil{
					return
				}
			
				}
		}
	}

				