package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func main() {
	var (
		icmp  ICMP
		laddr net.IPAddr = net.IPAddr{IP: net.ParseIP("192.168.22.148")}  //***IP地址改成你自己的网段***
		raddr net.IPAddr = net.IPAddr{IP: net.ParseIP("192.168.22.1")}
	)
	//如果你要使用网络层的其他协议还可以设置成 ip:ospf、ip:arp 等
	conn, err := net.DialIP("ip4:icmp", &laddr, &raddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	//开始填充数据包
	icmp.Type = 8 //8->echo message  0->reply message
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	var (
		buffer bytes.Buffer
	)
	//先在buffer中写入icmp数据报求去校验和
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	//然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("send icmp packet success!")
}
