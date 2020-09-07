package pperf

import (
	"fmt"
	"net"
	"time"
)

func tester(c net.Conn, duration int) Results {
	defer c.Close()
	c.SetDeadline(time.Now().Add(10 * time.Second))

	// server is given 0 for duration, but force an upper limit to test duration
	if duration == 0 {
		// if server closes connection first, client will receive error
		duration = 60
	}

	var r received
	remoteAddress := ""
	quit := make(chan bool)
	go sender(quit, c, &r)
	for r.secs < duration {
		rErr := receive(c, &r)
		if rErr != nil {
			return Results{Err: rErr}
		}

		if r.hdr.Command == CMD_IP {
			remoteAddress = string(r.data)
		}
		if r.hdr.Command == CMD_Err {
			return Results{Err: fmt.Errorf("Remote error: %s", string(r.data))}
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
