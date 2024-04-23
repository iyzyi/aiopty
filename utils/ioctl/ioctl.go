//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd
// +build darwin dragonfly freebsd linux netbsd openbsd

package ioctl

import (
	"syscall"
)

func Ioctl(fd, req, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall(syscall.SYS_IOCTL, fd, req, arg)
	if e1 != 0 {
		err = e1
	}
	return
}
