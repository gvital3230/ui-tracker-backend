// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
)

type Hub struct {
	clients    map[*Client]map[string]bool
	broadcast  chan TrackMessage
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]map[string]bool),
		broadcast:  make(chan TrackMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = make(map[string]bool)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			if message.State {
				h.clients[message.Client][message.ItemId] = true
			} else {
				if _, ok := h.clients[message.Client][message.ItemId]; ok {
					delete(h.clients[message.Client], message.ItemId)
				}
			}
			fmt.Println(h.clients[message.Client])
			m, _ := json.Marshal(h.clients)
			fmt.Println(m)
			for client := range h.clients {
				select {
				case client.send <- m:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
