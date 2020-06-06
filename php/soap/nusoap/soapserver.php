<?php

// 在PHP7.0下已无法调通

require_once('./lib/nusoap.php');
$soap = new soap_server;

$soap->soap_defencoding = 'UTF-8';
$soap->decode_utf8 = false;
$soap->xml_encoding = 'UTF-8';
$soap->configureWSDL('reverse', 'add2numbers');
//$server->configureWSDL('add2numbers');

$soap->register('reverse',
    array('str' => 'xsd:string'),
    array('return' => 'xsd:string')
);
$soap->register('add2numbers',
    array('num1' => 'xsd:int', 'num2' => 'xsd:int'),
    array('return' => 'xsd:int')
);

//Check variable set?
$HTTP_RAW_POST_DATA = isset($HTTP_RAW_POST_DATA) ? $HTTP_RAW_POST_DATA : '';

//service handle the client data
$soap->service($HTTP_RAW_POST_DATA);
//echo $server;
/*
* RPC function
* @param $name
*/
function reverse($str)
{
    $retval = "";
    if (strlen($str) < 1) {
        return new soap_fault('Client', '', 'Invalid string');
    }
    for ($i = 1; $i <= strlen($str); $i++) {
        $retval .= $str[(strlen($str) - $i)];
    }
    return $retval;
}

function add2numbers($num1, $num2)
{
    if (trim($num1) != intval($num1)) {
        return new soap_fault('Client', '', 'The first number is invalid');
    }
    if (trim($num2) != intval($num2)) {
        return new soap_fault('Client', '', 'The second number is invalid');
    }
    return ($num1 + $num2);
}

?>
