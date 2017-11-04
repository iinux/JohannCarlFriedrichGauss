<?php
$filename="gzread.php.gz";
$zd=gzopen($filename,"r");
$contents=gzread($zd,10000);
gzclose($zd);
echo $contents;
?>
