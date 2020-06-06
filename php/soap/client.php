<?php

try {
    $client = new SOAPClient(null, array(
        'location' => 'http://192.168.188.218:9200/server.php', // 设置server路径
        'uri'      => 'http://192.168.188.218:9200/server.php',
        'login'    => 'fdipzone', // HTTP auth login
        'password' => '123456' // HTTP auth password
    ));

    echo $client->strtolink('http://blog.csdn.net/fdipzone') . '<br>';               // 直接调用server方法
    echo $client->__soapCall('strtolink', array('http://blog.csdn.net/fdipzone')); // 间接调用server方法
} catch (SOAPFault $e) {
    throw $e;
    echo $e->getMessage();
}
