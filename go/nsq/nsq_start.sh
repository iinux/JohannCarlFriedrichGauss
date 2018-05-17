#!/usr/bin/env bash

#服务启动
#注意更改一下 --data-path 所指定的数据存放路径，否则会无法运行。
echo '删除日志文件'
rm -f nsqlookupd.log
rm -f nsqd1.log
rm -f nsqd2.log
rm -f nsqadmin.log

echo '启动nsq服务'
nohup nsqlookupd >nsqlookupd.log 2>&1&

echo '启动nsqd服务'
nohup nsqd --lookupd-tcp-address=0.0.0.0:4160 -tcp-address="0.0.0.0:4150"  --data-path=~/nsqd1  >nsqd1.log 2>&1&
nohup nsqd --lookupd-tcp-address=0.0.0.0:4160 -tcp-address="0.0.0.0:4152" -http-address="0.0.0.0:4153" --data-path=~/nsqd2 >nsqd2.log 2>&1&

echo '启动nsqdadmin服务'
nohup nsqadmin --lookupd-http-address=0.0.0.0:4161 >nsqadmin.log 2>&1&