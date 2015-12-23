package templates

const Index = `
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="content-type" content="text/html; charset=UTF-8">
<link href="css/main.css" rel="stylesheet">
<script language="Javascript" type="text/javascript" src="http://code.jquery.com/jquery-1.9.1.js"></script>
</head>
<body>
<div id="svgInlineDiv">
</div>
<script type="text/javascript">
		function json_ws( on_msg ) {
		        if (!("WebSocket" in window)) {
		            alert("Use a browser supporting websockets");
		        }				
		              var conn = new WebSocket("ws://{{.Host}}/ws?device={{.Device}}&as=578&ds=340&cs=204&od=68&d=680.0");
		
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
		
		function updateSiloLevels( percent ) {
			if (percent > 90) { 
			      $('#level10').show();
			      $('#level9').show();
			      $('#level8').show();
			      $('#level7').show();
			      $('#level6').show();
			      $('#level5').show();
			      $('#level4').show();
			      $('#level3').hide();
			}
			else if (percent > 80) { 
			      $('#level10').hide();
			      $('#level9').show();
			      $('#level8').show();
			      $('#level7').show();
			      $('#level6').show();
			      $('#level5').show();
			      $('#level4').show();
			      $('#level3').hide();
			}
			else if (percent > 70) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').show();
			      $('#level7').show();
			      $('#level6').show();
			      $('#level5').show();
			      $('#level4').show();
			      $('#level3').hide();
			}
			else if (percent > 60) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').show();
			      $('#level6').show();
			      $('#level5').show();
			      $('#level4').show();
			      $('#level3').hide();
			}
			else if (percent > 50) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').hide();
			      $('#level6').show();
			      $('#level5').show();
			      $('#level4').show();
			      $('#level3').hide();
			}
			else if (percent > 40) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').hide();
			      $('#level6').hide();
			      $('#level5').show();
			      $('#level4').show();
			      $('#level3').hide();
			}
			else if (percent > 30) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').hide();
			      $('#level6').hide();
			      $('#level5').hide();
			      $('#level4').show();
			      $('#level3').hide();
			}
			else if (percent > 20) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').hide();
			      $('#level6').hide();
			      $('#level5').hide();
			      $('#level4').hide();
			      $('#level3').show();
			}
			else if (percent > 10) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').hide();
			      $('#level6').hide();
			      $('#level5').hide();
			      $('#level4').hide();
			      $('#level3').hide();
			      $('#level2').show();
			      $('#level1').show();				
			      $('#level0').show();
			}
			else if (percent > 5) { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').hide();
			      $('#level6').hide();
			      $('#level5').hide();
			      $('#level4').hide();
			      $('#level3').hide();
			      $('#level2').hide();
			      $('#level1').show();				
			      $('#level0').show();
			}
			else { 
			      $('#level10').hide();
			      $('#level9').hide();
			      $('#level8').hide();
			      $('#level7').hide();
			      $('#level6').hide();
			      $('#level5').hide();
			      $('#level4').hide();
			      $('#level3').hide();
			      $('#level2').hide();
			      $('#level1').hide();
			      $('#level0').show();
			}
			
		}
		
		/* This websocket sends homogenous messages in the form
         * {timestamp: 1234567, analog: {capactity: 3.3, battery: 2.3, temp: 20}}
         * where timestamp is a Unix timestamp
         */
        json_ws(function(data) {
			// Accesss
			$('#ubkName tspan').text('#Silo1');
			$('#ubkAvailable tspan').text( data.analog["percentage"] + '%');
			$('#ubkCapacity tspan').html(data.analog["percentage"] + 'm&#179;');
			$('#ubkContent tspan').text(data.analog["weight"] +'Kg');
			$('#ubkVolume tspan').html(data.analog["volume"] + '0m&#179;');
			updateSiloLevels(data.analog["percentage"]);
        });

		(function ($) {
		
		/**
		* @function
		* @property {object} jQuery plugin which runs handler function once specified element is inserted into the DOM
		* @param {function} handler A function to execute at the time when the element is inserted
		* @param {bool} shouldRunHandlerOnce Optional: if true, handler is unbound after its first invocation
		* @example $(selector).waitUntilExists(function);
		*/
		
		$.fn.waitUntilExists    = function (handler, shouldRunHandlerOnce, isChild) {
		    var found       = 'found';
		    var $this       = $(this.selector);
		    var $elements   = $this.not(function () { return $(this).data(found); }).each(handler).data(found, true);
		
		    if (!isChild)
		    {
		        (window.waitUntilExists_Intervals = window.waitUntilExists_Intervals || {})[this.selector] =
		            window.setInterval(function () { $this.waitUntilExists(handler, shouldRunHandlerOnce, true); }, 500)
		        ;
		    }
		    else if (shouldRunHandlerOnce && $elements.length)
		    {
		        window.clearInterval(window.waitUntilExists_Intervals[this.selector]);
		    }
		
		    return $this;
		}
		
		}(jQuery));

		function inlineSVG()
		{
		    var SVGFile="images/silo.svg"
		    var loadXML = new XMLHttpRequest;
		    function handler(){
		    if(loadXML.readyState == 4 && loadXML.status == 200)
		        svgInlineDiv.innerHTML=loadXML.responseText
		    }
		    if (loadXML != null){
		        loadXML.open("GET", SVGFile, true);
		        loadXML.onreadystatechange = handler;
		        loadXML.send();
		    }
		}
		
		function resizeSVG(){
			var svg = $('#ubkSilo')[0];

            var w = svg.getAttribute('width').replace('px', '');
            var h = svg.getAttribute('height').replace('px', '');

            svg.removeAttribute('width');
            svg.removeAttribute('height');

            svg.setAttribute('viewbox', '0 0 ' + w + ' ' + h);
            svg.setAttribute('preserveAspectRatio', 'xMinYMin meet')

            $(svg)
                .css('width', '20%')
                .css('height', '50%')
                .css('background-color', 'white');			
			}
		
        $(document).ready(function () {
			inlineSVG();
			$('#ubkSilo').waitUntilExists(resizeSVG);
        });
			
</script>

</body>
</html>
`
