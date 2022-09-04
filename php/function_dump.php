<?php
//摘自：http://www.dewen.org/q/10775
function a() {
}

class b {
    public function f() {
    }
}

function function_dump($funcName) {
    try {
        if(is_array($funcName)) {
            $func = new ReflectionMethod($funcName[0], $funcName[1]);
            $funcName = $funcName[1];
        } else {
            $func = new ReflectionFunction($funcName);
        }
    } catch (ReflectionException $e) {
        echo $e->getMessage();
        return;
    }
    $start = $func->getStartLine() - 1;
    $end =  $func->getEndLine() - 1;
    $filename = $func->getFileName();
    echo "function $funcName defined by $filename($start - $end)\n";
}
function_dump('a');
function_dump(array('b', 'f'));
$b = new b();
function_dump(array($b, 'f'));

