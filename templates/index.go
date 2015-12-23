package templates

const Index = `
<!DOCTYPE html>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
 <meta http-equiv="content-type" content="text/html; charset=UTF-8">
	<script type="application/javascript" src="http://{{.Host}}/scripts/svidget.min.js"></script>
	<script type="application/javascript" src="http://{{.Host}}/scripts/snap.svg-min.js"></script>
</head>
<body>

<object id="siloChart1" role="svidget" data="http://{{.Host}}/images/silo.svg" type="image/svg+xml" width="400" height="400">
	<param name="hostname" value="{{.Host}}" />
	<param name="deviceNumber" value="UBKD334F21E-3C23-11E5-8494-C3AD4A89321E" />
	<param name="nameSilo" value="Silo1" />
	<param name="alcadaSilo" value="578" />
	<param name="diametreSilo" value="340" />
	<param name="ConeSilo" value="204" />
	<param name="OffsetDevice" value="68"  />
	<param name="Densitat" value="680.0" />
	<param name="Updated" value="1" />
</object>

<object id="siloChart2" role="svidget" data="http://{{.Host}}/images/silo.svg" type="image/svg+xml" width="400" height="400">
	<param name="hostname" value="{{.Host}}" />
	<param name="deviceNumber" value="UBKD334F21E-3C23-11E5-8494-C3AD4A89321E" />
	<param name="nameSilo" value="Silo2" />
	<param name="alcadaSilo" value="878" />
	<param name="diametreSilo" value="340" />
	<param name="ConeSilo" value="204" />
	<param name="OffsetDevice" value="68"  />
	<param name="Densitat" value="680.0" />
	<param name="Updated" value="1" />
</object>

<object id="siloChart3" role="svidget" data="http://{{.Host}}/images/silo.svg" type="image/svg+xml" width="400" height="400">
	<param name="hostname" value="{{.Host}}" />
	<param name="deviceNumber" value="UBKD334F21E-3C23-11E5-8494-C3AD4A89321E" />
	<param name="nameSilo" value="Silo3" />
	<param name="alcadaSilo" value="678" />
	<param name="diametreSilo" value="440" />
	<param name="ConeSilo" value="204" />
	<param name="OffsetDevice" value="68"  />
	<param name="Densitat" value="680.0" />
	<param name="Updated" value="1" />
</object>

<object id="siloChart4" role="svidget" data="http://{{.Host}}/images/silo.svg" type="image/svg+xml" width="400" height="400">
	<param name="hostname" value="{{.Host}}" />
	<param name="deviceNumber" value="UBKD334F21E-3C23-11E5-8494-C3AD4A89321E" />
	<param name="nameSilo" value="Silo4" />
	<param name="alcadaSilo" value="278" />
	<param name="diametreSilo" value="140" />
	<param name="ConeSilo" value="204" />
	<param name="OffsetDevice" value="58"  />
	<param name="Densitat" value="480.0" />
	<param name="Updated" value="1" />
</object>
</body>
</html>
`
