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
	<param name="deviceNumber" value="UBK83CD7740-BCBB-43B0-B747-BE0152BE728E" />
	<param name="nameSilo" value="Silo1" />
	<param name="alcadaSilo" value="140" />
	<param name="diametreSilo" value="200" />
	<param name="ConeSilo" value="150" />
	<param name="OffsetDevice" value="0"  />
	<param name="Densitat" value="625.0" />
	<param name="Updated" value="1" />
</object>

<object id="siloChart2" role="svidget" data="http://{{.Host}}/images/silo.svg" type="image/svg+xml" width="400" height="400">
	<param name="hostname" value="{{.Host}}" />
	<param name="deviceNumber" value="UBK39F8600B-BA09-4885-B4B3-B77702BCC0DB" />
	<param name="nameSilo" value="Silo2" />
	<param name="alcadaSilo" value="140" />
	<param name="diametreSilo" value="200" />
	<param name="ConeSilo" value="150" />
	<param name="OffsetDevice" value="0"  />
	<param name="Densitat" value="625.0" />
	<param name="Updated" value="1" />
</object>

<object id="siloChart3" role="svidget" data="http://{{.Host}}/images/silo.svg" type="image/svg+xml" width="400" height="400">
	<param name="hostname" value="{{.Host}}" />
	<param name="deviceNumber" value="UBKD334F21E-3C23-11E5-8494-C3AD4A89321E" />
	<param name="nameSilo" value="Silo3" />
	<param name="alcadaSilo" value="140" />
	<param name="diametreSilo" value="200" />
	<param name="ConeSilo" value="150" />
	<param name="OffsetDevice" value="0"  />
	<param name="Densitat" value="625.0" />
	<param name="Updated" value="1" />
</object>


</body>

</html>
`
