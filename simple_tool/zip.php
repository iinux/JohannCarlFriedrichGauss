<?php
//获取文件列表
function list_dir($dir)
{
    $result = array();
    if (is_dir($dir)) {
        $file_dir = scandir($dir);
        foreach ($file_dir as $file) {
            if ($file == '.' || $file == '..') {
                continue;
            } elseif (is_dir($dir . $file)) {
                $result = array_merge($result, list_dir($dir . $file . '/'));
            } else {
                array_push($result, $dir . $file);
            }
        }
    }
    return $result;
}

//获取列表
$datalist = list_dir('./');
$filename = "./bak.zip"; //最终生成的文件名（含路径）
if (!file_exists($filename)) {
    //重新生成文件
    $zip = new ZipArchive();//使用本类，linux需开启zlib，windows需取消php_zip.dll前的注释
    if ($zip->open($filename, ZipArchive::CREATE) !== TRUE) {
        exit('无法打开文件，或者文件创建失败');
    }
    foreach ($datalist as $val) {
        if (file_exists($val)) {
            $zip->addFile($val, basename($val));//第二个参数是放在压缩包中的文件名称，如果文件可能会有重复，就需要注意一下
        }
    }
    $zip->close();//关闭
}
if (!file_exists($filename)) {
    exit("无法找到文件"); //即使创建，仍有可能失败。。。。
}
header("Cache-Control: public");
header("Content-Description: File Transfer");
header('Content-disposition: attachment; filename=' . basename($filename)); //文件名
header("Content-Type: application/zip"); //zip格式的
header("Content-Transfer-Encoding: binary"); //告诉浏览器，这是二进制文件
header('Content-Length: ' . filesize($filename)); //告诉浏览器，文件大小
@readfile($filename);


//一、解压缩zip文件
$zip = new ZipArchive;//新建一个ZipArchive的对象
if ($zip->open('test.zip') === TRUE) {
    $zip->extractTo('images');//假设解压缩到在当前路径下images文件夹内
    $zip->close();//关闭处理的zip文件
}

//二、将文件压缩成zip文件
$zip = new ZipArchive;
if ($zip->open('test.zip', ZipArchive::OVERWRITE) === TRUE) {
    $zip->addFile('image.txt');//假设加入的文件名是image.txt，在当前路径下
    $zip->close();
}

//三、文件追加内容添加到zip文件
$zip = new ZipArchive;
$res = $zip->open('test.zip', ZipArchive::CREATE);
if ($res === TRUE) {
    $zip->addFromString('test.txt', 'file content goes here');
    $zip->close();
    echo 'ok';
} else {
    echo 'failed';
}

//四、将文件夹打包成zip文件
function addFileToZip($path, $zip)
{
    $handler = opendir($path); //打开当前文件夹由$path指定。
    while (($filename = readdir($handler)) !== false) {
        if ($filename != "." && $filename != "..") {//文件夹文件名字为'.'和‘..’，不要对他们进行操作
            if (is_dir($path . "/" . $filename)) {// 如果读取的某个对象是文件夹，则递归
                addFileToZip($path . "/" . $filename, $zip);
            } else { //将文件加入zip对象
                $zip->addFile($path . "/" . $filename);
            }
        }
    }
    @closedir($path);
}

$zip = new ZipArchive();
if ($zip->open('images.zip', ZipArchive::OVERWRITE) === TRUE) {
    addFileToZip('images/', $zip); //调用方法，对要打包的根目录进行操作，并将ZipArchive的对象传递给方法
    $zip->close(); //关闭处理的zip文件
}