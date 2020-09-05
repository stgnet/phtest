package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func send(c net.Conn, h header, d []byte) error {
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))
	h.Magic = MAGIC
	h.Size = uint16(headerSize + len(d))
	h.Crc = crc16(d)
	buf := new(bytes.Buffer)
	hwErr := binary.Write(buf, binary.BigEndian, &h)
	if hwErr != nil {
		return hwErr
	}
	for buf.Len() < headerSize {
		buf.WriteByte(0)
	}
	_, dwErr := buf.Write(d)
	if dwErr != nil {
		return dwErr
	}
	sent, cErr := c.Write(buf.Bytes())
	if cErr != nil {
		return cErr
	}
	if sent != int(h.Size) {
		return errors.New("Incorrect size sent")
	}
	return nil
}

func sendErr(c net.Conn, err error) {
	sErr := send(c, header{Command: CMD_Err}, []byte(err.Error()))
	if sErr != nil {
		// connection will next close, so just log it
		log.Error(sErr)
	}
}
