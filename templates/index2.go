package templates

const Index2 = `
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
		    value: 0,
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
				
		        var context = cubism.context().serverDelay(0).step(1e3).size(1440);
		        
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
		                title: 'RSSI',
		                unit: 'dB',
		                extent: [-256, 0]
		            },
		            battery: {
		                title: 'Battery',
		                unit: 'mA',
		                extent: [0, 256]
		            },
		            temp: {
		                title: 'LQ',
		                unit: '\u00b0C',
		                extent: [0, 255]
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

<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-70941826-1', 'auto');
  ga('send', 'pageview');

</script>
</body>
</html>
`
