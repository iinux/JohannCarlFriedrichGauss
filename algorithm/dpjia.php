<?php
function f1()
{
    echo (strtotime('2016-2-23') - strtotime('2015-1-8')) / (60 * 60 * 24) . " days";
}

function my_array_filter($arr, $want, $misc)
{
    $result = array();
    foreach ($want as $name) {
        if (array_key_exists($name, $arr))
            $result[$name] = $arr[$name];
        else
            $result[$name] = NULL;
    }

}

function f2()
{
    $array = array(
        'color'    => 'red',
        'shape'    => 'round',
        'radius'   => '10',
        'diameter' => '20'
    );
    $wanted_keys = array('color', 'shape', 'height');
    $result = my_array_filter($array, $wanted_keys, NULL);
}

$arr = array(89, 27, 109, 78, 32, 12, 0, 88, 63, 2, 5, 56);

function f3_1()
{
    for ($i = 0; $i < count($arr); $i++) {
        for ($j = $i + 1; $j < count($arr); $j++) {
            if ($arr[$i] > $arr[$j]) {
                $t = $arr[$i];
                $arr[$i] = $arr[$j];
                $arr[$j] = $t;
            }
        }
    }
}

function f3_2()
{
    for ($i = 1; $i < count($arr); $i++) {
        for ($j = $i - 1; $j >= 0; $j--) {
            if ($arr[$i] < $arr[$j]) {
                $t = $arr[$i];
                $arr[$i] = $arr[$j];
                $arr[$j] = $t;
                $i--;
            } else {
                break;
            }
        }
    }
}

function my_shuffle()
{
    $arr = range(1, 10);
    array_push($arr, 'J');
    array_push($arr, 'Q');
    array_push($arr, 'K');

    $result = range(-54, -1, 1);

    for ($i = 0; $i < 4; $i++) {
        foreach ($arr as $el) {
            while (true) {
                $index = rand(0, 53);
                if ($result[$index] < 0) {
                    $result[$index] = $el;
                    break;
                }
            }
        }
    }
    while (true) {
        $index = rand(0, 53);
        if ($result[$index] < 0) {
            $result[$index] = 'SG';
            break;
        }
    }
    while (true) {
        $index = rand(0, 53);
        if ($result[$index] < 0) {
            $result[$index] = 'BG';
            break;
        }
    }
    var_dump($result);


}

my_shuffle();

function f5()
{
    $len = count($arr);
    for ($i = 0; $i < $len - 1; $i++) {
        $index = rand($i + 1, $len - 1);
        $t = $arr[$i];
        $arr[$i] = $arr[$index];
        $arr[$index] = $t;
    }

    $i = 0;
    while ($i < $len - 1) {
        $index = rand($i + 1, $len - 1);
        $t = $arr[$i];
        $arr[$i] = $arr[$index];
        $arr[$index] = $t;
        $i++;
    }
}

var_dump(array_fill(0, 3, 0));

?> 