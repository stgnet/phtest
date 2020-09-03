package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func logSendErr(c net.Conn, err error) error {
	log.Error(err)
	sendErr(c, err)
	return err
}

type received struct {
	hdr   header
	data  []byte
	total int64
	start time.Time
	elms  int64 // elapsed milliseconds
	bps   int64 // bytes per second
}

func readex(c net.Conn, size int) ([]byte, error) {
	buf := make([]byte, size)
	got, rErr := c.Read(buf)
	if rErr != nil {
		return nil, rErr
	}
	if got < size {
		extra, xErr := readex(c, size-got)
		if xErr != nil {
			return nil, xErr
		}
		buf = append(buf[:got], extra...)
	}
	return buf, nil
}

func receive(c net.Conn, r *received) error {
	c.SetDeadline(time.Now().Add(10 * time.Second))

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
	} else {
		r.total += int64(r.hdr.Size)
		r.elms = time.Now().Sub(r.start).Milliseconds()
		if r.elms > 0 {
			r.bps = 1000 * r.total / r.elms
		}
	}
	return nil
}
