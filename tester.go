package pperf

import (
	"fmt"
	"log"
	"net"
	"time"
)

func tester(c net.Conn) Results {
	defer c.Close()
	c.SetDeadline(time.Now().Add(10 * time.Second))

	var r received
	remoteAddress := ""
	quit := make(chan bool)
	go sender(quit, c, &r)
	for r.secs < 5 {
		rErr := receive(c, &r)
		if rErr != nil {
			log.Printf("Closing on receive error %v", rErr)
			return Results{Err: rErr}
		}
		switch r.hdr.Command {
		case CMD_Test:
			// nothing to do
		case CMD_IP:
			remoteAddress = string(r.data)
		case CMD_Err:
			return Results{Err: fmt.Errorf("Server: %s", string(r.data))}
		case CMD_End:
			log.Printf("Received end")
		default:
			return Results{Err: fmt.Errorf("Unknown command: %v", r.hdr.Command)}
		}
		if r.hdr.Command == CMD_End || r.hdr.Command == CMD_Err {
			break
		}
	}
	close(quit)
	send(c, header{Command: CMD_End}, []byte{})
	return Results{
		Address: remoteAddress,
		Upload: Stats{
			TotalBytes:          r.send.total,
			ElapsedMilliseconds: r.send.elms,
			BytesPerSecond:      r.send.bps,
			MbitsPerSecond:      mbps(r.send.bps),
		},
		Download: Stats{
			TotalBytes:          r.recv.total,
			ElapsedMilliseconds: r.recv.elms,
			BytesPerSecond:      r.recv.bps,
			MbitsPerSecond:      mbps(r.recv.bps),
		},
	}
}
