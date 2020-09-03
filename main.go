package main

import (
	"flag"
	"net"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	test := flag.String("test", "stg.net", "test with this server")
	server := flag.Bool("server", false, "Start server")
	port := flag.Int("port", 5048, "Port to listen/send to")

	flag.Parse()

	if *server {
		l, lErr := net.Listen("tcp", ":"+strconv.Itoa(*port))
		if lErr != nil {
			panic(lErr)
		}

		for {
			c, cErr := l.Accept()
			if cErr != nil {
				log.Error(cErr)
				continue
			}
			go service(c)
			log.Infof("Closed connection")
		}
	}

	c, dErr := net.Dial("tcp", *test+":"+strconv.Itoa(*port))
	if dErr != nil {
		panic(dErr)
	}

	c.SetDeadline(time.Now().Add(10 * time.Second))
	count := 0
	for count < 100000 {
		hdr := header{
			Count: uint32(count),
		}
		count += 1
		sErr := send(c, hdr, randbytes(BLOCKSIZE-headerSize))
		if sErr != nil {
			panic(sErr)
		}
	}
}
