package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func lineNo() {
	_, file, line, _ := runtime.Caller(1) // 0 表示当前函数的调用信息
	fmt.Printf("当前行号: %s:%d\n", file, line)
}

func main() {
	// 配置抓包参数
	device := "en0" // 网卡名称，根据实际情况修改
	//filter := "tcp port 443" // 过滤443端口的TCP流量
	//filter := "tcp"
	snapshotLen := int32(65536)
	promiscuous := false
	timeout := 30 * time.Second

	// 打开网卡
	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal("打开网卡失败:", err)
	}
	defer handle.Close()

	// 设置BPF过滤器
	//if err := handle.SetBPFFilter(filter); err != nil {
	//	log.Fatal("设置过滤器失败:", err)
	//}

	fmt.Println("开始捕获HTTPS Client Hello消息...")
	fmt.Println("按Ctrl+C停止")

	// 创建数据包源
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	//packetSource.NoCopy = true
	//ctx, cancel := context.WithCancel(context.Background())
	//_ = ctx
	//defer cancel()

	// 处理数据包
	for packet := range packetSource.Packets() {
		processPacket(packet)
	}
}

func processPacket(packet gopacket.Packet) {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	// 获取IP层
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	ip6Layer := packet.Layer(layers.LayerTypeIPv6)
	if ipLayer == nil && ip6Layer == nil {
		return
	}
	ip, _ := ipLayer.(*layers.IPv4)
	ip6, _ := ip6Layer.(*layers.IPv6)

	if ip != nil && ip.FragOffset > 0 {
		fmt.Printf("%s\n", ip.DstIP)
	}

	// 获取TCP层
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}
	tcp, _ := tcpLayer.(*layers.TCP)

	// 检查是否包含应用层数据
	if len(tcp.Payload) == 0 {
		return
	}

	// 检查是否为TLS Client Hello (内容类型0x16，握手类型0x01)
	if tcp.Payload[0] != 0x16 {
		return
	}
	if len(tcp.Payload) < 5 {
		return
	}
	if tcp.Payload[5] != 0x01 { // Client Hello
		return
	}

	if tcp.DstPort == 80 {
		return
	}

	// 提取SNI
	sni := extractSNI(tcp.Payload)
	if sni == "" {
		if ip != nil {
			fmt.Printf("[-] %d %s %d %d %s %s %s %s\n",
				ip.Id,
				ip.Flags.String(),
				ip.FragOffset,
				ip.Length,
				ip.SrcIP,
				ip.DstIP,
				tcp.SrcPort,
				tcp.DstPort)
		} else if ip6 != nil {
			fmt.Printf("[-] %s %s %s %s\n",
				ip6.SrcIP,
				ip6.DstIP,
				tcp.SrcPort,
				tcp.DstPort)
		}
		return
	}

	// 打印结果
	if ip != nil {
		fmt.Printf("[+] Client Hello 检测到 - 去向 IP: %s, SNI: %s\n", ip.DstIP, sni)
	} else {
		fmt.Printf("[+] Client Hello 检测到 - 去向 IP: %s, SNI: %s\n", ip6.DstIP, sni)
	}
}

func extractSNI(payload []byte) string {
	// TLS记录层：跳过5字节头(ContentType[1], Version[2], Length[2])
	if len(payload) < 5 {
		return ""
	}
	recordLength := (int(payload[3])<<8 | int(payload[4]))
	if len(payload) < 5+recordLength {
		// return ""
	}
	// handshake := payload[5 : 5+recordLength]
	handshake := payload[5:]

	// 握手消息：跳过4字节头(HandshakeType[1], Length[3])
	if len(handshake) < 4 {
		return ""
	}
	handshakeLength := int(handshake[1])<<16 | int(handshake[2])<<8 | int(handshake[3])
	if len(handshake) < 4+handshakeLength {
		// return ""
	}
	// clientHello := handshake[4 : 4+handshakeLength]
	clientHello := handshake[4:]

	// 跳过随机数(32字节)和会话ID
	if len(clientHello) < 38 {
		return ""
	}
	sessionIDLength := int(clientHello[34])
	if len(clientHello) < 33+sessionIDLength {
		return ""
	}
	cipherSuites := clientHello[35+sessionIDLength:]

	// 跳过密码套件
	if len(cipherSuites) < 2 {
		return ""
	}
	cipherSuitesLength := int(cipherSuites[0])<<8 | int(cipherSuites[1])
	if len(cipherSuites) < 2+cipherSuitesLength+1 {
		return ""
	}
	compressionMethodsLength := int(cipherSuites[2+cipherSuitesLength])
	extensions := cipherSuites[2+cipherSuitesLength+1+compressionMethodsLength:]

	// 解析扩展
	if len(extensions) < 2 {
		return ""
	}
	// extensionsLength := int(extensions[0])<<8 | int(extensions[1])
	// extensions = extensions[2 : 2+extensionsLength]
	extensions = extensions[2:]

	// 查找SNI扩展(类型0x0000)
	for len(extensions) >= 4 {
		extType := (int(extensions[0]) << 8) | int(extensions[1])
		extLength := (int(extensions[2]) << 8) | int(extensions[3])

		if len(extensions) < 4+extLength {
			break
		}

		if extType == 0x0000 { // SNI扩展
			sniList := extensions[4 : 4+extLength]
			if len(sniList) < 3 {
				break
			}
			nameType := sniList[2]
			nameLength := int(sniList[3])<<8 | int(sniList[4])
			if len(sniList) < 3+nameLength {
				break
			}
			if nameType == 0x00 {
				return string(sniList[5 : 5+nameLength])
			}
			break
		}

		extensions = extensions[4+extLength:]
	}

	return ""
}
