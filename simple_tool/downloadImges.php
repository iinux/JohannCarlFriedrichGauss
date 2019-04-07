<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2019/4/4
 * Time: 10:29
 */

define('DEBUG', false);
$htmlFile = 'action.txt';
$urlsRegex = '/data-src=\'(https?:\/\/.*?\.(gif|png|jpg|jpeg))\'/';

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

$fileContent = file_get_contents($htmlFile);
preg_match_all($urlsRegex, $fileContent, $matches);
debug($matches);
foreach ($matches[1] as $url) {
    download($url);
}

function download($url)
{
    preg_match('/^https?:\/\/.*\/(.*?\.(gif|png|jpg|jpeg))/', $url, $matches);
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
