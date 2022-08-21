package ws

import (
	"encoding/json"
)

type Hub struct {
	clients    map[*Client]bool
	track      chan TrackMessage
	register   chan *Client
	unregister chan *Client
	dashboard  *Dashboard
	broadcast  chan *Dashboard
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		track:      make(chan TrackMessage),
		broadcast:  make(chan *Dashboard),
		dashboard: &Dashboard{
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
				h.dashboard.Unregister(client.VisitorId)
			}
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.broadcastDashboardState()
		case message := <-h.track:
			h.dashboard.Track(message)
			h.broadcastDashboardState()
		}
	}
}

func (h Hub) broadcastDashboardState() {
	m, _ := json.Marshal(h.dashboard.ActiveSessions)
	h.sendDashboardBroadcast(m)
}

func (h Hub) sendDashboardBroadcast(m []byte) {
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
