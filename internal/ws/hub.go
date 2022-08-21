package ws

import (
	"encoding/json"
)

type Hub struct {
	clients        map[*Client]bool
	broadcast      chan TrackMessage
	register       chan *Client
	unregister     chan *Client
	dashboardState *Dashboard
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan TrackMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		dashboardState: &Dashboard{
			ActiveSessions: make(DashBoardSessions),
		},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if client.VisitorId != "" {
				h.dashboardState.Unregister(client.VisitorId)
			}
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.broadcastDashboardState()
		case message := <-h.broadcast:
			h.dashboardState.Track(message)
			h.broadcastDashboardState()
		}
	}
}

func (h Hub) broadcastDashboardState() {
	m, _ := json.Marshal(h.dashboardState.ActiveSessions)
	h.sendBroadcast(m)
}

func (h Hub) sendBroadcast(m []byte) {
	for client := range h.clients {
		//send active sessions stats only to dashboard users
		if client.dashboardClient {
			select {
			case client.send <- m:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}
