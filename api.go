package pperf

import (
	"log"
	"net"
	"strconv"
)

type API struct {
	Server    bool
	Target    string // ip address of server
	Port      int
	Seconds   int    // seconds to run test
	Interface string // interface name to bind to (linux only)
}

type Results struct {
	Address             string
	TotalBytes          uint64
	ElapsedMilliseconds uint64
	BytesPerSecond      uint64
	MbitsPerSecond      float64
	Err                 error
}

func Pperf(api API) Results {
	if api.Server {
		l, lErr := net.Listen("tcp", ":"+strconv.Itoa(api.Port))
		if lErr != nil {
			panic(lErr)
		}

		for {
			c, cErr := l.Accept()
			if cErr != nil {
				log.Println(cErr)
				continue
			}
			go service(c)
		}
	}

	c, dErr := net.Dial("tcp", api.Target+":"+strconv.Itoa(api.Port))
	if dErr != nil {
		return Results{Err: dErr}
	}
	return service(c)
}
