<?php
/**
 * Created by PhpStorm.
 * User: nalux
 * Date: 2017/5/17
 * Time: 22:52
 */

$num1 = $argv[1];
$num2 = $argv[2];
$num3 = $argv[3];
$num4 = $argv[4];
foreach (['+','-','*','/'] as $op1) {
    foreach (['+','-','*','/'] as $op2) {
        foreach (['+','-','*','/'] as $op3) {
            $exp = "\$a=$num1$op1$num2$op2$num3$op3$num4;";
            eval($exp);
            if ($a == 517) {
                var_dump($exp);
            }
        }
    }
}
