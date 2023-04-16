# trickle
使用 libc 才有效，使用ldd命令查看

# wondershaper
apt install wondershaper

# TC（Traffic Control）

```bash
yum -y install iproute-tc
yum -y install kernel-modules-extra

tc qdisc ls
tc qdisc ls dev wlp6s0
tc -s qdisc ls dev wlp6s0

tc qdisc show dev wlp6s0

tc qdisc add dev wlp6s0 root netem delay 1000ms
tc qdisc add dev wlp6s0 root netem delay 1000ms 500ms
tc qdisc add dev wlp6s0 root netem delay 1000ms 500ms 20%
tc qdisc add dev wlp6s0 root netem delay 1000ms 500ms distribution normal

tc qdisc add dev wlp6s0 root netem loss 50%
# 包损
tc qdisc add dev wlp6s0 root netem corrupt 50%
tc qdisc add dev wlp6s0 root netem duplicate 50%

tc qdisc del dev wlp6s0 root netem delay 1000ms
tc qdisc del dev wlp6s0 root netem
tc qdisc del dev wlp6s0 root

tc qdisc add dev wlp6s0 root tbf rate 1mbit burst  32kbit latency 400ms
tc qdisc add dev wlp6s0 root tbf rate 10kbit burst 10kbit limit 10kbit
tc qdisc change dev wlp6s0 root netem delay 1000ms
```

# unit
## 带宽或者流速单位：
* kbps                            千字节／秒
* mbps                           兆字节／秒
* kbit                             KBits／秒
* mbit                            MBits／秒
* bps或者一个无单位数字      字节数／秒

## 数据的数量单位：
* kb或者k                      千字节
* mb或者m                    兆字节
* mbit                          兆bit
* kbit                           千bit
* b或者一个无单位数字       字节数

## 时间的计量单位：
* s、sec或者secs                              秒
* ms、msec或者msecs                       分钟
* us、usec、usecs或者一个无单位数字    微秒

# action
* add
  * 在一个节点里加入一个QDisc、类或者过滤器。添加时，需要传递一个祖先作为参数，传递参数时既可以使用ID也可以直接传递设备的根。如果要建立一个QDisc或者过滤器，可以使用句柄(handle)来命名；如果要建立一个类，可以使用类识别符(classid)来命名。
* remove
  * 删除有某个句柄(handle)指定的QDisc，根QDisc(root)也可以删除。被删除QDisc上的所有子类以及附属于各个类的过滤器都会被自动删除。
* change
  * 以替代的方式修改某些条目。除了句柄(handle)和祖先不能修改以外，change命令的语法和add命令相同。
* replace
  * 对一个现有节点进行近于原子操作的删除／添加。如果节点不存在，这个命令就会建立节点。
* link
  * 只适用于DQisc，替代一个现有的节点。

# 流量控制包括以下几种方式：
## SHAPING(限制)
当流量被限制，它的传输速率就被控制在某个值以下。限制值可以大大小于有效带宽，这样可以平滑突发数据流量，使网络更为稳定。shaping（限制）只适用于向外的流量。
## SCHEDULING(调度)
通过调度数据包的传输，可以在带宽范围内，按照优先级分配带宽。SCHEDULING(调度)也只适于向外的流量。
## POLICING(策略)
SHAPING用于处理向外的流量，而POLICIING(策略)用于处理接收到的数据。
## DROPPING(丢弃)
如果流量超过某个设定的带宽，就丢弃数据包，不管是向内还是向外。

# 流量的处理由三种对象控制，它们是：
* qdisc(排队规则, queueing discipline)
* class(类别)
* filter(过滤器)

## QDISC的类别如下：
* CLASSLESS QDisc(不可分类QDisc)
  * [p|b]fifo
    * 使用最简单的qdisc，纯粹的先进先出。只有一个参数：limit，用来设置队列的长度,pfifo是以数据包的个数为单位；bfifo是以字节数为单位。
  * pfifo_fast
    * 在编译内核时，如果打开了高级路由器(Advanced Router)编译选项，pfifo_fast就是系统的标准QDISC。它的队列包括三个波段(band)。在每个波段里面，使用先进先出规则。而三个波段(band)的优先级也不相同，band 0的优先级最高，band 2的最低。如果band里面有数据包，系统就不会处理band 1里面的数据包，band 1和band 2之间也是一样。数据包是按照服务类型(Type of Service,TOS)被分配多三个波段(band)里面的。
  * red
    * red是Random Early Detection(随机早期探测)的简写。如果使用这种QDISC，当带宽的占用接近于规定的带宽时，系统会随机地丢弃一些数据包。它非常适合高带宽应用。
  * sfq
    * sfq是Stochastic Fairness Queueing的简写。它按照会话(session--对应于每个TCP连接或者UDP流)为流量进行排序，然后循环发送每个会话的数据包。
  * tbf
    * tbf是Token Bucket Filter的简写，适合于把流速降低到某个值。
* CLASSFUL QDISC(分类QDisc)
  * CBQ
    * CBQ是Class Based Queueing(基于类别排队)的缩写。它实现了一个丰富的连接共享类别结构，既有限制(shaping)带宽的能力，也具有带宽优先级管理的能力。带宽限制是通过计算连接的空闲时间完成的。空闲时间的计算标准是数据包离队事件的频率和下层连接(数据链路层)的带宽。
  * HTB
    * HTB是Hierarchy Token Bucket的缩写。通过在实践基础上的改进，它实现了一个丰富的连接共享类别体系。使用HTB可以很容易地保证每个类别的带宽，虽然它也允许特定的类可以突破带宽上限，占用别的类的带宽。HTB可以通过TBF(Token Bucket Filter)实现带宽限制，也能够划分类别的优先级。
  * PRIO
    * PRIO QDisc不能限制带宽，因为属于不同类别的数据包是顺序离队的。使用PRIO QDisc可以很容易对流量进行优先级管理，只有属于高优先级类别的数据包全部发送完毕，才会发送属于低优先级类别的数据包。为了方便管理，需要使用iptables或者ipchains处理数据包的服务类型(Type Of Service,ToS)。

## CLASS(类)
某些QDisc(排队规则)可以包含一些类别，不同的类别中可以包含更深入的QDisc(排队规则)，通过这些细分的QDisc还可以为进入的队列的数据包排队。通过设置各种类别数据包的离队次序，QDisc可以为设置网络数据流量的优先级。
## FILTER(过滤器)
Filter(过滤器)用于为数据包分类，决定它们按照何种QDisc进入队列。无论何时数据包进入一个划分子类的类别中，都需要进行分类。分类的方法可以有多种，使用fileter(过滤器)就是其中之一。使用filter(过滤器)分类时，内核会调用附属于这个类(class)的所有过滤器，直到返回一个判决。如果没有判决返回，就作进一步的处理，而处理方式和QDISC有关。需要注意的是，filter(过滤器)是在QDisc内部，它们不能作为主体。

## 具体操作
 Linux流量控制主要分为建立队列、建立分类和建立过滤器三个方面。
 ### 基本实现步骤为：
1. 针对网络物理设备（如以太网卡eth0）绑定一个队列QDisc；
1. 在该队列上建立分类class；
1. 为每一分类建立一个基于路由的过滤器filter；
1. 最后与过滤器相配合，建立特定的路由表。
### 环境模拟实例:
流量控制器上的以太网卡(eth0) 的IP地址为192.168.1.66，在其上建立一个CBQ队列。假设包的平均大小为1000字节，包间隔发送单元的大小为8字节，可接收冲突的发送最长包数目为20字节。

假如有三种类型的流量需要控制:
1. 是发往主机1的，其IP地址为192.168.1.24。其流量带宽控制在8Mbit，优先级为2；
1. 是发往主机2的，其IP地址为192.168.1.30。其流量带宽控制在1Mbit，优先级为1；
1. 是发往子网1的，其子网号为192.168.1.0，子网掩码为255.255.255.0。流量带宽控制在1Mbit，优先级为6。

#### 建立队列
一般情况下，针对一个网卡只需建立一个队列。

将一个cbq队列绑定到网络物理设备eth0上，其编号为1:0；网络物理设备eth0的实际带宽为10 Mbit，包的平均大小为1000字节；包间隔发送单元的大小为8字节，最小传输包大小为64字节。

`tc qdisc add dev eth0 root handle 1: cbq bandwidth 10Mbit avpkt 1000 cell 8 mpu 64`

#### 建立分类
分类建立在队列之上。

一般情况下，针对一个队列需建立一个根分类，然后再在其上建立子分类。对于分类，按其分类的编号顺序起作用，编号小的优先；一旦符合某个分类匹配规则，通过该分类发送数据包，则其后的分类不再起作用。

* 创建根分类1:1；分配带宽为10Mbit，优先级别为8。
  * `tc class add dev eth0 parent 1:0 classid 1:1 cbq bandwidth 10Mbit rate 10Mbit maxburst 20 allot 1514 prio 8 avpkt 1000 cell 8 weight 1Mbit`
  * 该队列的最大可用带宽为10Mbit，实际分配的带宽为10Mbit，可接收冲突的发送最长包数目为20字节；最大传输单元加MAC头的大小为1514字节，优先级别为8，包的平均大小为1000字节，包间隔发送单元的大小为8字节，相应于实际带宽的加权速率为1Mbit。
* 创建分类1:2，其父分类为1:1，分配带宽为8Mbit，优先级别为2。
  * `tc class add dev eth0 parent 1:1 classid 1:2 cbq bandwidth 10Mbit rate 8Mbit maxburst 20 allot 1514 prio 2 avpkt 1000 cell 8 weight 800Kbit split 1:0 bounded`
  * 该队列的最大可用带宽为10Mbit，实际分配的带宽为 8Mbit，可接收冲突的发送最长包数目为20字节；最大传输单元加MAC头的大小为1514字节，优先级别为1，包的平均大小为1000字节，包间隔发送单元的大小为8字节，相应于实际带宽的加权速率为800Kbit，分类的分离点为1:0，且不可借用未使用带宽。
* 创建分类1:3，其父分类为1:1，分配带宽为1Mbit，优先级别为1。
  * `tc class add dev eth0 parent 1:1 classid 1:3 cbq bandwidth 10Mbit rate 1Mbit maxburst 20 allot 1514 prio 1 avpkt 1000 cell 8 weight 100Kbit split 1:0`
  * 该队列的最大可用带宽为10Mbit，实际分配的带宽为 1Mbit，可接收冲突的发送最长包数目为20字节；最大传输单元加MAC头的大小为1514字节，优先级别为2，包的平均大小为1000字节，包间隔发送单元的大小为8字节，相应于实际带宽的加权速率为100Kbit，分类的分离点为1:0。
* 创建分类1:4，其父分类为1:1，分配带宽为1Mbit，优先级别为6。
  * `tc class add dev eth0 parent 1:1 classid 1:4 cbq bandwidth 10Mbit rate 1Mbit maxburst 20 allot 1514 prio 6 avpkt 1000 cell 8 weight 100Kbit split 1:0`
  * 该队列的最大可用带宽为10Mbit，实际分配的带宽为1Mbit，可接收冲突的发送最长包数目为20字节；最大传输单元加MAC头的大小为1514字节，优先级别为6，包的平均大小为1000字节，包间隔发送单元的大小为8字节，相应于实际带宽的加权速率为100Kbit，分类的分离点为1:0。
### 建立过滤器
过滤器主要服务于分类。

一般只需针对根分类提供一个过滤器，然后为每个子分类提供路由映射。

* 1） 应用路由分类器到cbq队列的根，父分类编号为1:0；过滤协议为ip，优先级别为100，过滤器为基于路由表。
  * `tc filter add dev eth0 parent 1:0 protocol ip prio 100 route`
* 2） 建立路由映射分类1:2, 1:3, 1:4
  * `tc filter add dev eth0 parent 1:0 protocol ip prio 100 route to 2 flowid 1:2`
  * `tc filter add dev eth0 parent 1:0 protocol ip prio 100 route to 3 flowid 1:3`
  * `tc filter add dev eth0 parent 1:0 protocol ip prio 100 route to 4 flowid 1:4`
### 建立路由
该路由是与前面所建立的路由映射一一对应。
* 1） 发往主机192.168.1.24的数据包通过分类2转发(分类2的速率8Mbit)
  * `ip route add 192.168.1.24 dev eth0 via 192.168.1.66 realm 2`
* 2） 发往主机192.168.1.30的数据包通过分类3转发(分类3的速率1Mbit)
  * `ip route add 192.168.1.30 dev eth0 via 192.168.1.66 realm 3`
* 3）发往子网192.168.1.0/24的数据包通过分类4转发(分类4的速率1Mbit)
  * `ip route add 192.168.1.0/24 dev eth0 via 192.168.1.66 realm 4`

注：一般对于流量控制器所直接连接的网段建议使用IP主机地址流量控制限制，不要使用子网流量控制限制。如一定需要对直连子网使用子网流量控制限制，则在建立该子网的路由映射前，需将原先由系统建立的路由删除，才可完成相应步骤。
### 监视
主要包括对现有队列、分类、过滤器和路由的状况进行监视。
* 1）显示队列的状况
  * 详细显示指定设备(这里为eth0)的队列状况
  * `tc -s qdisc ls dev eth0`
* 2）显示分类的状况
  * 简单显示指定设备(这里为eth0)的分类状况
  * `tc class ls dev eth0`
  * 详细显示指定设备(这里为eth0)的分类状况
  * `tc -s class ls dev eth0`
  * 这里主要显示了通过不同分类发送的数据包，数据流量，丢弃的包数目，超过速率限制的包数目等等。其中根分类(class cbq 1:0)的状况应与队列的状况类似。
* 3）显示过滤器的状况
  * `tc -s filter ls dev eth0`

### 维护

主要包括对队列、分类、过滤器和路由的增添、修改和删除。

增添动作一般依照"队列->分类->过滤器->路由"的顺序进行；修改动作则没有什么要求；删除则依照"路由->过滤器->分类->队列"的顺序进行。

#### 1）队列的维护

一般对于一台流量控制器来说，出厂时针对每个以太网卡均已配置好一个队列了，通常情况下对队列无需进行增添、修改和删除动作了。

#### 2）分类的维护

* 增添
  * 增添动作通过tc class add命令实现，如前面所示。
* 修改
  * 修改动作通过tc class change命令实现，如下所示：
  * tc class change dev eth0 parent 1:1 classid 1:2 cbq bandwidth 10Mbit rate 7Mbit maxburst 20 allot 1514 prio 2 avpkt 1000 cell 8 weight 700Kbit split 1:0 bounded
  * 对于bounded命令应慎用，一旦添加后就进行修改，只可通过删除后再添加来实现。
* 删除
  * 删除动作只在该分类没有工作前才可进行，一旦通过该分类发送过数据，则无法删除它了。因此，需要通过shell文件方式来修改，通过重新启动来完成删除动作。

#### 3）过滤器的维护

增添

增添动作通过tc filter add命令实现，如前面所示。

修改

修改动作通过tc filter change命令实现，如下所示：

·tc filter change dev eth0 parent 1:0 protocol ip prio 100 route to 10 flowid 1:8

删除

删除动作通过tc filter del命令实现，如下所示：

·tc filter del dev eth0 parent 1:0 protocol ip prio 100 route to 10

#### 4）与过滤器一一映射路由的维护

增添

增添动作通过ip route add命令实现，如前面所示。

修改

修改动作通过ip route change命令实现，如下所示：

·ip route change 192.168.1.30 dev eth0 via 192.168.1.66 realm 8

删除

删除动作通过ip route del命令实现，如下所示：

·ip route del 192.168.1.30 dev eth0 via 192.168.1.66 realm 8

·ip route del 192.168.1.0/24 dev eth0 via 192.168.1.66 realm 4

```bash
# ！/bin/sh
touch  /var/lock/subsys/local
echo  1  > /proc/sys/net/ipv4/ip_forward  #（激活转发）

route add default  gw  10.0.0.0  #(这是加入电信网关，如果你已设了不用这条）

DOWNLOAD=640Kbit    #(640/8 =80K ,我这里限制下载最高速度只能80K）
UPLOAD=640Kbit          #(640/8 =80K,上传速度也限制在80K）
INET=192.168.0.          #(设置网段，根据你的情况填）
IPS=1                          #(这个意思是从192.168.0.1开始）
IPE=200                        #(我这设置是从IP为192.168.0.1-200这个网段限速，根据自已的需要改）
ServerIP=253                #(网关IP）
IDEV=eth0
ODEV=eth1

tc  qdisc  del  dev  $IDEV root handle 10:
tc  qdisc  del  dev  $ODEV  root handle  20:
tc  qdisc  add  dev $IDEV  root  handle  10: cbq  bandwidth  100Mbit avpkt  1000
tc  qdisc  add  dev  $ODEV  root  handle  20: cbq bandwidth  1Mbit  avpkt  1000
tc  class  add  dev $IDEV  parent 10:0  classid  10:1  cbq  bandwidth  100Mbit  rate 100Mbit  allot 1514  weight  1Mbit  prio  8  maxburst  20  avpkt 1000
tc  class  add  dev  $ODEV  parent  20:0  classid  20:1 cbq  bandwidth  1Mbit  rate  1Mbit  allot  1514  weitht  10Kbit  prio  8  maxburst  20  avpkt 1000

COUNTER=$IPS
while  [  $COUNTER  -le  $IPE  ]
do
    tc  class  add  dev  $IDEV  parent  10:1  classid  10:1$COUNTER  cbq  banwidth  100Mbit  rate $DOWNLOAD  allot  1514  weight  20Kbit  prio  5  maxburst  20  avpkt  1000  bounded
    tc  qdisc  add  dev  $IDEV  parent  10:1$COUNTER  sfq  quantum  1514b  perturb15
    tc  filter  add  dev  $IDEV  parent  10:0  protocol  ip  prio  100  u32  match  ipdst  $INET$COUNTER  flowid  10:1$COUNTER
    COUNTER=`expr $COUNTER + 1`
done

iptables  -t  nat  -A  POSTROUTING  -o  eth1  -s  192.168.0.0/24  -J  MASQUERADE
```


```bash
#!/bin/sh
tc qdisc del dev eth7 root &> /dev/null
tc qdisc del dev eth8 root &> /dev/null

#Add qdisc
tc qdisc add dev eth7 root handle 10: htb default 9998
tc qdisc add dev eth8 root handle 10: htb default 9998

#Add htb root node
tc class add dev eth7 parent 10: classid 10:9999 htb rate 1000000kbit ceil 1000000kbit
tc class add dev eth8 parent 10: classid 10:9999 htb rate 1000000kbit ceil 1000000kbit

#Add htb fake default node here
tc class add dev eth7 parent 10:9999 classid 10:9998 htb rate 1000000kbit ceil 1000000kbit
tc class add dev eth8 parent 10:9999 classid 10:9998 htb rate 1000000kbit ceil 1000000kbit

#Add rule node
tc class add dev eth7 parent 10:9999 classid 10:3 htb rate 1kbit ceil 50kbit
tc filter add dev eth7 parent 10: protocol ip handle 3 fw classid 10:3
tc class add dev eth8 parent 10:9999 classid 10:3 htb rate 1kbit ceil 50kbit
tc filter add dev eth8 parent 10: protocol ip handle 3 fw classid 10:3

#Add htb real default node here
tc class change dev eth7 classid 10:9998 htb rate 1kbit ceil 1000000kbit
tc class change dev eth8 classid 10:9998 htb rate 1kbit ceil 1000000kbit
```
