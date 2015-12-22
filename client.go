package main

import (
	"encoding/json"
	"github.com/davidraba/go-iot/models"
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	ws     *websocket.Conn
	send   chan models.SiloData
	sn     string
	status string
	config models.SiloMeasureSimpleRequest
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
		case data, ok := <-c.send:
			if !ok { // Si no Ok, no està viu, tanca la connexió
				c.ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.ws.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return
				}
				return
			}
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))

			// Compute distance according client configuration
			c.config.Distancia = data.Distance
			volume := c.config.EvalDistance()

			t := time.Now()
			wsData := models.WebsocketData{
				Timestamp: t.Unix() * 1000,
				Analog: models.AnalogData{
					Percentage:    volume.ContentPerc,
					Capacity:      volume.SiloCapacityM3,
					WeightCurrent: volume.ContentWeightKg,
					VolumeCurrent: volume.ContentVolumeM3,
				},
			}

			message, _ := json.Marshal(wsData)
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
