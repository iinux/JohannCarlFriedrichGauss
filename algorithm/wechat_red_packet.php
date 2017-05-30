<?php
// $bonus_total总数 $bonus_count个数 $bonus_type类型
function randBonus($bonus_total = 0, $bonus_count = 3, $bonus_type = 1)
{
    $bonus_items = array(); // 结果
    $bonus_balance = $bonus_total; // 余额
    $bonus_avg = number_format($bonus_total / $bonus_count, 2); // 平均数
    $i = 0;
    while ($i < $bonus_count) {
        if ($i < $bonus_count - 1) {
            $rand = $bonus_type ? (rand(1, $bonus_balance * 100 - 1) / 100) : $bonus_avg;
            $bonus_items[] = $rand;
            $bonus_balance -= $rand;
        } else {
            $bonus_items[] = $bonus_balance;
        }
        $i++;
    }
    return $bonus_items;
}

$bonus_items = randBonus(100, 3, 1);
var_dump($bonus_items);
var_dump(array_sum($bonus_items));

function sendRandBonus($total = 0, $count = 3, $type = 1)
{
    // $total总数 $count个数 $type类型
    if ($type == 1) {
        $input = range(0.01, $total, 0.01);
        if ($count > 1) {
            $rand_keys = (array)array_rand($input, $count - 1);
            $last = 0;
            foreach ($rand_keys as $i => $key) {
                $current = $input[$key] - $last;
                $items[] = $current;
                $last = $input[$key];
            }
        }
        $items[] = $total - array_sum($items);
    } else {
        $avg = number_format($total / $count, 2);
        $i = 0;
        while ($i < $count) {
            $items[] = $i < $count - 1 ? $avg : ($total - array_sum($items));
            $i++;
        }
    }
    return $items;
}

$items = sendRandBonus(100, 3, 1);
var_dump($items);