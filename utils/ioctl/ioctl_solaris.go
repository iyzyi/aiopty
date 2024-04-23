package ioctl

import (
	"syscall"
	"unsafe"
)

//go:cgo_import_dynamic libc_ioctl ioctl "libc.so"
//go:linkname procIoctl libc_ioctl
var procIoctl uintptr

// sysvicall6 is implemented in asm_solaris_amd64.s which has been copied from sys/unix.
func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

func Ioctl(fd, req, args uintptr) (err error) {
	_, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procIoctl)), 3, fd, req, args, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
