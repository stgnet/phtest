package main

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func service(c net.Conn) {
	defer c.Close()
	c.SetDeadline(time.Now().Add(10 * time.Second))
	log.Infof("Connection from %v", c.RemoteAddr())

	var r received
	quit := make(chan bool)
	go sender(quit, c, &r)
	for r.secs < 5 {
		rErr := receive(c, &r)
		if rErr != nil {
			log.Infof("Closing on receive error %v", rErr)
			break
		}

	}
	close(quit)
	log.Infof("total=%v elms=%v bps=%v mpbs=%s", r.total, r.elms, r.bps, mbps(r.bps))
}
