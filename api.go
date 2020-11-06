package pperf

import (
	"fmt"
	"net"
	"strconv"
	"syscall"
	"time"
)

type API struct {
	Server    bool
	Target    string // ip address of server to test against
	Port      int
	Seconds   int    // number of seconds to run test
	Interface string // interface name to bind to (linux only)
}

type Stats struct {
	TotalBytes          uint64
	ElapsedMilliseconds uint64
	BytesPerSecond      uint64
	MbitsPerSecond      float64
}

type Results struct {
	Address  string
	Download Stats
	Upload   Stats
}

func Run(api API) (*Results, error) {
	if api.Port == 0 {
		api.Port = 5048
	}
	if api.Seconds == 0 {
		api.Seconds = 15
	}
	if api.Server {
		l, lErr := net.Listen("tcp", ":"+strconv.Itoa(api.Port))
		if lErr != nil {
			panic(lErr)
		}

		for {
			c, cErr := l.Accept()
			if cErr != nil {
				fmt.Println(cErr)
				continue
			}
			go tester(c, 0)
		}
	}

	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true,
	}
	if api.Interface != "" {
		dialer.Control = func(network, address string, c syscall.RawConn) (err error) {
			err1 := c.Control(func(fd uintptr) {
				err = bindInterface(int(fd), api.Interface)
				if err != nil {
					return
				}
			})
			if err != nil {
				return err
			}
			return err1
		}
	}
	c, dErr := dialer.Dial("tcp", api.Target+":"+strconv.Itoa(api.Port))
	if dErr != nil {
		return nil, dErr
	}
	result, err := tester(c, api.Seconds)
	if err != nil {
		return nil, err
	}

	return result, nil
}
