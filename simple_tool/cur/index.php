<?php

function dd($var) {
    var_dump($var);
    die(0);
}

function jdd($var) {
    echo json_encode($var);
    die(0);
}

function checkCookieSwitch($name) {
    if (!empty($_COOKIE[$name])) {
        return true;
    } else {
        return false;
    }
}

// config
$protocol = 'https://';
$domain = 'www.baidu.com';
$authUsername = 'admin';
$authPassword = 'admin';
$resolve = '';
$headerValueReplace = [
    'baidu.com' => 'iinux.cn',
];
if (file_exists('config.php')) {
    require 'config.php';
}

// auth
if (!empty($_COOKIE['username']) && !empty($_COOKIE['password'])
    && $_COOKIE['username'] == $authUsername && $_COOKIE['password'] == $authPassword
) {
} elseif (!empty($_SERVER['PHP_AUTH_USER']) && !empty($_SERVER['PHP_AUTH_PW'])
    && $_SERVER['PHP_AUTH_USER'] == $authUsername && $_SERVER['PHP_AUTH_PW'] == $authPassword
) {
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
if (!empty($_COOKIE['resolve'])) {
    $resolve = $_COOKIE['resolve'];
}

if (checkCookieSwitch('phpinfo')) {
    phpinfo();
    die(0);
}
if (checkCookieSwitch('displayError')) {
    ini_set('display_errors', 1);
}

// var
$requestUri = $_SERVER['REQUEST_URI'];
$method = $_SERVER['REQUEST_METHOD'];
$uri = $protocol . $domain . $requestUri;

// header process
$headers = getallheaders();
$oldHost = $headers['Host'];
$headers['Host'] = $domain;
if (!empty($headers['Referer'])) {
    $headers['Referer'] = str_replace($oldHost, $domain, $headers['Referer']);
}

if (checkCookieSwitch('dumpUH')) {
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

if (!empty($resolve)) {
    curl_setopt($ch, CURLOPT_DNS_USE_GLOBAL_CACHE, false);
    curl_setopt($ch, CURLOPT_RESOLVE, [
        $resolve,
    ]);
}

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

if (checkCookieSwitch('skipSslCheck')) {
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, 0);
}
if (checkCookieSwitch('debug')) {
    curl_setopt($ch, CURLOPT_VERBOSE, true);
    $verbose = fopen('php://temp', 'w+');
    curl_setopt($ch, CURLOPT_STDERR, $verbose);
}

$content = curl_exec($ch);

if ($content === false) {
    echo 'Curl error ' . curl_errno($ch) . ' : ' . curl_error($ch);
    if (checkCookieSwitch('debug')) {
        rewind($verbose);
        $verboseLog = stream_get_contents($verbose);
        echo " Verbose information:\n<pre>", htmlspecialchars($verboseLog), "</pre>\n";
    }
    die(0);
}

$headerSize = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
$contentHeader = substr($content, 0, $headerSize);
$contentBody = substr($content, $headerSize);

if (checkCookieSwitch('dumpRH')) {
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
        foreach ($headerValueReplace as $k => $v) {
            $headerValue = str_replace($k, $v, $headerValue);
        }
        header("$headerName: $headerValue", false);
    }
}

echo $contentBody;

curl_close($ch);

