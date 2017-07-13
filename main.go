package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gitlab.identitylabs.mx/jescarri/bravia-remote/pkg/client"
)

func main() {
	remote := &bravia.Remote{}
	remote, err := remote.NewRemote("192.168.1.9", "0000")
	if err != nil {
		panic(err)
	}
	spew.Dump(remote)
	err = remote.SendCode("AAAAAQAAAAEAAAASAw==")

	fmt.Println(err)
}
