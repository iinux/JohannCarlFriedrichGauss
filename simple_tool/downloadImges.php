<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2019/4/4
 * Time: 10:29
 */

define('DEBUG', false);
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

$fileContent = file_get_contents($htmlFile);
preg_match_all($urlsRegex, $fileContent, $matches);
debug($matches);
foreach ($matches[1] as $url) {
    if (class_exists(AsyncDownload::class)) {
        $thread = new AsyncDownload($url);
        $thread->start();
        // $thread->join();
    } else {
        download($url);
    }
}

function download($url)
{
    global $suffixesStr;
    $fileRegex = '/^https?:\/\/.*\/(.*?\.(' . $suffixesStr . '))/';
    preg_match($fileRegex, $url, $matches);
    debug($matches);
    $fileName = $matches[1];
    echo "$url $fileName\n";
    if (file_exists($fileName)) {
        echo "exist\n";
    } else {
        echo "not exist, downloading...\n";
        $ret = true;
        if (!DEBUG) {
            $ret = file_put_contents($fileName, file_get_contents($url));
        }
        if ($ret === false) {
            dd(__LINE__);
        } else {
            echo "downloaded\n";
        }
    }

}
