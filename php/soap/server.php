<?php

class soapHandle
{
    public function strtolink($url = '')
    {
        return sprintf('<a href="%s">%s</a>', $url, $url);
    }
}

if ($_SERVER['PHP_AUTH_USER'] != 'fdipzone' || $_SERVER['PHP_AUTH_PW'] != '123456') {
    header('WWW-Authenticate: Basic realm="MyFramework Realm"');
    header('HTTP/1.0 401 Unauthorized');
    echo "You must enter a valid login ID and password to access this resource.\n";
    exit;
}

try {
    $server = new SOAPServer(null, array('uri' => 'http://192.168.188.218:9200/server.php'));
    $server->setClass('soapHandle'); //设置处理的class
    $server->handle();
} catch (SOAPFault $f) {
    echo $f->getMessage();
}
