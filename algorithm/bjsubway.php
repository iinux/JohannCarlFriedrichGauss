<?php
/**
 * Created by PhpStorm.
 * User: qzhang
 * Date: 2016/5/19
 * Time: 9:18
 */
$website = 'Subway Calculation';

$basePrice = isset($_POST['base_price']) ? $_POST['base_price'] : '';
$maxPrice = 9;

function dd($var)
{
    var_dump($var);
    die(0);
}

class Discount {
    public $sum = 0;
    public $discount_100 = 0.8;
    public $discount_150 = 0.5;

    public function add($price) {
        $this->sum += $this->getDiscountPrice($price);
    }

    public function getDiscountPrice($price) {
        if ($this->sum < 100) {

        } else if ($this->sum < 150) {
            $price *= $this->discount_100;
        } else if ($this->sum < 400) {
            $price *= $this->discount_150;
        }
        return $price;
    }
}

class AliDiscount extends Discount {
    public function __construct()
    {
        $this->sum += $this->cardFee * $this->cardFeeDiscount;
    }

    public $discountRate = 0.8;
    public $cardFee = 19.99;
    public $cardFeeDiscount = 0.2;

    public function getDiscountPrice($price)
    {
        $price *= $this->discountRate;

        return parent::getDiscountPrice($price);
    }
}

class JdDiscount extends Discount {
    public function __construct()
    {
        $this->sum += $this->cardFee * $this->cardFeeDiscount;
    }

    public $discountFee = 1.5;
    public $discountSum = 0;
    public $maxDiscount = 90;
    public $cardFee = 20;
    public $cardFeeDiscount = 0.7245;

    public function add($price) {
        $this->sum += $this->getDiscountPrice($price, true);
        var_dump($this->sum);
    }

    public function getDiscountPrice($price, $updateDiscountSum = false)
    {
        if ($this->discountSum < $this->maxDiscount) {
            $price -= $this->discountFee;
            if ($updateDiscountSum) {
                $this->discountSum += $this->discountFee;
            }
        }

        return parent::getDiscountPrice($price);
    }
}

class SevenClockDiscount extends Discount {

    public $discountRate = 0.5;

    public function getDiscountPrice($price)
    {
        $price *= $this->discountRate;

        return parent::getDiscountPrice($price);
    }
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
                $day = 1;
                $sevenClockDiscount = new SevenClockDiscount();
                while ($day <= 31) {
                    $sevenClockDiscount->add($basePrice);
                    $sevenClockDiscount->add($basePrice);
                    echo "{'day':$day,'price':{$sevenClockDiscount->sum},'discount':'beforeSeven'},\n";
                    $day++;
                }

                $day = 1;
                $aliDiscount = new AliDiscount();
                while ($day <= 31) {
                    $aliDiscount->add($basePrice);
                    $aliDiscount->add($basePrice);
                    echo "{'day':$day,'price':{$aliDiscount->sum},'discount':'alipay'},\n";
                    $day++;
                }

                $day = 1;
                $jdDiscount = new JdDiscount();
                while ($day <= 31) {
                    $jdDiscount->add($basePrice);
                    $jdDiscount->add($basePrice);
                    echo "{'day':$day,'price':{$jdDiscount->sum},'discount':'jd'},\n";
                    $day++;
                }

                $day = 1;
                $noneDiscount = new Discount();
                while ($day <= 31) {
                    $noneDiscount->add($basePrice);
                    $noneDiscount->add($basePrice);
                    echo "{'day':$day,'price':{$noneDiscount->sum},'discount':'none'},\n";
                    $day++;
                }
            } else {
                $price = 3;
                while ($price <= $maxPrice) {
                    $day = 1;
                    $noneDiscount = new Discount();
                    while ($day <= 31) {
                        $noneDiscount->add($basePrice);
                        $noneDiscount->add($basePrice);
                        echo "{'day':$day,'price':{$noneDiscount->sum},'Base price':'$price'},\n";
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
