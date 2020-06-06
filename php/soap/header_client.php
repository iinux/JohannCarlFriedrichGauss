<?php

$config = array(
    'location' => 'http://192.168.188.218:9200/header_server.php',
    'uri'      => 'http://192.168.188.218:9200/header_server.php',
    'login'    => 'fdipzone',
    'password' => '123456',
    'trace'    => true,
);

try {
    $auth = array('fdipzone', '654321');

    // no wsdl
    $client = new SOAPClient(null, $config);
    // $header = new SOAPHeader('http://demo.fdipzone.com/soap/sever.php', 'auth', $auth, false, SOAP_ACTOR_NEXT);
    $header = new SOAPHeader('http://demo.fdipzone.com/soap/sever.php', 'auth', $auth, false, 'actor');
    $client->__setSoapHeaders(array($header));

    $revstring = $client->revstring('123456');
    $strtolink = $client->__soapCall('strtolink', array('http://blog.csdn.net/fdipzone', 'fdipzone blog', 1));
    $uppcase = $client->__soapCall('uppcase', array('Hello World'));

    echo $revstring . '<br>';
    echo $strtolink . '<br>';
    echo $uppcase . '<br>';
} catch (SOAPFault $e) {
    echo $e->getMessage();
}
