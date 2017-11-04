<?php
    $website = 'Perorsoft Search Engine(PHP Version)';
    $fontBeautify = false;
	function post_curl($queryString)
	{
		$log = date('Y-m-d H:m:s').' '.$_SERVER['REMOTE_ADDR'].' '.$_SERVER['REMOTE_PORT'].' '.$queryString." <br />\n";
		$log = urldecode($log);
		file_put_contents('searchlog.txt', $log, FILE_APPEND);
		
		$url = "http://www.google.com.hk/search?".$queryString;
		$this_header = array(
			"Accept-Language: zh-CN,zh;q=0.8"
		);
		$ch = curl_init();
		curl_setopt($ch, CURLOPT_HTTPHEADER, $this_header);
		curl_setopt($ch, CURLOPT_URL, $url);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);//设置是将结果保存到字符串中还是输出到屏幕上，1表示将结果保存到字符串
		curl_setopt($ch, CURLOPT_HEADER, 0);//显示返回的Header区域内容
		//curl_setopt($ch, CURLOPT_BINARYTRANSFER, true) ;
		//curl_setopt($ch, CURLOPT_ENCODING, 'gzip,deflate');
		//curl_setopt($ch, CURLOPT_FOLLOWLOCATION,true);//使用自动跳转
		//curl_setopt($ch, CURLOPT_USERAGENT,"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)");
		//curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, 0); // 对认证证书来源的检查
		//curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, 1); // 从证书中检查SSL加密算法是否存在
		//curl_setopt($ch, CURLOPT_USERAGENT, $_SERVER['HTTP_USER_AGENT']); // 模拟用户使用的浏览器
		//curl_setopt($ch, CURLOPT_REFERER, $ref);
		//curl_setopt($ch, CURLOPT_COOKIEFILE,$GLOBALS['cookie_file']); // 读取上面所储存的Cookie信息
		//curl_setopt($ch, CURLOPT_COOKIEJAR, $GLOBALS['cookie_file']); // 存放Cookie信息的文件名称
		//curl_setopt($ch, CURLOPT_TIMEOUT, 30); // 设置超时限制防止死循环
		//if (curl_errno($curl)) {
		//	echo 'Errno'.curl_error($curl);
		//}
		
		//$post_data = array ("username" => "bob","key" => "12345");
		// post数据
		//curl_setopt($ch, CURLOPT_POST, 1);
		// post的变量
		//curl_setopt($ch, CURLOPT_POSTFIELDS, $post_data);
		
		$output = curl_exec($ch);
		//$info = curl_getinfo($ch);
		curl_close($ch);
		
		$output = str_replace('/url?q=', '', $output);
		$output = str_replace('search?', 'search.php?', $output);
		$output = str_replace('action="/search"', 'action="/search.php"', $output);
		
		$output = mb_convert_encoding($output, 'utf-8', 'gbk'); //加上这行
		
		//打印获得的数据
		echo $output;
		//$str = gzdecode($output);
	}
	if ($_SERVER['REQUEST_METHOD'] == "POST"){//echo $_POST['q'];}if (false){
		$my_curl = curl_init();    //初始化一个curl对象
		$q = $_POST['q'];
		$q = str_replace(' ', '+', $q);
		post_curl('q='.$q);
	}
?>

<?php
	if ($_SERVER['REQUEST_METHOD'] == "GET"){
	    if (count($_GET) >= 1){
	        post_curl($_SERVER['QUERY_STRING']);
	        return ;
	    }
?>

        $serverParams = $this->serverParams;
        // 如果有HTTP_X_WAP_PROFILE则一定是移动设备
        if (isset ($serverParams['HTTP_X_WAP_PROFILE'])) {
            return true;
        }
        // 如果via信息含有wap则一定是移动设备,部分服务商会屏蔽该信息
        if (isset ($serverParams['HTTP_VIA'])) {
            return stristr($serverParams['HTTP_VIA'], "wap") ? true : false;// 找不到为false,否则为true
        }
        // 判断手机发送的客户端标志,兼容性有待提高
        if (isset ($serverParams['HTTP_USER_AGENT'])) {
            $clientKeywords = array(
                'mobile',
                'nokia',
                'sony',
                'ericsson',
                'mot',
                'samsung',
                'htc',
                'sgh',
                'lg',
                'sharp',
                'sie-',
                'philips',
                'panasonic',
                'alcatel',
                'lenovo',
                'iphone',
                'ipod',
                'blackberry',
                'meizu',
                'android',
                'netfront',
                'symbian',
                'ucweb',
                'windowsce',
                'palm',
                'operamini',
                'operamobi',
                'openwave',
                'nexusone',
                'cldc',
                'midp',
                'wap'
            );
            // 从HTTP_USER_AGENT中查找手机浏览器的关键字
            if (preg_match("/(" . implode('|', $clientKeywords) . ")/i", strtolower($serverParams['HTTP_USER_AGENT']))) {
                return true;
            }
        }
        if (isset ($serverParams['HTTP_ACCEPT'])) { // 协议法，因为有可能不准确，放到最后判断
            $accept = $serverParams['HTTP_ACCEPT'];
            // 如果只支持wml并且不支持html那一定是移动设备
            // 如果支持wml和html但是wml在html之前则是移动设备
            if ((strpos($accept, 'vnd.wap.wml') !== false) &&
                (strpos($accept, 'text/html') === false || strpos($accept, 'vnd.wap.wml') < strpos($accept, 'text/html'))
            ) {
                return true;
            }
        }
        return false;

<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title><?php echo $website ?></title>
	<?php if ($fontBeautify) { ?>
	<link href="https://fonts.googleapis.com/css?family=Architects+Daughter" rel="stylesheet" type="text/css">
	<?php } ?>

	<style type="text/css">

	::selection { background-color: #E13300; color: white; }
	::-moz-selection { background-color: #E13300; color: white; }

	body {
		background-color: #fff;
		margin: 40px;
		font: 13px/20px normal Helvetica, Arial, sans-serif;
		color: #4F5155;
	}

	a {
		color: #003399;
		background-color: transparent;
		font-weight: normal;
	}

	h1 {
		color: #444;
		background-color: transparent;
		border-bottom: 1px solid #D0D0D0;
		font-size: 19px;
		font-weight: normal;
		<?php if ($fontBeautify) { ?>
		font-family: 'Architects Daughter', 'Helvetica Neue', Helvetica, Arial, serif;
		<?php } ?>
		margin: 0 0 14px 0;
		padding: 14px 15px 10px 15px;
	}

	code {
		font-family: Consolas, Monaco, Courier New, Courier, monospace;
		font-size: 12px;
		background-color: #f9f9f9;
		border: 1px solid #D0D0D0;
		color: #002166;
		display: block;
		margin: 14px 0 14px 0;
		padding: 12px 10px 12px 10px;
	}

	#body {
		margin: 0 15px 0 15px;
	}

	p.footer {
		text-align: right;
		font-size: 11px;
		border-top: 1px solid #D0D0D0;
		line-height: 32px;
		padding: 0 10px 0 10px;
		margin: 20px 0 0 0;
	}

	#container {
		margin: 10px;
		border: 1px solid #D0D0D0;
		box-shadow: 0 0 8px #D0D0D0;
	}
	</style>
</head>
<body>

<div id="container">
	<h1><?php echo $website ?></h1>

	<div id="body">
		<form method="post" enctype="multipart/form-data" id="container" style="text-align: center;padding: 180px 0px 180px 0px;">
			<input id="id_content" name="q" type="text" style="padding: 10px;" size="100">
			<input type="submit" value="Search" style="padding: 10px 40px 10px 40px;margin-left: 36px;">
		</form>
	</div>

	<p class="footer">@2016 Perorsoft</p>
</div>

</body>
</html>

<?php } ?>