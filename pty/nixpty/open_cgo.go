//go:build (darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris) && cgo
// +build darwin dragonfly freebsd linux netbsd openbsd solaris
// +build cgo

package nixpty

/*
#define _XOPEN_SOURCE 600	// X/Open 6, incorporating POSIX 2004
#include <fcntl.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/ioctl.h>
*/
import "C"

import (
	"github.com/iyzyi/aiopty/utils/log"
	"os"
	"syscall"
)

// Open returns a control pty(ptm) and the linked process tty(pts).
func open() (ptm *os.File, pts *os.File, err error) {
	log.Debug("Supported by CGO")

	ptmFd, err := C.posix_openpt(syscall.O_RDWR)
	if ptmFd < 0 {
		return
	}

	res, err := C.grantpt(ptmFd)
	if res < 0 {
		C.close(ptmFd)
		return
	}

	res, err = C.unlockpt(ptmFd)
	if res < 0 {
		C.close(ptmFd)
		return
	}

	ptsname := C.GoString(C.ptsname(ptmFd))
	ptm = os.NewFile(uintptr(ptmFd), "ptm")
	pts, err = os.OpenFile(ptsname, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0)
	return
}
