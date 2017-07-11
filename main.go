package main

import (
	"fmt"
	"gitlab.identitylabs.mx/jescarri/bravia-remote/pkg/client"
)

func main() {
	remote := &bravia.Remote{}
	remote = remote.NewRemote("192.168.1.9", "0000")
	err := remote.SendCode("AAAAAQAAAAEAAAASAw==")
	fmt.Println(err)
}
