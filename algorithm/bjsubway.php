<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2016/5/19
 * Time: 9:18
 */
define('DISCOUNT_BEFORE_SEVEN_CLOCK', 0.5);
define('DISCOUNT_ALIPAY', 0.8);
define('DISCOUNT_ALIPAY_FEE', 19.99);
define('DISCOUNT_ALIPAY_FEE_DISCOUNT', 0.2);
define('DISCOUNT_100', 0.8);
define('DISCOUNT_150', 0.5);
$website = 'Subway Calculation';

$basePrice = isset($_POST['base_price']) ? $_POST['base_price'] : '';
$maxPrice = 9;

function discount($basePrice, $sum, $beforeSevenClock = false, $alipayDiscount = false)
{
    if ($beforeSevenClock) {
        $basePrice *= DISCOUNT_BEFORE_SEVEN_CLOCK;
    }
    if ($alipayDiscount) {
        $basePrice *= DISCOUNT_ALIPAY;
    }
    if ($sum < 100) {

    } else if ($sum < 150) {
        $basePrice *= DISCOUNT_100;
    } else if ($sum < 400) {
        $basePrice *= DISCOUNT_150;
    }
    return $basePrice;
}

?>
    <html>
    <head>
        <meta charset="utf-8">
        <title><?php echo $website ?></title>
        <style type="text/css">

            ::selection {
                background-color: #E13300;
                color: white;
            }

            ::-moz-selection {
                background-color: #E13300;
                color: white;
            }

            body {
                background-color: #fff;
                margin: 40px;
                font: 13px/20px normal Helvetica, Arial, sans-serif;
                color: #4F5155;
            }

            a {
                color: #003399;
                background-color: transparent;
                font-weight: normal;
            }

            h1 {
                color: #444;
                background-color: transparent;
                border-bottom: 1px solid #D0D0D0;
                font-size: 19px;
                font-weight: normal;
                margin: 0 0 14px 0;
                padding: 14px 15px 10px 15px;
            }

            code {
                font-family: Consolas, Monaco, Courier New, Courier, monospace;
                font-size: 12px;
                background-color: #f9f9f9;
                border: 1px solid #D0D0D0;
                color: #002166;
                display: block;
                margin: 14px 0 14px 0;
                padding: 12px 10px 12px 10px;
            }

            #body {
                margin: 0 15px 0 15px;
            }

            p.footer {
                text-align: right;
                font-size: 11px;
                border-top: 1px solid #D0D0D0;
                line-height: 32px;
                padding: 0 10px 0 10px;
                margin: 20px 0 0 0;
            }

            #container {
                margin: 10px;
                border: 1px solid #D0D0D0;
                box-shadow: 0 0 8px #D0D0D0;
            }
        </style>
        <?php
        if ($_SERVER['REQUEST_METHOD'] == "POST") {
            ?>
            <!-- <script src="https://as.alipayobjects.com/g/datavis/g2/1.2.1/index.js"></script> -->
            <script src=g2-1.2.1.js></script>
            <?php
        }
        ?>

    </head>
<body>
<div id="container">
    <h1><?php echo $website ?></h1>

    <div id="body">
<?php
if ($_SERVER['REQUEST_METHOD'] == "GET") {
    ?>
    <form action="" method="post">
        BASE PRICE:&nbsp;<input name="base_price" type="text"/><br/>
        <!-- Before Seven Clock:&nbsp;<input name="before_seven_clock" type="checkbox" /><br /> -->
        <input type="submit" value="submit"/>
    </form>
<?php
return;
} else {
?>
    <!-- Step 2: 创建图表容器 -->

    <div id="c1"></div>
    <script>
        //Step 3: 创建图表

        var chart = new G2.Chart({
            id: 'c1',
            width: 1000,
            height: 500
        });
        //Step 4: 指定数据源

        var data = [
            <?php
            if ($basePrice != '') {
                $sum = 0;
                $day = 1;
                while ($day <= 31) {
                    $sum += discount($basePrice, $sum, true, false);
                    $sum += discount($basePrice, $sum, true, false);
                    echo "{'day':$day,'price':$sum,'discount':'beforeSeven'},\n";
                    $day++;
                }

                $sum = DISCOUNT_ALIPAY_FEE * DISCOUNT_ALIPAY_FEE_DISCOUNT;
                $day = 1;
                while ($day <= 31) {
                    $sum += discount($basePrice, $sum, false, true);
                    $sum += discount($basePrice, $sum, false, true);
                    echo "{'day':$day,'price':$sum,'discount':'alipay'},\n";
                    $day++;
                }

                $sum = 0;
                $day = 1;
                while ($day <= 31) {
                    $sum += discount($basePrice, $sum, false, false);
                    $sum += discount($basePrice, $sum, false, false);
                    echo "{'day':$day,'price':$sum,'discount':'none'},\n";
                    $day++;
                }
            } else {
                $price = 3;
                while ($price <= $maxPrice) {
                    $sum = 0;
                    $day = 1;
                    while ($day <= 31) {
                        $sum += discount($price, $sum);
                        $sum += discount($price, $sum);
                        echo "{'day':$day,'price':$sum,'Base price':'$price'},\n";
                        $day++;
                    }
                    $price++;
                }
            }
            ?>
        ];
        chart.source(data);
        //Step 5: 指定绘制的图形
        <?php
        if ($basePrice != ''){
        ?>
        chart.line().position('day*price').color('discount');
        <?php
        }else{
        ?>
        chart.line().position('day*price').color('Base price');
        <?php
        }
        ?>
        //Step 6: 绘制图表

        chart.render();
    </script>

    </div>

    <p class="footer">@2019 Perorsoft</p>
    </div>

    </body>
    </html>
    <?php
}
?>
