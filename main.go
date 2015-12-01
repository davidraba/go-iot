// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// UBIKWA SLU
// Websocket server to provide realtime updates on user allower silos
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	addr      = flag.String("addr", ":80", "http service address")
	homeTempl = template.Must(template.New("").Parse(homeHTML))
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
	content:    "",
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
	s := &http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: r,
	}

	//--- RUN Server ---------------------------------------------------------
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

const homeHTML = `
<!doctype html>
<html>

<head>
    <title>SmartSilo</title>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<script language="Javascript" type="text/javascript" src="http://code.jquery.com/jquery-1.9.1.js"></script>
	<script src="//cdnjs.cloudflare.com/ajax/libs/d3/3.5.9/d3.min.js"></script>
	<script src="//cdnjs.cloudflare.com/ajax/libs/raphael/2.1.4/raphael-min.js"></script>
	<script src="//cdn.jsdelivr.net/justgage/1.0.1/justgage.min.js"></script>
	<script src="//cdnjs.cloudflare.com/ajax/libs/cubism/1.6.0/cubism.v1.min.js"></script>

    <script language='javascript'>
    /* 
    Array.push() appends to the end, and returns the new length
    Array.pop() removes the last and returns it
    Array.shift() removes the first and returns it
    Array.unshift() appends to the front and returns the new length
    */
    var data;
	$(function() {

		var gageValue = 0;
		var g;
		
		var g = new JustGage({
		    id: "g1",
		    value: getRandomInt(0, gageValue),
		    min: 0,
		    max: gageValue,
		    relativeGaugeSize: true,
			gaugeWidthScale: 0.2,
			title: "Silo 4",
			levelColors: [
		        "#00fff6",
		        "#ff00fc",
		        "#1200ff"
		      ],
			startAnimationTime: 1,
		    startAnimationType: "linear",
		    refreshAnimationTime: 1,
		    refreshAnimationType: "linear",				
		    donut: true
		});
		
		function updateGage(n) {
		 g.refresh(n);
		 gageValue = n;
		}

		function json_ws( on_msg ) {
		        if (!("WebSocket" in window)) {
		            alert("Use a browser supporting websockets");
		        }				
		              var conn = new WebSocket("ws://{{.Host}}/ws?device={{.Device}}");
		
		              conn.onclose = function(evt) {
		                  document.getElementById("fileData").textContent = 'Connection closed';
		              }
		
		              conn.onmessage = function(evt) {		 					
					try {
						data = JSON.parse(evt.data);
					}
		            catch (SyntaxError) {
		                console.log("Invalid data: " + evt.data);
		                return;
		            }
		            if (data)
		                on_msg(data);						
		              }
		
			window.onbeforeunload = function() {
		        			 conn.onclose = function() {};
		        			 conn.close();
			}
		}
	    var analog_keys = ['capacity', 'battery', 'temp'];	

		    (function() {
		        function make_realtime(key) {
		            var buf = [], callbacks = [];
		            return {
		                data: function(ts, val) {
		                    buf.push({ts: ts, val: val});
		                    callbacks = callbacks.reduce(function(result, cb) {
		                        if (!cb(buf))
		                            result.push(cb);
		                        return result
		                    }, []);
		                },
		                add_callback: function(cb) {
		                    callbacks.push(cb);
		                }
		            }
		        };
		        var realtime = {
		            capacity: make_realtime('capacity'),
		            battery: make_realtime('battery'),
		            temp: make_realtime('temp'),
		        };
				
		        /* This websocket sends homogenous messages in the form
		         * {timestamp: 1234567, analog: {capactity: 3.3, battery: 2.3, temp: 20}}
		         * where timestamp is a Unix timestamp
		         */
		        json_ws(function(data) {
		            analog_keys.map(function (key) {
		                realtime[key].data(data.timestamp, data.analog[key]);
						if (key == 'capacity') {
							updateGage(data.analog[key].toFixed(2))
						}
		            });
		        });
				
		        var context = cubism.context().serverDelay(5).step(1e3).size(1440);
		        
				var metric = function (key, title) {
		            var rt = realtime[key];
		            return context.metric(function (start, stop, step, callback) {
		                start = start.getTime();
		                stop = stop.getTime();
		                rt.add_callback(function(buf) {
		                    if (!(buf.length > 1 && 
		                          buf[buf.length - 1].ts > stop + step)) {
		                        // Not ready, wait for more data
		                        return false;
		                    }
		                    var r = d3.range(start, stop, step);
		                    /* Don't like using a linear search here, but I don't
		                     * know enough about cubism to really optimize. I had
		                     * assumed that once a timestamp was requested, it would
		                     * never be needed again so I could drop it. That doesn't
		                     * seem to be true!
		                     */
		                    var i = 0;
		                    var point = buf[i];
		                    callback(null, r.map(function (ts) {
		                        if (ts < point.ts) {
		                            // We have to drop points if no data is available
		                            return null;
		                        }
		                        for (; buf[i].ts < ts; i++);
		                        return buf[i].val;
		                    }));
		                    // opaque, but this tells the callback handler to
		                    // remove this function from its queue
		                    return true;
		                });
		            }, title);
		        };
		        ['top', 'bottom'].map(function (d) {
		            d3.select('#charts').append('div')
		                .attr('class', d + ' axis')
		                .call(context.axis().ticks(12).orient(d));
		        });
		        d3.select('#charts').append('div').attr('class', 'rule')
		            .call(context.rule());
		        charts = {
		            capacity: {
		                title: 'Capacity',
		                unit: 'V',
		                extent: [0, 100]
		            },
		            battery: {
		                title: 'Battery',
		                unit: 'mA',
		                extent: [3000, 5000]
		            },
		            temp: {
		                title: 'Temperature',
		                unit: '\u00b0C',
		                extent: [-20, 60]
		            }
		        };
		        Object.keys(charts).map(function (key) {
		            var cht = charts[key];
		            var num_fmt = d3.format('.3r');
		            d3.select('#charts')
		                .insert('div', '.bottom')
		                .datum(metric(key, cht.title))
		                .attr('class', 'horizon')
		                .call(context.horizon()
		                    .extent(cht.extent)
		                    .title(cht.title)
		                    .format(function (n) { 
		                        return num_fmt(n) + ' ' + cht.unit; 
		                    })
		                );
		        });
		        context.on('focus', function (i) {
		            if (i !== null) {
		                d3.selectAll('.value').style('right',
		                                             context.size() - i + 'px');
		            }
		            else {
		                d3.selectAll('.value').style('right', null)
		            }
		        });	
		    })();
		});
        </script>
	
    <style>
	@import url(//fonts.googleapis.com/css?family=Yanone+Kaffeesatz:400,700);
    body {
		font-family: "Helvetica Neue", Helvetica, sans-serif;
        text-align: center;
        padding: 0px;
        margin: 0px;
    }
    /* clearfix */

    .clear:before,
    .clear:after {
        content: "";
        display: table;
    }

    .clear:after {
        clear: both;
    }

    .clear {
        *zoom: 1;
    }

    .gauge {
        display: block;
        float: left;
    }

    #g1 {
        width: 20%;
    }

    #g2 {
        width: 60%;
    }

    #g3 {
        width: 20%;
    }
    </style>
</head>

<body>
    <pre id="fileData">{{.Data}}</pre>
	<div class="container">
		 <div id="g1" class="gauge" data-value="50" data-min="0" data-max="100" data-gaugeWidthScale="0.6"></div>
	</div>
    <div id='charts'></div>
</body>

</html>
`
