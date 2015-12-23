// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// UBIKWA SLU
// Websocket server to provide realtime updates on user allower silos
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"text/template"
	"time"

	"github.com/davidraba/go-iot/templates"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	homeTempl = template.Must(template.New("").Parse(templates.Index))
	filename  string
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
	}
)

const (
	// Time allowed to write the file to the client.
	writeWait = 5 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 10 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 2 * time.Second

	// Maximum message size, here 1kB.
	maxMessageSize = 1024 * 1024

	maxLength = 1 << 20
)

var (
	MaxWorker = os.Getenv("MAX_WORKERS")
	MaxQueue  = os.Getenv("MAX_QUEUE")

	numprocs = runtime.NumCPU() * 2
	maxqueue = 20
)

// Connexion hub to handle realtime notifications based on device serialnumber
var h = hub{
	broadcast:  make(chan string),
	unicast:    make(chan DirectMessage),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
	sn:         "",
}

// Declare global jobQueue to handle running tasks
var jobQueue *Dispatcher

//TODO: Another hub that handle websockets created on longrunning tasks and
//enable user to retrieve realtime notifications of ongoing process carried on that process.
func main() {

	// EXECUTE and RUN go tool pprof -pdf ./ubk-wserver cpu.prof/cpu.pprof > test.pdf
	/*	cfg := profile.Config{
			CPUProfile:     true,
			MemProfile:     true,
			ProfilePath:    "./cpu.prof",
			NoShutdownHook: false, // do not hook SIGINT
		}
		p := profile.Start(&cfg)
		defer p.Stop()
	*/

	//init system signal channel, need for Ctrl+C and, kill and other system signal handling
	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	nu := NewUpdateController()
	// -----------------------------------------------------------------------
	// start dispatcher with NUMPROC workers (goroutines) and
	// jobsQueue channel size MAXQUEUE
	jobQueue = NewDispatcher(numprocs, maxqueue)
	go func() {
		// catch quit signal
		s := <-sChan
		log.Println("os.Signal", s, "received, finishing application...")
		// stop dispatcher
		jobQueue.Stop()
		os.Exit(1)
	}()
	go jobQueue.Run()

	// Run Hub Websocket connections -----------------------------------------
	go h.Run()

	//------------------------------------------------------------------------
	// Instantiate a new router
	r := gin.Default()

	r.StaticFile("/favicon.ico", "./images/favicon.ico")
	r.Static("/css", "./css")
	r.Static("/images", "./images")
	r.Static("/scripts", "./scripts")

	// -----------------------------------------------------------------------
	// ROUTE LIST
	// -----------------------------------------------------------------------
	// Index served from static file
	r.GET("/", func(c *gin.Context) { serveHome(c.Writer, c.Request) })

	// Webservice handler served
	r.GET("/ws", func(c *gin.Context) { serveWs(c.Writer, c.Request) })

	// Generate JSON with historical data, and compute distance according camera matrix,
	// and unpdate ContextBroker with computed distance to be displayed
	r.POST("/distance/update", nu.OnUpdateContext) // New distance derived from context updates

	// -----------------------------------------------------------------------
	// LAUNCH SERVER
	// -----------------------------------------------------------------------
	port := "8081"
	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}
	s := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	//--- RUN Server ---------------------------------------------------------
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
