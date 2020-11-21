<?php
/**
 * Created by PhpStorm.
 * User: iinux
 * Date: 2020/1/25
 * Time: 2:21 PM
 */

/**
 * @param $var
 */
function dd($var) {
    var_dump(($var));
    die(0);
}

$mac = '00:24:9b:3e:94:e3';
$mac = $argv[1];

if (strpos($mac, '-')) {
    $macs = explode('-', $mac);
} else {
    $macs = explode(':', $mac);
}
$macs = array_merge(array_slice($macs, 0, 3), ['ff', 'fe'], array_slice($macs, 3));
$macs[0] = dechex(hexdec($macs[0]) ^ 0b00000010);
$ipv6 = 'fe80:';
$odd = true;
foreach ($macs as $item) {
    if ($odd) {
        $ipv6 .= ":$item";
    } else {
        $ipv6 .= "$item";
    }
    $odd = !$odd;
}
$ipv6 = strtoupper($ipv6);
dd($ipv6);
//refer
//https://www.vultr.com/resources/mac-converter/
//https://baike.baidu.com/item/EUI-64
