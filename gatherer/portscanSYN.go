package gatherer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	snapshot_len int32         = 1024
	promiscuous  bool          = false
	timeout      time.Duration = -1 * time.Second
	//packageconut int           = 0
	//handle       *pcap.Handle
	filter string
	device string
)

var ipandports = make(chan map[string]string, 200)
var results = make(chan map[string]string, 200)

func finddevice() {
	fmt.Printf("[*]%c[1;32;40m%s%c[0m\n", 0x1B, "scanning,please waiting.....", 0x1B)
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[*]found device")
	for _, device := range devices {
		fmt.Println("***************************")
		fmt.Println("[*]", device.Name)
		fmt.Println("[*]", device.Description)
		for _, addr := range device.Addresses {
			fmt.Println("[*]", addr.IP)
			fmt.Println("[*]", addr.Broadaddr)
			fmt.Println("[*]", addr.Netmask)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("[*]输入网络设备")
	dev, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	device = string(dev)
}

func captureflow() {
	var ipandport = make(map[string]string)
	handle, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatal()
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range source.Packets() {
		networklayer := packet.NetworkLayer()
		srcHost := networklayer.NetworkFlow().Src().String()
		transportlayer := packet.TransportLayer()
		srcPort := transportlayer.TransportFlow().Src().String()
		ipandport[srcHost] = srcPort
		results <- ipandport
	}
}

func ports() []string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("[*]输入端口,端口段为以下格式,如:10-65535;或者使用:10,15,443")
	Ports, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	var portslice []string
	if strings.Contains(string(Ports), "-") {
		Portslice := strings.Split(string(Ports), "-")
		fmt.Println(Ports)
		fmt.Println(Portslice)
		port1, _ := strconv.Atoi(Portslice[0])
		port2, _ := strconv.Atoi(Portslice[1])
		for i := port1; i <= port2; i++ {
			portslice = append(portslice, strconv.Itoa(i))
		}
	} else if strings.Contains(string(Ports), ",") {
		portslice = strings.Split(string(Ports), "-")
		fmt.Println(Ports)
		fmt.Println(portslice)
	} else if strings.Contains(string(Ports), "-") && strings.Contains(string(Ports), ",") {
		fmt.Println("[*]暂时不支持这样的格式")
	} else {
		fmt.Println("[*]IP:", string(Ports))
		portslice = append(portslice, string(Ports))
	}
	return portslice
}

func ips() []string {
	var ipslice []string
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("[*]输入IP或者IP段,IP段为以下格式,如:10.1.100.1-10.1.100.254;或者使用:10.1.100.1,10.1.100.2,目前只支持10.1.100.1/16范围的")
	IPs, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(string(IPs), "-") {
		IPslice := strings.Split(string(IPs), "-")
		fmt.Println(IPs)
		fmt.Println(IPslice)
		ips1 := strings.Split(IPslice[0], ".")
		ips2 := strings.Split(IPslice[1], ".")
		if ips1[0] != ips2[0] || ips1[1] != ips2[1] {
			fmt.Println("[*]请查看输入IP网段的范围")
		} else if ips1[2] == ips2[2] {
			IPc1, _ := strconv.Atoi(ips1[3])
			IPc2, _ := strconv.Atoi(ips2[3])
			for i := IPc1; i <= IPc2; i++ {
				newIP := ips1[0] + "." + ips1[1] + "." + ips1[2] + "." + strconv.Itoa(i)
				ipslice = append(ipslice, newIP)
			}
		} else if ips1[2] != ips2[2] {
			IPc1, _ := strconv.Atoi(ips1[3])
			IPc2, _ := strconv.Atoi(ips2[3])
			IPb1, _ := strconv.Atoi(ips1[2])
			IPb2, _ := strconv.Atoi(ips2[3])
			for i := IPb1; i <= IPb2; i++ {
				for t := IPc1; t <= IPc2; i++ {
					newIP := ips1[0] + "." + ips1[1] + "." + strconv.Itoa(i) + "." + strconv.Itoa(t)
					ipslice = append(ipslice, newIP)
				}
			}
		}
	} else if strings.Contains(string(IPs), ",") {
		ipslice = strings.Split(string(IPs), ",")
		fmt.Println(IPs)
		fmt.Println(ipslice)
	} else if strings.Contains(string(IPs), "-") && strings.Contains(string(IPs), ",") {
		fmt.Println("[*]暂时不支持这样的格式")
	} else {
		fmt.Println("[*]IP:", string(IPs))
		ipslice = append(ipslice, string(IPs))
	}

	return ipslice
}

func workers(ipandports chan map[string]string, results chan map[string]string) {
	for ipandport := range ipandports {
		addres := ipandport["ip"] + ":" + ipandport["port"]
		conn, err := net.DialTimeout("tcp", addres, 1000*time.Millisecond)
		if err != nil {
			log.Fatal(err)
		}
		conn.Close()
	}
}

func getFilter(ips, ports []string) {
	var portFilter string
	var ipFilter string
	if len(ips) == 1 {
		ipFilter = " port " + ips[0]
	} else {
		for key, ip := range ips {
			if key == 0 {
				ipFilter += ("(port " + ip)
			}
			ipFilter += (" or port " + ip)
		}
		ipFilter += ")"
	}

	if len(ports) == 1 {
		portFilter = " port " + ports[0]
	} else {
		for key, ports := range ports {
			if key == 0 {
				portFilter += ("(port " + ports)
			}
			portFilter += (" or port " + ports)
		}
		portFilter += ")"
	}

	filter = fmt.Sprintf("tcp and %s and %s and (tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18)", ipFilter, portFilter)
	fmt.Println(filter)
}

func Scan_SYN() {
	var opens []map[string]string
	var map1 map[string]string
	finddevice()
	iiip := ips()
	ppport := ports()
	getFilter(iiip, ppport)

	go captureflow()

	go func() {
		for _, val := range iiip {
			for _, val1 := range ppport {
				map1["ip"] = val
				map1["port"] = val1
				ipandports <- map1
			}
		}
	}()

	for i := 0; i < cap(ipandports); i++ {
		go workers(ipandports, results)
	}

	for rel := range results {
		opens = append(opens, rel)
	}
	close(ipandports)
	close(results)

	if len(opens) == 0 {
		fmt.Println("[*]没有找到开放端口")
	} else {
		for _, ipandport1 := range opens {
			fmt.Printf("[*]%s:%s", ipandport1["ip"], ipandport1["port"])
		}
	}
}
