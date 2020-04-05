<?php

class sftpData
{
    // 初始配置为NULL
    private $config = NULL;
    // 连接为NULL
    private $conn = NULL;
    //sftp resource
    private $resSftp = NULL;

    // 初始化
    public function __construct($config)
    {
        $this->config = $config;
        $this->connect();
    }

    public function connect()
    {

        $this->conn = ssh2_connect($this->config['host'], $this->config['port']);
        if (ssh2_auth_password($this->conn, $this->config['username'], $this->config['password'])) {
            $this->resSftp = ssh2_sftp($this->conn);

        } else {
            echo "用户名或密码错误";
        }

    }

    // 下载文件
    public function downftp($remote, $local)
    {
        return copy("ssh2.sftp://{$this->resSftp}" . $remote, $local);
    }

    // 文件上传
    public function upftp($local, $remote, $file_mode = 0777)
    {
        return copy($local, "ssh2.sftp://{$this->resSftp}" . $remote);
        // return ssh2_scp_send($this->ressftp,$local,$remote,$file_mode);
    }

    //创建目录
    public function ssh2_sftp_mchkdir($path)  //使用创建目录循环
    {
        ssh2_sftp_mkdir($this->resSftp, $path, 0777, true);
    }

    //判段目录是否存在
    public function ssh2_dir_exits($dir)
    {
        return file_exists("ssh2.sftp://{$this->resSftp}" . $dir);
    }
}

function getSftpData()
{
    $time = date('Ymd', time());
    $config = [
        "host"     => "baidu.com",
        "username" => "hello",
        "port"     => "8000",
        "password" => "xxxxxx"
    ];
    $localPath = 'a.go';   //本地文件目录
    $serverPath = "/{$config['username']}/$time";  //远程目录（需要上传到的目录）
    try {
        $sftp = new \sftpData($config);

        $re = $sftp->ssh2_dir_exits($serverPath);
        //如果目录存在直接上传
        if (!$re) {
            $sftp->ssh2_sftp_mchkdir($serverPath);
        }
        // $sftp->upftp("$localPath", "$serverPath/$localPath");
        $sftp->downftp("$serverPath/$localPath", "$localPath");
    } catch (\Exception $e) {
        die('连接失败！');
    }
}

getSftpData();
