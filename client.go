package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

type TrackMessage struct {
	Client  *Client
	Visitor string `json:"visitor"`
	ItemId  string `json:"item_id"`
	State   bool   `json:"state"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readMessage() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var m TrackMessage
		m.Client = c
		err = json.Unmarshal(message, &m)

		if err != nil {
			log.Printf("unmarsharll error: %v", err)
			break
		}
		c.hub.broadcast <- m
	}
}

func (c *Client) writeMessage() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			//c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			//// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			if err := w.Close(); err != nil {
				return
			}
		}
	}

}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	godotenv.Load()

	upgrader := websocket.Upgrader{

		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
			return strings.Contains(allowedOrigins, origin)
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.readMessage()
	go client.writeMessage()

}
