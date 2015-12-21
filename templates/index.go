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
		// Accesss
		$('#ubkName tspan').text('#Silo1');
		$('#ubkAvailable tspan').text('82.1%');
		$('#ubkCapacity tspan').html('26.00 m&#179;');
		$('#ubkContent tspan').text('154889 Kg');
		$('#ubkVolume tspan').html('22.10 m&#179;');
		$('#level10').hide();
			
</script>

</body>
</html>
`
