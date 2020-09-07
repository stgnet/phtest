package pperf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func logSendErr(c net.Conn, err error) error {
	log.Println(err)
	sendErr(c, err)
	return err
}

type received struct {
	hdr   header
	data  []byte
	start time.Time
	total uint64
	elms  uint64 // elapsed milliseconds
	bps   uint64 // bytes per second
	secs  int    // elapsed seconds
	last  int    // last elapsed seconds (tick trigger)
}

// bytes per second
func bps(total uint64, elms uint64) uint64 {
	if elms == 0 {
		return 0
	}
	return 1000 * total / elms
}

// mega-bits-per-second
func mbps(bps uint64) float64 {
	return float64(bps) * 0.000008
}

func readex(c net.Conn, size int) ([]byte, error) {
	buf := make([]byte, size)
	have, rErr := c.Read(buf)
	if rErr != nil {
		return nil, rErr
	}
	for have < size {
		need := size - have
		extra := make([]byte, need)
		add, aErr := c.Read(extra)
		if aErr != nil {
			return nil, aErr
		}
		buf = append(buf[:have], extra[:add]...)
		have += add
	}

	/*
		prior recursive version for reference:
		if have < size {
			extra, xErr := readex(c, size-have)
			if xErr != nil {
				return nil, xErr
			}
			buf = append(buf[:got], extra...)
		}
	*/
	return buf, nil
}

func receive(c net.Conn, r *received) error {
	c.SetReadDeadline(time.Now().Add(10 * time.Second))

	hbuf, rErr := readex(c, headerSize) //c.Read(hbuf)
	if rErr != nil {
		return logSendErr(c, fmt.Errorf("read header: %w", rErr))
	}
	hrErr := binary.Read(bytes.NewReader(hbuf), binary.BigEndian, &r.hdr)
	if hrErr != nil {
		return logSendErr(c, fmt.Errorf("convert header: %w: buf=%+v", hrErr, hbuf))
	}
	if r.hdr.Magic != MAGIC {
		return logSendErr(c, errors.New("received message with wrong magic code"))
	}
	if int(r.hdr.Size) < headerSize || int(r.hdr.Size) > BLOCKSIZE {
		return logSendErr(c, errors.New("received message with invalid header size"))
	}
	// log.Infof("Received %+v", r.hdr)

	dsize := int(r.hdr.Size) - headerSize
	if dsize > 0 {
		dbuf, dErr := readex(c, dsize)
		if dErr != nil {
			return logSendErr(c, fmt.Errorf("read data: %w", dErr))
		}
		if r.hdr.Crc != crc16(dbuf) {
			return logSendErr(c, errors.New("received wrong data crc"))
		}
		r.data = dbuf
	}

	// update calculations
	if r.hdr.Count == 0 {
		// reset values -- first packet doesn't count
		r.total = 0
		r.start = time.Now()
		r.last = 0
		r.secs = 0
		r.bps = 0
	} else {
		r.last = r.secs
		r.secs = int(time.Now().Sub(r.start).Seconds())

		r.total += uint64(r.hdr.Size)
		r.elms = uint64(time.Now().Sub(r.start).Milliseconds())
		r.bps = bps(r.total, r.elms)
	}
	return nil
}
