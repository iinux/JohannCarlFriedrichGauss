<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2018/12/5
 * Time: 10:22
 */

// 在PHP7.0下已无法调通

require_once "./lib/nusoap.php";

$client = new nusoap_client("http://192.168.188.218:9200/soapserver.php");
$client->soap_defencoding = 'UTF-8';
$client->decode_utf8 = false;
$client->xml_encoding = 'UTF-8';

$str = "This string will be reversed";
$params1 = array('str' => $str);
$reversed = $client->call('reverse', $params1);
echo "If you reverse '$str', you get '$reversed'<br>\n";

$n1 = 5;
$n2 = 14;
$params2 = array('num1' => $n1, 'num2' => $n2);
$added = $client->call('add2numbers', $params2);
echo "If you add $n1 and $n2 you get $added<br>\n";
