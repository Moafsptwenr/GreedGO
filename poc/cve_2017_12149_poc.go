package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func Cve_2017_12149_poc() {
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Println("[*] 输入测试的url: ")
	url, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	var client http.Client
	target := string(url)
	url1 := target + "/invoker/readonly"
	fmt.Println(target)
	fmt.Println(url1)
	request1, err := http.NewRequest("POST", url1, nil)
	request1.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:63.0) Gecko/20100101 Firefox/63.0")
	request1.Header.Add("Accept", "*/*")
	request1.Header.Add("Content-Type", "application/json")
	request1.Header.Add("X-Requested-With", "XMLHttpRequest")
	request1.Header.Add("Connection", "close")
	request1.Header.Add("Cache-Control", "no-cache")
	if err != nil {
		panic(err)
	}
	resp1, err := client.Do(request1)
	if err != nil {
		fmt.Printf("[-] %s没有cve-2017-12149\n", target)
	}
	if resp1.Status == "500 Internal Server Error" {
		fmt.Printf("[+] %s有cve-2017-12149\n", target)
	} else {
		fmt.Printf("[-] %s没有cve-2017-12149\n", target)
	}
}
