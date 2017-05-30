<?php
//摘自：http://www.dewen.org/q/10775
function a() {
}

class b {
    public function f() {
    }
}

function function_dump($funcname) {
    try {
        if(is_array($funcname)) {
            $func = new ReflectionMethod($funcname[0], $funcname[1]);
            $funcname = $funcname[1];
        } else {
            $func = new ReflectionFunction($funcname);
        }
    } catch (ReflectionException $e) {
        echo $e->getMessage();
        return;
    }
    $start = $func->getStartLine() - 1;
    $end =  $func->getEndLine() - 1;
    $filename = $func->getFileName();
    echo "function $funcname defined by $filename($start - $end)\n";
}
function_dump('a');
function_dump(array('b', 'f'));
$b = new b();
function_dump(array($b, 'f'));
?>
