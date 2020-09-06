package main

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func service(c net.Conn) {
	defer c.Close()
	c.SetDeadline(time.Now().Add(10 * time.Second))
	log.Infof("Connection from %s", c.RemoteAddr().String())

	var r received
	quit := make(chan bool)
	go sender(quit, c, &r)
	for r.secs < 5 {
		rErr := receive(c, &r)
		if rErr != nil {
			log.Infof("Closing on receive error %v", rErr)
			break
		}
		switch r.hdr.Command {
		case CMD_Test:
			// nothing to do
		case CMD_IP:
			log.Infof("Remote IP: %s", string(r.data))
		case CMD_Err:
			log.Errorf("Remote ERROR: %s", string(r.data))
		case CMD_End:
			log.Infof("Received end")
		default:
			log.Errorf("Unknown command received: %v", r.hdr.Command)
		}
		if r.hdr.Command == CMD_End || r.hdr.Command == CMD_Err {
			break
		}
	}
	close(quit)
	send(c, header{Command: CMD_End}, []byte{})
	log.Infof("total=%v elms=%v bps=%v mpbs=%s", r.total, r.elms, r.bps, mbps(r.bps))
}
