<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="initial-scale=1.0, user-scalable=no" />
	<style type="text/css">
		body, html,#allmap {width: 100%;height: 100%;overflow: hidden;margin:0;font-family:"微软雅黑";}
	</style>
	<script type="text/javascript" src="http://api.map.baidu.com/api?v=2.0&ak=V0WaQASP9HPLzaEFKkAkNkgG8Kk5EKVn"></script>
	<title>根据城市名设置地图中心点</title>
</head>
<body>
	<div id="allmap"></div>
</body>
</html>
<script type="text/javascript">
	// 百度地图API功能
	//var map = new BMap.Map("allmap");  // 创建Map实例
	var map = new BMap.Map("allmap",{minZoom:4,maxZoom:24}); 
	map.centerAndZoom("南京",15);      // 初始化地图,用城市名设置地图中心点
	map.enableScrollWheelZoom(true);
	map.setZoom(14);   
	map.enableScrollWheelZoom(true);
	var new_point = new BMap.Point(120.01, 30.22);
	var marker = new BMap.Marker(new_point);  // 创建标注
	map.addOverlay(marker);              // 将标注添加到地图中
	map.panTo(new_point);      
</script>

