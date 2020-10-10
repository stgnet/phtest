package pperf

// const header size because unsafe.Sizeof() has different padding per platform
const headerSize = 32

type header struct {
	Magic   uint16 // 2 magic bytes
	Size    uint16 // 2 bytes in this message (including data beyond struct)
	Count   uint32 // 4 message counter (0=first message)
	Total   uint64 // 8 total bytes received (from other side)
	Elapsed uint64 // 8 ms elapsed since first message received
	Crc     uint16 // 2 crc of data bytes beyond header
	Command uint8  // 1 command code or indicatation of data contents
	Offset  uint8  // 1 offset from header to start of data portion
}

const (
	MAGIC     = uint16(0x5005)
	CMD_Test  = 0 // regular test data
	CMD_IP    = 1 // IP address:port string reported in data
	CMD_Err   = 2 // error occurred, message follows
	CMD_End   = 3 // close connection normally
	BLOCKSIZE = 4096
)
