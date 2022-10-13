package quantityofflow

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

func Flowfilter(handle *pcap.Handle) *pcap.Handle {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("[*]输入过滤规则,如:tcp and port 80")
	filter, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter(string(filter))
	if err != nil {
		log.Fatal(err)
	}
	return handle
}

func FindDevice() {
	devices, err := pcap.FindAllDevs()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n[*]devices found")

	for _, device := range devices {
		fmt.Println("-----------------")
		fmt.Printf("[*]%c[1;31;32m%s%c[0m\n", 0x1B, "Deivce Name:"+device.Name, 0x1B)
		fmt.Println("[*]device Description:\t" + device.Description)
		for _, address := range device.Addresses {
			fmt.Println("[*]ip:", address.IP)
			fmt.Println("[*]net mask:", address.Netmask)
		}
	}
}

var (
	snapshot_len int32         = 1024
	promiscuous  bool          = false
	timeout      time.Duration = -1 * time.Second
	packageconut int           = 0
	handle       *pcap.Handle
	pcapfile     string = "test.pcap"
	// rel          string
)

// type Analysis struct{
// 	snapshot_len int32         = 1024
// 	promiscuous  bool          = false
// 	timeout      time.Duration = -1 * time.Second
// 	packageconut int           = 0
// 	handle       *pcap.Handle
// 	pcapfile     string = "test1.pcap"
// }

var source = make(chan gopacket.Packet, 20)

func GetFlow(filer1, parse1, write1 string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("[*]输入网络设备的名称：")
	device, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(device))
	if handle, err = pcap.OpenLive(string(device), snapshot_len, promiscuous, timeout); err != nil {
		log.Fatal(err)
	} else {
		defer handle.Close()
		if filer1 == "yes" {
			fmt.Println("======过滤....")
			handle = Flowfilter(handle)
		}
		packagesource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packagesource.Packets() {
			// fmt.Printf("\n%s\n", packet)
			if write1 == "yes" {
				source <- packet
			}

			if parse1 == "yes" {
				Parsepcap(packet)
			}
		}
	}
}

func Writepcap() {
	f, _ := os.Create(pcapfile)
	w := pcapgo.NewWriter(f)

	w.WriteFileHeader(uint32(snapshot_len), layers.LinkTypeEthernet)
	defer f.Close()

	for {
		a := <-source
		fmt.Println(a)
		w.WritePacket(a.Metadata().CaptureInfo, a.Data())
		packageconut++
	}
}

func Readpcap() {
	handle1, err := pcap.OpenOffline(pcapfile)
	if err != nil {
		log.Fatal(err)
	}

	pcapfilesource := gopacket.NewPacketSource(handle1, handle1.LinkType())
	for packet := range pcapfilesource.Packets() {
		Parsepcap(packet)
	}
}

func Parsepcap(packet gopacket.Packet) {
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetpacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("==========================")
		fmt.Println("[*]ethernet packet dected")
		fmt.Printf("[*]source MAC: %s\n", ethernetpacket.SrcMAC)
		fmt.Printf("[*]Destination MAC: %s\n", ethernetpacket.DstMAC)
		fmt.Printf("[*]Ethernet type: %s\n", ethernetpacket.EthernetType)
		fmt.Printf("[*]Ethernet length: %d\n", ethernetpacket.Length)
		// fmt.Println("[*]payload: ", ethernetpacket.Payload)
		fmt.Println()
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("==========================")
		fmt.Println("[*]IPv4 dected")
		ipPacket, _ := ipLayer.(*layers.IPv4)
		fmt.Printf("[*]from %s to %s\n", ipPacket.SrcIP, ipPacket.DstIP)
		fmt.Printf("[*]Protocol: %s\n", ipPacket.Protocol)
		fmt.Println("[*]Contets: ", ipPacket.Contents)
		fmt.Println()
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("==========================")
		fmt.Println("[*]TCP dected")
		tcpPacket, _ := tcpLayer.(*layers.TCP)
		fmt.Printf("[*]From port %s to port %s\n", tcpPacket.SrcPort, tcpPacket.DstPort)
	}

	// aLayer := packet.Layer(layers.LayerType)

	fmt.Println("[*]All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}
	app1 := packet.ApplicationLayer()
	if app1 != nil {
		fmt.Println("[*]payload:", string(app1.Payload()))
		// fmt.Println(app1.)
		if strings.Contains(string(app1.Payload()), "sqlmap") {
			fmt.Printf("[*]%c[1;31;41m%s%c[0m\n", 0x1B, "SQLmap 探测", 0x1B)
		} else {
			rule_sql(string(app1.Payload()))
		}
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}

func rule_sql(payload string) {
	var rel string
	sql_keywords := "txt\\sql_keywords.txt"
	file, err := os.Open(sql_keywords)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	fd := bufio.NewScanner(file)

	var rule_array = []string{`={1}[\d+%*[A-z]*]*-{2}`, `={1}\d+%*\d+\){1}[\d+%*[A-z]*]*-{2}`, `={1}.*-{2}`, `={1}.*={1}`}
	for _, rule_choose := range rule_array {
		reg1 := regexp.MustCompile(rule_choose)
		result := reg1.FindAllStringSubmatch(payload, -1)
		if result == nil {
			fmt.Println("规则未匹配到特殊字符")
		} else {
			for _, val := range result {
				for _, val1 := range val {
					rel += val1
				}
			}
		}

		rel2, err := url.QueryUnescape(rel)

		if err != nil {
			fmt.Printf("[*]%c[1;33;40m%s%c[0m\n", 0x1B, "payload中有特殊字符未能解码:"+rel, 0x1B)
			// return
		}

		fmt.Println(rel)
		rel2 = strings.ToUpper(rel2)
		if rel != "" {
			for fd.Scan() {
				if strings.Contains(rel2, fd.Text()) {
					fmt.Printf("[*]%c[1;31;41m%s%c[0m\n", 0x1B, "疑似SQL注入 payload:"+rel, 0x1B)
					fmt.Printf("[*]%c[1;31;41m%s%c[0m\n", 0x1B, "疑似SQL注入 解码后payload:"+rel2, 0x1B)
					fmt.Printf("[*]%c[1;31;41m%s%c[0m\n", 0x1B, "疑似SQL注入 关键字:"+fd.Text(), 0x1B)
					break
				}
			}
			break
		}
	}
}

// func rule_xss() {

// }

func rule_ftp() {

}

func CreateAndSendPacket() {

}
