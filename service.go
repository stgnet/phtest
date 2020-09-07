package pperf

import (
	"fmt"
	"log"
	"net"
	"time"
)

func service(c net.Conn) Results {
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
	log.Printf("total=%v elms=%v bps=%v mpbs=%.2f", r.total, r.elms, r.bps, mbps(r.bps))
	return Results{
		Address:             remoteAddress,
		TotalBytes:          r.total,
		ElapsedMilliseconds: r.elms,
		BytesPerSecond:      r.bps,
		MbitsPerSecond:      mbps(r.bps),
	}
}
