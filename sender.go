package pperf

import (
	"net"
)

func sender(quit chan bool, c net.Conn, r *received) error {
	// for first sync (count=0) packet, send along the IP address
	iErr := send(c, header{Command: CMD_IP}, []byte(c.RemoteAddr().String()))
	if iErr != nil {
		return iErr
	}
	count := uint32(1)
	for {
		sErr := send(c,
			header{
				Count:   count,
				Total:   r.total,
				Elapsed: r.elms,
			},
			randbytes(BLOCKSIZE-headerSize),
		)
		if sErr != nil {
			return sErr
		}
		count += 1

		select {
		case <-quit:
			return nil
		default:
		}
	}
}
