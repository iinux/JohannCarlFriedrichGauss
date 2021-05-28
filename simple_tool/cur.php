<?php

function dd($var) {
    var_dump($var);
    die(0);
}

function jdd($var) {
    echo json_encode($var);
    die(0);
}

// config
$protocol = 'https://';
$domain = 'duckduckgo.com';
$authUsername = 'admin';
$authPassword = 'admin';
$host = $domain;

// auth
if (!empty($_COOKIE['username']) && !empty($_COOKIE['password'])
    && $_COOKIE['username'] == $authUsername && $_COOKIE['password'] == $authPassword
) {
} elseif (!empty($_SERVER['PHP_AUTH_USER']) && !empty($_SERVER['PHP_AUTH_PW'])
    && $_SERVER['PHP_AUTH_USER'] == $authUsername && $_SERVER['PHP_AUTH_PW'] == $authPassword) {
} else {
    header('WWW-Authenticate: Basic realm="My Realm"');
    header('HTTP/1.0 401 Unauthorized');
    echo 'Text to send if user hits Cancel button';
    exit;
}

if (!empty($_COOKIE['domain'])) {
    $domain = $_COOKIE['domain'];
}
if (!empty($_COOKIE['protocol'])) {
    $protocol = $_COOKIE['protocol'];
}
if (!empty($_COOKIE['host'])) {
    $host = $_COOKIE['host'];
}
if (!empty($_COOKIE['phpinfo'])) {
    phpinfo();
    die(0);
}
if (!empty($_COOKIE['displayError'])) {
    ini_set('display_errors', 1);
}

// var
$requestUri = $_SERVER['REQUEST_URI'];
$method = $_SERVER['REQUEST_METHOD'];
$uri = $protocol . $domain . $requestUri;

// header process
$headers = getallheaders();
$oldHost = $headers['Host'];
$headers['Host'] = $host;
if (!empty($headers['Referer'])) {
    $headers['Referer'] = str_replace($oldHost, $host, $headers['Referer']);
}

if (!empty($_COOKIE['dumpUH'])) {
    jdd([
        'uri'     => $uri,
        'headers' => $headers,
    ]);
}

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $uri);
curl_setopt($ch, CURLOPT_HEADER, 1);
curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

$responseHeaders = [];
curl_setopt($ch, CURLOPT_HEADERFUNCTION,
    function ($curl, $header) use (&$responseHeaders) {
        $len = strlen($header);
        $header = explode(':', $header, 2);
        if (count($header) < 2) // ignore invalid headers
            return $len;

        $responseHeaders[strtolower(trim($header[0]))][] = trim($header[1]);

        return $len;
    }
);

if (!empty($_COOKIE['skipSslCheck'])) {
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, 0);
}

$content = curl_exec($ch);

if ($content === false) {
    echo 'Curl error: ' . curl_error($ch);
    die(0);
}

$headerSize = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
$contentHeader = substr($content, 0, $headerSize);
$contentBody = substr($content, $headerSize);

if (!empty($_COOKIE['dumpRH'])) {
    jdd($responseHeaders);
}

foreach ($responseHeaders as $headerName => $headerValues) {
    $lower = strtolower($headerName);
    if (in_array($lower, [
        'server',
        'date',
        'transfer-encoding',
        'alt-svc',
        'content-security-policy',
    ])) {
        continue;
    }
    foreach ($headerValues as $headerValue) {
        $map = [
            'baidu.com' => 'iinux.cn',
        ];
        foreach ($map as $k => $v) {
            $headerValue = str_replace($k, $v, $headerValue);
        }
        header("$headerName: $headerValue", false);
    }
}

echo $contentBody;

curl_close($ch);

