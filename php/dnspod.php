<?php

require_once 'passwords.php';

function postData($url, $data)
{
    $data = http_build_query($data);
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, "$url");
// curl_setopt($ch, CURLOPT_HEADER, 1);
// curl_setopt($ch, CURLOPT_VERBOSE, 1);
    curl_setopt($ch, CURLOPT_POST, 1);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
//    curl_setopt($ch, CURLOPT_USERAGENT, '');
    curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
    $response = curl_exec($ch);
    curl_close($ch);
    return $response;
}

$data = [
    'login_token' => $dnspodToken,
    'format' => 'json',
];

$action = 'list';
$res = '{}';

if (!empty($argv[1])) {
    $action = $argv[1];
}

$domain_id = '92247536';

switch ($action) {
    case 'list_domain':
        $res = postData("https://dnsapi.cn/Domain.List", $data);
        print_r(json_decode($res));
        break;
    case 'list':
    case 'list_record':
        $data['domain_id'] = $domain_id;
        $res = postData("https://dnsapi.cn/Record.List", $data);
        $cols = [
            'id' => "%s\t",
            //'ttl' => "%s\t",
            'value' => "%-20s\t",
            //'enabled' => "%s\t",
            //'status' => "%s\t",
            //'updated_on' => "%s\t",
            'name' => "%s\t\t",
            //'line' => "%s\t",
            'line_id' => "%s\t\t",
            'type' => "%s\t\t",
            //'weight' => "%s\t",
            //'monitor_status' => "%s\t",
            //'remark' => "%s\t",
            //'use_aqb' => "%s\t",
            'mx' => "%s\t",
            //'hold' => "%s\t",
        ];
        foreach ($cols as $k => $v) {
            printf($v, $k);
        }
        echo "\n";
        foreach (json_decode($res, true)['records'] as $record) {
            foreach ($cols as $k => $v) {
                printf($v, $record[$k]);
            }
            echo "\n";
        }
        break;
    case 'add':
        $data['domain_id'] = $domain_id;
        $data['record_line_id'] = 0;
        $data['sub_domain'] = $argv[2];
        $data['record_type'] = $argv[3];
        $data['value'] = $argv[4];
        $res = postData("https://dnsapi.cn/Record.Create", $data);
        print_r(json_decode($res));
        break;
    case 'edit':
        $data['domain_id'] = $domain_id;
        $data['record_line_id'] = 0;
        $data['record_id'] = $argv[2];
        $data['sub_domain'] = $argv[3];
        $data['record_type'] = $argv[4];
        $data['value'] = $argv[5];
        $res = postData("https://dnsapi.cn/Record.Modify", $data);
        print_r(json_decode($res));
        break;
    case 'del':
        $data['domain_id'] = $domain_id;
        $data['record_id'] = $argv[2];
        $res = postData("https://dnsapi.cn/Record.Remove", $data);
        print_r(json_decode($res));
        break;
}

