<?php
function greet($param)
{
	$value = 'Hello '.$param;
	return new SoapParam($value, 'greetReturn');
}

$server = new SoapServer(null, [
	'uri' => 'http://localhost/php-soap/non-wsdl/helloService'
]);
$server->addFunction('greet');
$server->handle();