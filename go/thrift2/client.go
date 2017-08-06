package main

import (
	"demo/rpc"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"os"
	"time"
)

func main() {
	startTime := currentTimeMillis()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "19090"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}

	useTransport, err := transportFactory.GetTransport(transport)
	client := rpc.NewRpcServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:19090", " ", err)
		os.Exit(1)
	}
	defer transport.Close()

	for i := 0; i < 1000; i++ {
		paramMap := make(map[string]string)
		paramMap["name"] = "qinerg"
		paramMap["passwd"] = "123456"
		r1, e1 := client.FunCall(currentTimeMillis(), "login", paramMap)
		fmt.Println(i, "Call->", r1, e1)
	}

	endTime := currentTimeMillis()
	fmt.Println("Program exit. time->", endTime, startTime, (endTime - startTime))
}

// 转换成毫秒
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
