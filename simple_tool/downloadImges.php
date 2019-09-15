<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2019/4/4
 * Time: 10:29
 */

define('DEBUG', false);
define('THREAD', false);
define('MULTI_CURL', false);
$htmlFile = 'action.txt';

$suffixes = ['gif', 'png', 'jpg', 'jpeg'];
$upCaseSuffixes = [];
foreach ($suffixes as $suffix) {
    $upCaseSuffixes[] = strtoupper($suffix);
}
$suffixes = array_merge($suffixes, $upCaseSuffixes);
$suffixesStr = implode('|', $suffixes);

$urlsRegex = '/data-src=\'(https?:\/\/.*?\.(' . $suffixesStr . '))\'/';

/**
 * @param $var
 */
function dd($var)
{
    var_dump($var);
    die(0);
}

function dump($var)
{
    var_dump($var);
}

function debug($var)
{
    if (DEBUG) {
        dump($var);
    }
}

if (class_exists(Thread::class)) {
    class AsyncDownload extends Thread
    {
        protected $url;

        public function __construct($url)
        {
            $this->url = $url;
        }

        public function run()
        {
            if ($this->url) {
                download($this->url);
            }
        }
    }
}

class MultiCurl
{
    protected $chs = [];
    protected $multiHandle = null;

    public function __construct()
    {
        $this->multiHandle = curl_multi_init();
    }

    public function addHandle($ch, $callback)
    {
        $this->chs[] = [$ch, $callback];
    }

    public function addUrl($url, $callback)
    {
        echo "addUrl: $url\n";
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_HEADER, 0);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_TIMEOUT, 10);
        $this->chs[] = [$ch, $callback];
    }

    public function run()
    {
        foreach ($this->chs as $ch) {
            curl_multi_add_handle($this->multiHandle, $ch[0]);
        }

        $running = null;
        // 执行批处理句柄
        do {
            sleep(1);
            curl_multi_exec($this->multiHandle, $running);
        } while ($running > 0);

        // 关闭全部句柄
        foreach ($this->chs as $ch) {
            curl_multi_remove_handle($this->multiHandle, $ch[0]);
        }
        curl_multi_close($this->multiHandle);

        foreach ($this->chs as $ch) {
            $content = curl_multi_getcontent($ch[0]);
            $ch[1]($content);
        }
    }
}

$mc = new MultiCurl();

$fileContent = file_get_contents($htmlFile);
preg_match_all($urlsRegex, $fileContent, $matches);
debug($matches);
foreach ($matches[1] as $url) {
    if (THREAD && class_exists(AsyncDownload::class)) {
        $thread = new AsyncDownload($url);
        $thread->start();
        // $thread->join();
    } else {
        download($url);
    }
}

if (MULTI_CURL) {
    $mc->run();
}

function download($url)
{
    global $suffixesStr, $mc;
    $fileRegex = '/^https?:\/\/.*\/(.*?\.(' . $suffixesStr . '))/';
    preg_match($fileRegex, $url, $matches);
    debug($matches);
    $fileName = $matches[1];
    echo "$url $fileName\n";
    if (file_exists($fileName)) {
        echo "exist\n";
    } else {
        echo "not exist, downloading...\n";

        $filePutContents = function ($content) use ($fileName){
            $ret = file_put_contents($fileName, $content);
            if ($ret === false) {
                dd(__LINE__);
            } else {
                echo "downloaded\n";
            }
        };

        if (!DEBUG) {
            if (MULTI_CURL) {
                $mc->addUrl($url, $filePutContents);
            } else {
                $content = file_get_contents($url);
                $filePutContents($content);
            }
        }

    }

}
