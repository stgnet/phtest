// +build linux

package pperf

import "syscall"

func bindInterface(fd int, ifname string) error {
	return syscall.SetsockoptString(fd, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ifname)
}
