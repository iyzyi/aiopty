//go:build (darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris) && !cgo
// +build darwin dragonfly freebsd linux netbsd openbsd solaris
// +build !cgo

package nixpty

import (
	"github.com/iyzyi/aiopty/pty/nixpty/native"
	"github.com/iyzyi/aiopty/utils/log"
	"os"
	"syscall"
)

// Open returns a control pty(ptm) and the linked process tty(pts).
func open() (ptm *os.File, pts *os.File, err error) {
	log.Debug("Supported by native GO")

	ptm, err = native.Openpt(syscall.O_RDWR)
	if err != nil {
		return
	}

	err = native.Grantpt(ptm)
	if err != nil {
		ptm.Close()
		return
	}

	err = native.Unlockpt(ptm)
	if err != nil {
		ptm.Close()
		return
	}

	ptsname, err := native.Ptsname(ptm)
	if err != nil {
		ptm.Close()
		return
	}

	pts, err = os.OpenFile(ptsname, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0)
	return
}
