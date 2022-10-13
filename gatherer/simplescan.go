package gatherer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
)

func worker(ports, results chan int, ip string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", ip, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func Simple_scan() {
	var ports = make(chan int, 200)
	var results = make(chan int, 200)
	var open_ports []int
	fmt.Println("输入IP:")
	reader := bufio.NewReader(os.Stdin)
	ip, _, err := reader.ReadLine()
	if err != nil {
		log.Fatalln("ip 输入错误")
	}

	if err != nil {
		log.Fatalln("read ip error")
	}

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, string(ip))
	}

	go func() {
		for i := 0; i <= 65535; i++ {
			ports <- i
		}
	}()

	for i := 0; i <= 65535; i++ {
		port := <-results
		if port != 0 {
			open_ports = append(open_ports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(open_ports)

	for _, port := range open_ports {
		fmt.Println("开放的端口：", port)
	}
}
