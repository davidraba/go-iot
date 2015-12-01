package main

type DirectMessage struct {
	sn   string
	mesg string
}

type hub struct {
	clients    map[*client]bool
	broadcast  chan string
	unicast    chan DirectMessage
	register   chan *client
	unregister chan *client
	content    string
	sn         string
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			c.send <- []byte(c.status)
			break

		case c := <-h.unregister:
			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
			}
			break

		case m := <-h.broadcast:
			h.content = m
			h.broadcastMessage()
			break

		case n := <-h.unicast:
			h.content = n.mesg
			h.sn = n.sn
			h.unicastMessage()
			break
		}

	}
}

func (h *hub) broadcastMessage() {
	for c := range h.clients {
		select {
		case c.send <- []byte(h.content):
			break

		// We can't reach the client
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}

func (h *hub) unicastMessage() {
	for c := range h.clients {
		if c.sn == h.sn {
			select {
			case c.send <- []byte(h.content):
				break

			// We can't reach the client
			default:
				close(c.send)
				delete(h.clients, c)
			}
		}
	}
}
