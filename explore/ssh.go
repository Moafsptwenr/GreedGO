package explore

import (
	"fmt"

	_ "golang.org/x/crypto/ssh"
)

func SSHExplore() {
	fmt.Printf("[*]%c[1;32;40m%s%c[0m\n", 0x1B, "SSH exploring,please waiting.....", 0x1B)
	//gossh.Dial()
}
