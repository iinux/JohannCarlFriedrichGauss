package main

import (
	"bytes"
	"encoding/binary"
	//"flag"
	. "fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type TCPHeader struct {
	SrcPort   uint16
	DstPort   uint16
	SeqNum    uint32
	AckNum    uint32
	Offset    uint8
	Flag      uint8
	Window    uint16
	Checksum  uint16
	UrgentPtr uint16
}

type PsdHeader struct {
	SrcAddr   uint32
	DstAddr   uint32
	Zero      uint8
	ProtoType uint8
	TcpLength uint16
}

func inetAddr(ipaddr string) uint32 {
	var (
		segments []string = strings.Split(ipaddr, ".")
		ip       [4]uint64
		ret      uint64
	)
	for i := 0; i < 4; i++ {
		ip[i], _ = strconv.ParseUint(segments[i], 10, 64)
	}
	ret = ip[3]<<24 + ip[2]<<16 + ip[1]<<8 + ip[0]
	return uint32(ret)
}

func htons(port uint16) uint16 {
	var (
		high uint16 = port >> 8
		ret  uint16 = port<<8 + high
	)
	return ret
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
	/*
	var si = flag.String("si", "", "source ip")
	var sp = flag.Int("sp", 0, "source port")
	var di = flag.String("di", "", "destination ip")
	var dp = flag.Int("dp", 0, "destination port")

	flag.Parse()

	 */
	si := os.Args[1]
	sp, _ := strconv.ParseUint(os.Args[2], 10, 16)
	di := os.Args[3]
	dp, _ := strconv.ParseUint(os.Args[4], 10, 16)
	seq, _ := strconv.ParseUint(os.Args[5], 10, 32)
	ack, _ := strconv.ParseUint(os.Args[6], 10, 32)

	var (
		msg       string
		psdHeader PsdHeader
		tcpHeader TCPHeader
	)

	// Printf("Input the content: ")
	// Scanf("%s", &msg)
	msg = "ok"

	/*填充TCP伪首部*/
	psdHeader.SrcAddr = inetAddr(si)
	psdHeader.DstAddr = inetAddr(di)
	psdHeader.Zero = 0
	psdHeader.ProtoType = syscall.IPPROTO_TCP
	psdHeader.TcpLength = uint16(unsafe.Sizeof(TCPHeader{})) + uint16(len(msg))

	/*填充TCP首部*/
	tcpHeader.SrcPort = uint16(sp)
	tcpHeader.DstPort = uint16(dp)
	tcpHeader.SeqNum = uint32(seq)
	tcpHeader.AckNum = uint32(ack)
	tcpHeader.Offset = uint8(uint16(unsafe.Sizeof(TCPHeader{}))/4) << 4
	tcpHeader.Flag = 2 //SYN
	tcpHeader.Flag = 4 //RST
	tcpHeader.Window = 60000
	tcpHeader.Checksum = 0

	/*buffer用来写入两种首部来求得校验和*/
	var (
		buffer bytes.Buffer
	)
	_ = binary.Write(&buffer, binary.BigEndian, psdHeader)
	_ = binary.Write(&buffer, binary.BigEndian, tcpHeader)
	tcpHeader.Checksum = CheckSum(buffer.Bytes())

	/*接下来清空buffer，填充实际要发送的部分*/
	buffer.Reset()
	_ = binary.Write(&buffer, binary.BigEndian, tcpHeader)
	_ = binary.Write(&buffer, binary.BigEndian, msg)

	/*下面的操作都是raw socket操作，大家都看得懂*/
	var (
		sockFd int
		addr   syscall.SockaddrInet4
		err    error
	)


	if sockFd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP); err != nil {
		Println("Socket() error: ", err.Error())
		return
	}
	defer syscall.Shutdown(sockFd, syscall.SHUT_RDWR)

	addr.Addr[3], addr.Addr[2], addr.Addr[1], addr.Addr[0] = byte(psdHeader.DstAddr&(0xff<<24)>>24),
		byte(psdHeader.DstAddr&(0xff<<16)>>16), byte(psdHeader.DstAddr&(0xff<<8)>>8), byte(psdHeader.DstAddr&0xff)
	addr.Port = int(tcpHeader.DstPort)

	if err = syscall.Sendto(sockFd, buffer.Bytes(), 0, &addr); err != nil {
		Println("Send to() error: ", err.Error())
		return
	}

	Println("Send success!")
}
