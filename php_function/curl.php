<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2018/4/25
 * Time: 14:12
 */

function curl_call($p1, $p2, $times = 1)
{
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_TIMEOUT, 5);
    curl_setopt($ch, CURLOPT_URL, 'http://demon.at');
    $curl_version = curl_version();
    if ($curl_version['version_number'] >= 462850) {
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT_MS, 20);
        curl_setopt($ch, CURLOPT_NOSIGNAL, 1);
    } else {
        throw new Exception('this curl version is too low, version_num : ' . $curl_version['version']);
    }
    $res = curl_exec($ch);
    curl_close($ch);
    if (false === $res) {
        if (curl_errno($ch) == CURLE_OPERATION_TIMEOUTED && $times != 5) {
            $times += 1;
            return curl_call($p1, $p2, $times);
        }
    }

    return $res;
}

var_dump(curl_version());