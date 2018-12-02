package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket"
)

func main1() {
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print device information
	fmt.Println("Devices found:")
	for _, d := range devices {
		fmt.Println("\nName: ", d.Name)
		fmt.Println("Description: ", d.Description)
		fmt.Println("Devices addresses: ", d.Description)

		for _, address := range d.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
}

func main() {
	// handle, err := pcap.OpenLive("\\Device\\NPF_{713C668E-58F6-4831-90A5-73FEEC913A39}", 1024, false, 30*time.Second)
	handle, err := pcap.OpenLive("\\Device\\NPF_{D7DE90C9-4C90-4894-848E-C1FA8A1EF79B}", 1024, false, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		log.Println(packet)
	}
}
