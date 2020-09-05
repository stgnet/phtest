package main

import (
	"net"
)

func sender(quit chan bool, c net.Conn, r *received) {
	count := uint32(0)
	for {
		sErr := send(c,
			header{
				Count:   count,
				Total:   r.total,
				Elapsed: r.elms,
				Test:    1,
			},
			randbytes(BLOCKSIZE-headerSize),
		)
		if sErr != nil {
			return
		}
		count += 1

		select {
		case <-quit:
			return
		default:
		}
	}
}
