package main

import (
	"gitlab.identitylabs.mx/jescarri/bravia-remote/pkg/client"
	"time"
)

func main() {
	remote := &bravia.Remote{}
	remote, err := remote.NewRemote("192.168.1.9", "0000")
	if err != nil {
		panic(err)
	}
	err = remote.Do("pause")
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	err = remote.Do("play")
	if err != nil {
		panic(err)
	}
}
