What is this
------------
This is a simple Go Client for Sony Bravia Remote Control API

Is based on the following work:

1. https://github.com/alanreid/bravia

2. https://github.com/breunigs/bravia-auth-and-remote

How to Use it
-------------

First enable basic auth into your TV, to do it follow this instructions

1. Navigate to: [Settings] → [Network] → [Home Network Setup] → [IP Control]
2. Set [Authentication] to [Normal and Pre-Shared Key]
3. There should be a new menu entry [Pre-Shared Key]. Set it to a four digit code, save those numbers.

A basic progam will look like:

when you call remote.NewRemote() you need to provide the ipaddress or hostname of the tv and the PSK you set on setp 3.

````
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
	err = remote.Do("Pause")
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	err = remote.Do("Play")
	if err != nil {
		panic(err)
	}
}
````


