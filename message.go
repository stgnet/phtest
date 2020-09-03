package main

import (
	"unsafe"
)

type header struct {
	Magic   uint16 // magic bytes 0x5048 "PH"
	Size    uint16 // bytes in this message (including data beyond struct)
	Count   uint32 // message counter (0=first message)
	Total   uint64 // total bytes received (from other side)
	Elapsed uint64 // ms elapsed since first message received
	Crc     uint16 // crc of bytes beyond header
	Test    uint8  // non-zero is test transmit stream enable
	Command uint8  // command/state
}

const (
	MAGIC     = uint16(0x5048)
	CMD_Noop  = 0 // no response needed
	CMD_Ping  = 1 // reply with pong to return current status
	CMD_Pong  = 2 // only sent by server
	CMD_Err   = 3 // error occurred, message follows
	CMD_End   = 4 // close connection normally
	BLOCKSIZE = 1024
)

var headerSize int

func init() {
	var hdr header
	headerSize = int(unsafe.Sizeof(hdr))
}

/*
	Test sequence:

	* Client connects, sends CMD_Noop with Test:1 to start DL
	* Server sends constant stream of CMD_Noop with random data
	* Client sends CMD_Noop with Test:1 every second to defeat timeout
	* Client sends CMD_Noop with Test:0 to stop test
	* Client starts sending stream of CMD_Noop with random data
	* Client uses CMD_Ping instead once a second to get results
	* Client sends CMD_End to shut down connection

*/
