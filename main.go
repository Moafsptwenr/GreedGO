package main

import (
	_ "AssassinateGO/exp/nc"
	"AssassinateGO/explore"
	"AssassinateGO/gatherer"
	"AssassinateGO/quantityofflow"
	"AssassinateGO/rob"
	"flag"
	"fmt"
)

var (
	GOOS          string
	GOARCH        string
	CGO_ENABLED   *int
	types         string
	FLOWFILTER    string
	FLOW          string
	WRITEPCAP     string
	READPCAP      string
	FINDDEVICE    string
	GETFLOW       string
	PARSEFLOW     string
	EXPLORE       string
	LOGIN_EXPLORE string
	GATHERER      string
	PORT_SCAN     string
	ROB           string
	ROB_DATA      string
	FTP_EXPLORE   string
	SSH_EXPLORE   string
	SUBDOMAIN     string
	DOMAIN        string
	WORDLIST      string
	WORKERS       int
	SERVERADDR    string
)

func flags() {

	flag.StringVar(&types, "type", "", "类型")
	flag.StringVar(&GOOS, "GOOS", "", "system type")
	flag.StringVar(&GOARCH, "GOARCH", "", "架构")
	CGO_ENABLED = flag.Int("CGO_ENABLED", 0, "CGO_ENABLED")
	flag.StringVar(&FLOW, "FLOW", "", "使用流量相关的功能参数为:yes,不适用不添加即可")
	flag.StringVar(&WRITEPCAP, "WRITEPCAP", "", "把流量包写入到pcap文件中,如果使用参数值为:yes,不使用可以不添加或者no")
	flag.StringVar(&READPCAP, "READPCAP", "", "把流量包从pcap文件中读出来,如果使用参数值为:yes,不使用可以不添加或者no")
	flag.StringVar(&FINDDEVICE, "FINDDEVICE", "", "发现本机的网络设备,如果使用参数值为:yes,不使用可以不添加或者no")
	flag.StringVar(&GETFLOW, "GETFLOW", "", "获取指定网络设备的流量包,如果使用参数值为:yes,不使用可以不添加或者no")
	flag.StringVar(&PARSEFLOW, "PARSEFLOW", "", "解析流量包,如果使用参数值为:yes,不使用可以不添加或者no")
	flag.StringVar(&FLOWFILTER, "FLOWFILTER", "", "是否对流量进行过滤(过滤为BPF语法),过滤为yes,不过滤为no或者不添加")
	flag.StringVar(&EXPLORE, "EXPLORE", "", "爆破功能,使用为yes,不使用为no或者不添加")
	flag.StringVar(&LOGIN_EXPLORE, "LOGIN_EXPLORE", "", "爆破,yes or no")
	flag.StringVar(&GATHERER, "GATHERER", "", "信息收集,yes or no")
	flag.StringVar(&PORT_SCAN, "PORT_SCAN", "", "syn端口扫描,使用simple scan 输入1,使用syn泛洪保护扫描输入2")
	flag.StringVar(&ROB, "ROB", "", "使用ROB功能,ROB输入yes")
	flag.StringVar(&ROB_DATA, "ROB_DATA", "", "获取路径文件或数据库中的数据,使用文件输入1,使用数据库输入2")
	flag.StringVar(&FTP_EXPLORE, "FTP_EXPLORE", "", "使用yes")
	flag.StringVar(&SSH_EXPLORE, "SSH_EXPLORE", "", "使用yes")
	flag.StringVar(&SUBDOMAIN, "SUBDOMAIN", "", "使用yes")
	flag.StringVar(&DOMAIN, "DOMAIN", "", "To guess the domain name of the subdomain")
	flag.StringVar(&WORDLIST, "WORDLIST", "", "The wordlist path to use")
	flag.IntVar(&WORKERS, "WORKERS", 100, "The amount of wirkers to use")
	flag.StringVar(&SERVERADDR, "SERVERADDR", "8.8.8.8:53", "The DNS server to use")

	flag.Parse()
	flag.Usage()
	fmt.Println("[*]使用GOLcML中的移植代码时,把主函数名从编号改为main再编译")
}

func main() {
	flags()
	if FLOW == "yes" {
		if WRITEPCAP == "yes" {
			// for i := 0; i < 10; i++ {
			// 	go quantityofflow.Writepcap()
			// }
			go quantityofflow.Writepcap()
		}
		if FINDDEVICE == "yes" {
			quantityofflow.FindDevice()
		}
		if GETFLOW == "yes" {
			quantityofflow.GetFlow(FLOWFILTER, PARSEFLOW, WRITEPCAP)
		}
		if READPCAP == "yes" {
			quantityofflow.Readpcap()
		}
	}

	if EXPLORE == "yes" {
		if LOGIN_EXPLORE == "yes" {
			explore.Login_explore()
		} else if FTP_EXPLORE == "yes" {
			explore.FtpExplore()
		} else if SSH_EXPLORE == "yes" {
			explore.SSHExplore()
		}
	}

	if GATHERER == "yes" {
		if PORT_SCAN == "1" {
			gatherer.Simple_scan()
		} else if PORT_SCAN == "2" {
			gatherer.Scan_SYN()
		} else if SUBDOMAIN == "3" {
			gatherer.CollectingSubdomainNames(WORKERS, DOMAIN, WORDLIST, SERVERADDR)
		}
	}

	if ROB == "yes" {
		if ROB_DATA == "1" {
			rob.RobFileWindows()
		}
	}
}
