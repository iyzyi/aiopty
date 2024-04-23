package native

import (
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"unsafe"
)

// See https://github.com/openbsd/src/blob/master/lib/libc/stdlib/posix_pty.c
// See https://man.openbsd.org/OpenBSD-5.5/ptm.4

func Openpt(flags int) (ptm *os.File, err error) {
	_ptm, err := os.OpenFile("/dev/ptm", flags, 0)
	if err != nil {
		return
	}

	var arg ptmget
	err = ioctl.Ioctl(_ptm.Fd(), PTMGET, uintptr(unsafe.Pointer(&arg)))
	if err != nil {
		return
	}
	ptm = os.NewFile(uintptr(arg.Cfd), ByteSliceToString(arg.Cn[:]))
	// Note: We can directly use arg.Sfd and arg.Sn to obtain the pts, but we choose
	// not to do so in order to maintain func interface consistency.
	return
}

func Grantpt(ptm *os.File) error {
	return nil
}

func Unlockpt(ptm *os.File) error {
	return nil
}

func Ptsname(ptm *os.File) (string, error) {
	// e.g. /dev/ptyp4 -> /dev/ttyp4
	ptsname := []byte(ptm.Name())
	ptsname[len("/dev/")] = 't'
	return string(ptsname), nil
}

const PTMGET = 0x40287401

type ptmget struct {
	Cfd int32
	Sfd int32
	Cn  [16]byte
	Sn  [16]byte
}
