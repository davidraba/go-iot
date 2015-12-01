package main

import (
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	ws     *websocket.Conn
	send   chan []byte
	sn     string
	status string
}

func (c *client) reader() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *client) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	c.status = string(message)
	return c.ws.WriteMessage(mt, message)
}

func (c *client) writer(serialnumber string) {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok { // Si no Ok, no està viu, tanca la connexió
				c.ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.ws.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return
				}
				return
			}
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-pingTicker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
