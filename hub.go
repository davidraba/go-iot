package main

import (
	"github.com/davidraba/go-iot/models"
)

type DirectMessage struct {
	sn   string
	data models.SiloData
}

type hub struct {
	clients    map[*client]bool
	broadcast  chan string
	unicast    chan DirectMessage
	register   chan *client
	unregister chan *client
	data       models.SiloData
	sn         string
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			c.send <- models.SiloData{}
			break

		case c := <-h.unregister:
			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
			}
			break

		case n := <-h.unicast:
			h.data = n.data
			h.sn = n.sn
			h.unicastMessage()
			break
		}

	}
}

func (h *hub) unicastMessage() {
	for c := range h.clients {
		if c.sn == h.sn {
			select {
			case c.send <- h.data:
				break

			// We can't reach the client
			default:
				close(c.send)
				delete(h.clients, c)
			}
		}
	}
}
