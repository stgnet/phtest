package pperf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

func send(c net.Conn, h header, d []byte) error {
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))
	h.Magic = MAGIC
	h.Offset = uint8(headerSize)
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

func sendErr(c net.Conn, err error) error {
	send(c, header{Command: CMD_Err}, []byte(err.Error()))
	// connection will be closed next anyway, so ignore any error sending
	return err
}
