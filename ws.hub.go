package main

import (
	"encoding/json"
	"time"
)

type Hub struct {
	clients    map[string]*Client
	broadcast  chan WSRawMessage
	register   chan *Client
	unregister chan *Client
}

func newWSHub() *Hub {
	return &Hub{
		broadcast:  make(chan WSRawMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			socketID := client.userID

			go WSPingPong(*client, Ping, socketID, time.Second*5)

			h.clients[socketID] = client
		case client := <-h.unregister:
			if client, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
		case packetRaw := <-h.broadcast:
			var packet WSPacket
			err := json.Unmarshal(packetRaw.Data, &packet)
			if err != nil {
				continue
			}

			packet.Client = packetRaw.Client

			switch packet.Type {
			case RequestChannelAccess:
				RequestChannelAccessChan <- packet
			case GamePing:
				GamePingChan <- packet
			case GameInteraction:
				GameInteractionChan <- packet
			default:
				continue
			}
		}
	}
}

func (h *Hub) Broadcast(message string) {
	for client := range h.clients {
		select {
		case h.clients[client].send <- []byte(message):
		default:
			close(h.clients[client].send)
			delete(h.clients, client)
		}
	}
}
