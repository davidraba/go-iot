package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/davidraba/go-iot/jobs"
	"github.com/davidraba/go-iot/models"
	"github.com/gorilla/websocket"
)

type PayloadCollection struct {
	Payloads []Payload `json:"payloads"`
}

type Payload struct {
	Data string `json:"data"`
}

//----------------------------------------------------------------------------
// Websocket
//----------------------------------------------------------------------------
func serveWs(w http.ResponseWriter, r *http.Request) {
	/*	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}*/
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		} else {
		}
		return
	}

	currentStatus := models.WebsocketData{
		Timestamp: time.Now().Unix() * 1000,
		Analog: models.AnalogData{
			Capacity:    0,
			Battery:     0,
			Temperature: 0,
		},
	}
	status_str := ""
	if b, err := json.Marshal(currentStatus); err == nil {
		status_str = string(b)
	}
	c := &client{
		send:   make(chan []byte, maxMessageSize),
		status: status_str,
		ws:     ws,
		sn:     r.FormValue("device"),
	}

	h.register <- c
	go c.writer(r.FormValue("device"))
	c.reader()
}

//----------------------------------------------------------------------------
// Home
//----------------------------------------------------------------------------
func serveHome(w http.ResponseWriter, r *http.Request) {
	var (
		p []byte
	)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if device := r.FormValue("device"); device != "" {
		p = []byte(device)
	} else {
		p = []byte("UBK2DD4EFCC-3C23-11E5-9AFC-1F78C59BFD03")
	}

	var v = struct {
		Host   string
		Data   string
		Device string
	}{
		r.Host,
		string(p),
		string(p),
	}
	homeTempl.Execute(w, &v)
}

//----------------------------------------------------------------------------
// Data Ingestion
//----------------------------------------------------------------------------
func dataIncomingHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	work := &jobs.OrionJob{Stream: b}
	// Push the work onto the queue.
	jobQueue.AddJob(work)

	w.WriteHeader(http.StatusOK)
}
