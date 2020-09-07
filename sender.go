package pperf

import (
	"net"
)

func sender(quit chan bool, c net.Conn, r *received) error {
	// for first sync (count=0) packet, send along the IP address
	addr, _, _ := net.SplitHostPort(c.RemoteAddr().String())
	iErr := send(c, header{Command: CMD_IP}, []byte(addr))
	if iErr != nil {
		return iErr
	}
	count := uint32(1)
	for {
		sErr := send(c,
			header{
				Count:   count,
				Total:   r.recv.total,
				Elapsed: r.recv.elms,
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
