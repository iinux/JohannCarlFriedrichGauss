<?php
//return phpinfo();
$website = 'shell';

function run_cmd()
{//system exec passthru ``
    $cmd = $_POST['q'];
    /* Add redirection so we can get stderr. */
    $handle = popen($cmd . ' 2>&1', 'r');
    //echo $handle.'<br />';
    echo "'$handle'; " . gettype($handle) . "<br />";
    return;
    while ($read = fread($handle, 2048)) {
        $read = str_replace('<', '&lt;', $read);
        $read = str_replace('>', '&gt;', $read);
        $read = str_replace('\r\n', '<br />', $read);
        $read = str_replace('\n', '<br />', $read);
        echo $read . '<br />';
    }
    pclose($handle);
}

function get_zip_originalsize($filename, $path)
{
    //先判断待解压的文件是否存在
    if (!file_exists($filename)) {
        die("文件 $filename 不存在！");
    }
    $starttime = explode(' ', microtime()); //解压开始的时间

    //将文件名和路径转成windows系统默认的gb2312编码，否则将会读取不到
    $filename = iconv("utf-8", "gb2312", $filename);
    $path = iconv("utf-8", "gb2312", $path);
    //打开压缩包
    $resource = zip_open($filename);
    $i = 1;
    //遍历读取压缩包里面的一个个文件
    while ($dir_resource = zip_read($resource)) {
        //如果能打开则继续
        if (zip_entry_open($resource, $dir_resource)) {
            //获取当前项目的名称,即压缩包里面当前对应的文件名
            $file_name = $path . zip_entry_name($dir_resource);
            //以最后一个“/”分割,再用字符串截取出路径部分
            $file_path = substr($file_name, 0, strrpos($file_name, "/"));
            //如果路径不存在，则创建一个目录，true表示可以创建多级目录
            if (!is_dir($file_path)) {
                mkdir($file_path, 0777, true);
            }
            //如果不是目录，则写入文件
            if (!is_dir($file_name)) {
                //读取这个文件
                $file_size = zip_entry_filesize($dir_resource);
                //最大读取6M，如果文件过大，跳过解压，继续下一个
                if ($file_size < (1024 * 1024 * 6)) {
                    $file_content = zip_entry_read($dir_resource, $file_size);
                    file_put_contents($file_name, $file_content);
                } else {
                    echo "<p> " . $i++ . " 此文件已被跳过，原因：文件过大， -> " . iconv("gb2312", "utf-8", $file_name) . " </p>";
                }
            }
            //关闭当前
            zip_entry_close($dir_resource);
        }
    }
    //关闭压缩包
    zip_close($resource);
    $endtime = explode(' ', microtime()); //解压结束的时间
    $thistime = $endtime[0] + $endtime[1] - ($starttime[0] + $starttime[1]);
    $thistime = round($thistime, 3); //保留3为小数
    echo "<p>解压完毕！，本次解压花费：$thistime 秒。</p>";
}

function my_zip()
{
//需开启配置 php_zip.dll
//phpinfo();
    header("Content-type:text/html;charset=utf-8");
    $size = get_zip_originalsize('Flat-UI-master.zip', 'zhangtemp/');
}

function run_cmd2()
{
    $phar = new PharData('song.tar.gz');
    //路径 要解压的文件 是否覆盖
    $phar->extractTo('c:/tmp', null, true);
}

function test__server_array()
{
    #测试网址:     http://localhost/blog/testurl.php?id=5

    //获取域名或主机地址
    echo $_SERVER['HTTP_HOST'] . "<br>"; #localhost

    //获取网页地址
    echo $_SERVER['PHP_SELF'] . "<br>"; #/blog/testurl.php

    //获取网址参数
    echo $_SERVER["QUERY_STRING"] . "<br>"; #id=5

    try {
        //获取用户代理
        if (array_key_exists('HTTP_REFERER', $_SERVER))
            echo $_SERVER['HTTP_REFERER'] . "<br>";
    } catch (Exception $e) {
    }
    //获取完整的url
    echo 'http://' . $_SERVER['HTTP_HOST'] . $_SERVER['REQUEST_URI'] . "<br>";
    echo 'http://' . $_SERVER['HTTP_HOST'] . $_SERVER['PHP_SELF'] . '?' . $_SERVER['QUERY_STRING'] . "<br>";
    #http://localhost/blog/testurl.php?id=5

    //包含端口号的完整url
    echo 'http://' . $_SERVER['SERVER_NAME'] . ':' . $_SERVER["SERVER_PORT"] . $_SERVER["REQUEST_URI"] . "<br>";
    #http://localhost:80/blog/testurl.php?id=5

    //只取路径
    $url = 'http://' . $_SERVER['SERVER_NAME'] . $_SERVER["REQUEST_URI"];
    echo dirname($url);
    #http://localhost/blog
}

function preg_match_test()
{
    //print(2)+3;
    //echo '1'.print(2)+3;
    //echo "1".(print 2)+3;
    $str = "";
    //preg_match("/^[".chr(0xa1)."-".chr(0xff)."A-Za-z0-9_]+$/",$str)//gb2312
    if (preg_match("/^[\x{4e00}-\x{9fa5}]+$/u", $str)) {
        print("该字符串全部是中文");
    } else {
        print("该字符串不全部是中文");
    }
}

?>
<?php
if (False && $_SERVER['REQUEST_METHOD'] == "GET") {
    ?>

    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="utf-8">
        <title><?php echo $website ?></title>

        <style type="text/css">

            ::selection {
                background-color: #E13300;
                color: white;
            }

            ::-moz-selection {
                background-color: #E13300;
                color: white;
            }

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
            <form method="post" enctype="multipart/form-data" id="container"
                  style="text-align: center;padding: 180px 0px 180px 0px;">
                <input id="id_content" name="q" type="text" style="padding: 10px;" size="100">
                <input type="submit" value="Enter" style="padding: 10px 40px 10px 40px;margin-left: 36px;">
            </form>
        </div>

        <p class="footer">@2016 Perorsoft</p>
    </div>

    </body>
    </html>

<?php } else {
    my_zip();
} ?>