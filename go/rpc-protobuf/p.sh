#!/bin/sh
protoc -I customer/ customer/customer.proto --go_out=plugins=grpc:customer
# refer http://www.jianshu.com/p/3139e8dd4dd1?hmsr=toutiao.io&utm_medium=toutiao.io&utm_source=toutiao.io
