<?php

// 服务器验证  
if ($_SERVER['PHP_AUTH_USER'] != 'fdipzone' || $_SERVER['PHP_AUTH_PW'] != '123456') {
    header('WWW-Authenticate: Basic realm="NMG Terry"');
    header('HTTP/1.0 401 Unauthorized');
    echo "You must enter a valid login ID and password to access this resource.\n";
    exit();
}

class SOAPHandle
{
    // header 驗證
    public function auth($auth)
    {
        if ($auth->string[0] != 'fdipzone' || $auth->string[1] != '654321') {
            throw new SOAPFault('Server', 'No Permission');
        }
    }

    // 反轉字符串
    public function revstring($str = '')
    {
        return strrev($str);
    }

    // 字符傳轉連接
    public function strtolink($str = '', $name = '', $openwin = 0)
    {
        $name = $name == '' ? $str : $name;
        $openwin_tag = $openwin == 1 ? ' target="_blank" ' : '';
        return sprintf('<a href="%s" %s>%s</a>', $str, $openwin_tag, $name);
    }

    // 字符串轉大寫
    public function uppcase($str)
    {
        return strtoupper($str);
    }
}

$config = array(
    'uri' => 'http://192.168.188.218:9200/header_server.php',
    // 'actor' => 'myserver',
);

$oHandle = new SOAPHandle;

// no wsdl mode
try {

    $server = new SOAPServer(null, $config);
    $server->setObject($oHandle);
    $server->handle();

} catch (SOAPFault $f) {
    echo $f->getMessage();
}
