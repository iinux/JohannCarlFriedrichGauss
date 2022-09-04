package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main() {
	iface := flag.String("i", "", "interface name")
	filePath := flag.String("f", "", "check file")
	flag.Parse()

	if *iface == "" {
		// 得到所有的(网络)设备
		devices, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatal(err)
		}
		// 打印设备信息
		fmt.Println("Devices found:")
		for _, device := range devices {
			fmt.Printf("Name: %s, Desc: %s\n", device.Name, device.Description)
			for _, address := range device.Addresses {
				fmt.Println("- IP address: ", address.IP)
				fmt.Println("- Subnet mask: ", address.Netmask)
			}
		}
	} else {
		// 打开某一网络设备
		handle, err = pcap.OpenLive(*iface, snapshotLen, promiscuous, timeout)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()
		// Use the handle as a packet source to process all packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			// Process packet here
			dnsLayer := packet.Layer(layers.LayerTypeDNS)
			if dnsLayer != nil {
				// fmt.Println(packet)
				dns, _ := dnsLayer.(*layers.DNS)
				if len(dns.Answers) > 0 {
					for _, answer := range dns.Answers {
						name := string(answer.Name)
						if name == "www.jetbrains.com" || name == "account.jetbrains.com" {
							cfr := checkFile(answer.IP.String(), *filePath)
							if cfr {
								fmt.Printf("block drop from any to %s\n", answer.IP)
							}
						}
					}
				}

			}
		}
	}
}

func checkFile(ip, filePath string) bool {
	if filePath == "" {
		return true
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return true
	}
	fileContent := string(b)
	if strings.Contains(fileContent, " "+ip+"\n") {
		return false
	} else if strings.Contains(fileContent, " "+ip+"\r") {
		return false
	} else if strings.Contains(fileContent, " "+ip+" ") {
		return false
	} else {
		return true
	}
}

var (
	snapshotLen int32         = 1024
	promiscuous bool          = false
	timeout     time.Duration = -1 * time.Second
	handle      *pcap.Handle
	err         error
)
