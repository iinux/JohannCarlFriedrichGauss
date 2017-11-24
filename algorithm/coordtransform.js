function Coordtransform() {
  //定义一些常量
  var x_PI = 3.14159265358979324 * 3000.0 / 180.0;
  var PI = 3.1415926535897932384626;
  var a = 6378245.0;
  var ee = 0.00669342162296594323;
  /**
   * 百度坐标系 (BD-09) 与 火星坐标系 (GCJ-02)的转换
   * 即 百度 转 谷歌、高德
   * @param bd_lon
   * @param bd_lat
   * @returns {*[]}
   */
  var bd09togcj02 = function bd09togcj02(bd_lon, bd_lat) {
    var bd_lon = +bd_lon;
    var bd_lat = +bd_lat;
    var x = bd_lon - 0.0065;
    var y = bd_lat - 0.006;
    var z = Math.sqrt(x * x + y * y) - 0.00002 * Math.sin(y * x_PI);
    var theta = Math.atan2(y, x) - 0.000003 * Math.cos(x * x_PI);
    var gg_lng = z * Math.cos(theta);
    var gg_lat = z * Math.sin(theta);
    return [gg_lng, gg_lat]
  };

  /**
   * 火星坐标系 (GCJ-02) 与百度坐标系 (BD-09) 的转换
   * 即谷歌、高德 转 百度
   * @param lng
   * @param lat
   * @returns {*[]}
   */
  var gcj02tobd09 = function gcj02tobd09(lng, lat) {
    var lat = +lat;
    var lng = +lng;
    var z = Math.sqrt(lng * lng + lat * lat) + 0.00002 * Math.sin(lat * x_PI);
    var theta = Math.atan2(lat, lng) + 0.000003 * Math.cos(lng * x_PI);
    var bd_lng = z * Math.cos(theta) + 0.0065;
    var bd_lat = z * Math.sin(theta) + 0.006;
    return [bd_lng, bd_lat]
  };

  /**
   * WGS84转GCj02
   * @param lng
   * @param lat
   * @returns {*[]}
   */
  var wgs84togcj02 = function wgs84togcj02(lng, lat) {
    var lat = +lat;
    var lng = +lng;
    if (out_of_china(lng, lat)) {
      return [lng, lat]
    } else {
      var dlat = transformlat(lng - 105.0, lat - 35.0);
      var dlng = transformlng(lng - 105.0, lat - 35.0);
      var radlat = lat / 180.0 * PI;
      var magic = Math.sin(radlat);
      magic = 1 - ee * magic * magic;
      var sqrtmagic = Math.sqrt(magic);
      dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * PI);
      dlng = (dlng * 180.0) / (a / sqrtmagic * Math.cos(radlat) * PI);
      var mglat = lat + dlat;
      var mglng = lng + dlng;
      return [mglng, mglat]
    }
  };

  /**
   * GCJ02 转换为 WGS84
   * @param lng
   * @param lat
   * @returns {*[]}
   */
  var gcj02towgs84 = function gcj02towgs84(lng, lat) {
    var lat = +lat;
    var lng = +lng;
    if (out_of_china(lng, lat)) {
      return [lng, lat]
    } else {
      var dlat = transformlat(lng - 105.0, lat - 35.0);
      var dlng = transformlng(lng - 105.0, lat - 35.0);
      var radlat = lat / 180.0 * PI;
      var magic = Math.sin(radlat);
      magic = 1 - ee * magic * magic;
      var sqrtmagic = Math.sqrt(magic);
      dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * PI);
      dlng = (dlng * 180.0) / (a / sqrtmagic * Math.cos(radlat) * PI);
      var mglat = lat + dlat;
      var mglng = lng + dlng;
      return [lng * 2 - mglng, lat * 2 - mglat]
    }
  };

  var transformlat = function transformlat(lng, lat) {
    var lat = +lat;
    var lng = +lng;
    var ret = -100.0 + 2.0 * lng + 3.0 * lat + 0.2 * lat * lat + 0.1 * lng * lat + 0.2 * Math.sqrt(Math.abs(lng));
    ret += (20.0 * Math.sin(6.0 * lng * PI) + 20.0 * Math.sin(2.0 * lng * PI)) * 2.0 / 3.0;
    ret += (20.0 * Math.sin(lat * PI) + 40.0 * Math.sin(lat / 3.0 * PI)) * 2.0 / 3.0;
    ret += (160.0 * Math.sin(lat / 12.0 * PI) + 320 * Math.sin(lat * PI / 30.0)) * 2.0 / 3.0;
    return ret
  };

  var transformlng = function transformlng(lng, lat) {
    var lat = +lat;
    var lng = +lng;
    var ret = 300.0 + lng + 2.0 * lat + 0.1 * lng * lng + 0.1 * lng * lat + 0.1 * Math.sqrt(Math.abs(lng));
    ret += (20.0 * Math.sin(6.0 * lng * PI) + 20.0 * Math.sin(2.0 * lng * PI)) * 2.0 / 3.0;
    ret += (20.0 * Math.sin(lng * PI) + 40.0 * Math.sin(lng / 3.0 * PI)) * 2.0 / 3.0;
    ret += (150.0 * Math.sin(lng / 12.0 * PI) + 300.0 * Math.sin(lng / 30.0 * PI)) * 2.0 / 3.0;
    return ret
  };

  /**
   * 判断是否在国内，不在国内则不做偏移
   * @param lng
   * @param lat
   * @returns {boolean}
   */
  var out_of_china = function out_of_china(lng, lat) {
    var lat = +lat;
    var lng = +lng;
    // 纬度3.86~53.55,经度73.66~135.05 
    return !(lng > 73.66 && lng < 135.05 && lat > 3.86 && lat < 53.55);
  };

  var lonLat2mercator = function lonLat2mercator(lon, lat){
    var x = (lon/180) * 20037508.34;
    if(lat > 85.05112){ lat = 85.05112;}
    if(lat < -85.05112){ lat = -85.05112;}
    lat = (Math.PI / 180.0) * lat;
    var y = 20037508.34 * Math.log(Math.tan(Math.PI / 4.0 + lat / 2.0)) / Math.PI;
    return [x,y];
  };
  
  var EARTH_RADIUS = 6378.137; // 地球半径
  function rad(d)
  {
     return d * Math.PI / 180.0;
  }
  var wgs_distance = function wgs_distance(lat1, lng1, lat2, lng2){
    var radLat1 = rad(lat1);
    var radLat2 = rad(lat2);
    var a = radLat1 - radLat2;
    var b = rad(lng1) - rad(lng2);
    var s = 2 * Math.asin(Math.sqrt(Math.pow(Math.sin(a/2),2) + Math.cos(radLat1)*Math.cos(radLat2)*Math.pow(Math.sin(b/2),2)));
    s = s * EARTH_RADIUS;
    s = Math.round(s * 10000) / 10000;
    return s;
  };

  // lon 经度，西经为负数
  // lat 纬度，南纬是负数
  var millerXY = function millerXY (lon, lat){
       var L = 6381372 * Math.PI * 2,     // 地球周长
           W = L,     // 平面展开后，x轴等于周长
           H = L / 2,     // y轴约等于周长一半
           mill = 2.3,      // 米勒投影中的一个常数，范围大约在正负2.3之间
           x = lon * Math.PI / 180,     // 将经度从度数转换为弧度
           y = lat * Math.PI / 180;     // 将纬度从度数转换为弧度
       // 这里是米勒投影的转换
       y = 1.25 * Math.log( Math.tan( 0.25 * Math.PI + 0.4 * y ) );
       // 这里将弧度转为实际距离
       x = ( W / 2 ) + ( W / (2 * Math.PI) ) * x;
       y = ( H / 2 ) - ( H / ( 2 * mill ) ) * y;
       // 转换结果的单位是公里
       // 可以根据此结果，算出在某个尺寸的画布上，各个点的坐标
       return [x, y];
       return {
            x : x,
            y : y
       };
  }

  var simple = function simple(lon, lat) {
    return [lon*85375.5033525, lat*111319.490821];
  }

  return {
    bd09togcj02: bd09togcj02,
    gcj02tobd09: gcj02tobd09,
    wgs84togcj02: wgs84togcj02,
    gcj02towgs84: gcj02towgs84,
    lonLat2mercator: lonLat2mercator,
    millerXY: millerXY,
    simple: simple,
    wgs_distance: wgs_distance
  }
}

coordtransform = new Coordtransform();
//百度经纬度坐标转国测局坐标
var bd09togcj02=coordtransform.bd09togcj02(116.404, 39.915);
//国测局坐标转百度经纬度坐标
var gcj02tobd09=coordtransform.gcj02tobd09(116.404, 39.915);
//wgs84转国测局坐标
var wgs84togcj02=coordtransform.wgs84togcj02(116.404, 39.915);
//国测局坐标转wgs84坐标
//var gcj02towgs84=coordtransform.gcj02towgs84(116.404, 39.915);
var gcj02towgs84=coordtransform.gcj02towgs84(116.37, 39.92);
var gcj02towgs84_2=coordtransform.gcj02towgs84(116.39, 39.92);

var mercator=coordtransform.lonLat2mercator(gcj02towgs84[0], gcj02towgs84[1]);
var mercator_2=coordtransform.lonLat2mercator(gcj02towgs84_2[0], gcj02towgs84_2[1]);
var mercator_dis=Math.sqrt(Math.pow(mercator[0]-mercator_2[0],2)+Math.pow(mercator[1]-mercator_2[1],2));
//console.log(mercator);
//console.log(mercator_2);
console.log(mercator_dis);

var miller=coordtransform.millerXY(gcj02towgs84[0], gcj02towgs84[1]);
var miller_2=coordtransform.millerXY(gcj02towgs84_2[0], gcj02towgs84_2[1]);
var miller_dis=Math.sqrt(Math.pow(miller[0]-miller_2[0],2)+Math.pow(miller[1]-miller_2[1],2));
console.log(miller_dis);

var simple=coordtransform.simple(gcj02towgs84[0], gcj02towgs84[1]);
var simple_2=coordtransform.simple(gcj02towgs84_2[0], gcj02towgs84_2[1]);
var simple_dis=Math.sqrt(Math.pow(simple[0]-simple_2[0],2)+Math.pow(simple[1]-simple_2[1],2));
console.log(simple_dis);

var wgs_dis=coordtransform.wgs_distance(gcj02towgs84[1], gcj02towgs84[0], gcj02towgs84_2[1], gcj02towgs84_2[0]);
console.log(wgs_dis);

//console.log(bd09togcj02);
//console.log(gcj02tobd09);
//console.log(wgs84togcj02);
console.log(gcj02towgs84);
console.log(gcj02towgs84_2);


//result
//bd09togcj02:   [ 116.39762729119315, 39.90865673957631 ]
//gcj02tobd09:   [ 116.41036949371029, 39.92133699351021 ]
//wgs84togcj02:  [ 116.41024449916938, 39.91640428150164 ]
//gcj02towgs84:  [ 116.39775550083061, 39.91359571849836 ]
