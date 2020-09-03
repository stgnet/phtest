package main

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func service(c net.Conn) {
	defer c.Close()

	log.Infof("Connection from %v", c.RemoteAddr())

	var r received
	for {
		rErr := receive(c, &r)
		if rErr != nil {
			log.Infof("Closing on receive error %v", rErr)
			break
		}

	}
}
