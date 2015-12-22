package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/davidraba/go-iot/jobs"
	"github.com/davidraba/go-iot/models"
	"github.com/davidraba/go-iot/util/conversions"
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
			Percentage:    0,
			Capacity:      0,
			WeightCurrent: 0,
			VolumeCurrent: 0,
		},
	}

	status_str := ""
	if b, err := json.Marshal(currentStatus); err == nil {
		status_str = string(b)
	}

	// Params POST:
	//   as: Altura Silo
	//   ds: Diametro Silo
	//   cs: Cono Silo
	//   d: Densidad material
	//   od : Offset distance. Distáncia desde el sensor al contenido(lleno)
	//   device: Numero de serie dispositivo
	// Sample:
	//   as=578&ds=340&cs=204&od=68&d=680.0&device=UBKD334F21E-3C23-11E5-8494-C3AD4A89321E
	as, _ := conversions.IntForValue(r.FormValue("as"))
	ds, _ := conversions.IntForValue(r.FormValue("ds"))
	cs, _ := conversions.IntForValue(r.FormValue("cs"))
	od, _ := conversions.IntForValue(r.FormValue("od"))
	dens, _ := conversions.Float64ForValue(r.FormValue("d"))

	fmt.Println(r.FormValue("as"), as, ds, cs, od, dens)
	q := models.SiloMeasureSimpleRequest{
		SiloHeightCm:         as,
		SiloDiameterCm:       ds,
		SiloHeightConeCm:     cs,
		SiloOffsetDistanceCm: od,
		ContentDensityKgm3:   dens,
	}

	c := &client{
		send:   make(chan models.SiloData, maxMessageSize),
		status: status_str,
		ws:     ws,
		sn:     r.FormValue("device"),
		config: q,
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
