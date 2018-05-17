package main

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

// 消费者
type ConsumerT struct{}

func main() {
	//InitConsumer("test", "test-channel", "10.4.123.218:4161")
	InitConsumer("mytopic", "test-channel", "10.4.123.218:4161")
	for {
		time.Sleep(time.Second * 10)
	}
}

// 处理消息
func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	fmt.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	return nil
}

// 初始化消费者
func InitConsumer(topic string, channel string, address string) {
	cfg := nsq.NewConfig()
	// 最大允许向两台NSQD服务器接受消息，默认是1，要特别注意
	//cfg.MaxInFlight = 2
	cfg.LookupdPollInterval = time.Second          // 设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		panic(err)
	}
	c.SetLogger(nil, 0)        // 屏蔽系统日志
	c.AddHandler(&ConsumerT{}) // 添加消费者接口

	/*
		// 对消息进行处理的具体方法
		receive:=func(msg *nsq.Message)error{
			fmt.Println(string(msg.Body)
			return nil
		}
		// 添加消息处理的具体实现
		c.AddHandler(nsq.HandlerFunc(receive))
	*/

	// 将消费者连接到具体的NSQD
	//if err := c1.ConnectToNSQD("127.0.0.1:4150"); err != nil {
	//  panic(err)
	//}
	// 或者，如果启动了Lookupd服务，可通过nsqlookupd再分发给具体的nsqd
	// 建立NSQLookupd连接
	if err := c.ConnectToNSQLookupd(address); err != nil {
		panic(err)
	}

	// 建立多个nsqd连接
	// if err := c.ConnectToNSQDs([]string{"127.0.0.1:4150", "127.0.0.1:4152"}); err != nil {
	//  panic(err)
	// }

	// 建立一个nsqd连接
	// if err := c.ConnectToNSQD("127.0.0.1:4150"); err != nil {
	//  panic(err)
	// }
}

func produce(tag string, addr string) {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second * 5)
		p.Publish("test", []byte(tag+":"+time.Now().String()))
	}
}
func consume() {
	config := nsq.NewConfig()
	// 注意MaxInFlight的设置，默认只能接受一个节点
	config.MaxInFlight = 2
	c, err := nsq.NewConsumer("test", "consum", config)
	if err != nil {
		panic(err)
	}
	hand := func(msg *nsq.Message) error {
		fmt.Println(string(msg.Body))
		return nil
	}
	c.AddHandler(nsq.HandlerFunc(hand))
	if err := c.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fmt.Println(err)
	}
}
