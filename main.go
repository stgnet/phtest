package main

import (
	"flag"
	"log"
	"net"
	"strconv"
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
				log.Println(cErr)
				continue
			}
			go service(c)
		}
	}

	c, dErr := net.Dial("tcp", *test+":"+strconv.Itoa(*port))
	if dErr != nil {
		panic(dErr)
	}
	service(c)
}
