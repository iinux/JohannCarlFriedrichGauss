<?php
function greet($param)
{
    $value = 'Hello '.$param->name;
    $result = array('greetReturn' => $value);
    return $result;
}

$server = new SoapServer('hello.wsdl');
$server->addFunction('greet');
$server->handle();