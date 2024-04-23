package native

import (
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"syscall"
	"unsafe"
)

// See https://www.unix.com/man-page/netbsd/4/ptm/

func Grantpt(ptm *os.File) error {
	return ioctl.Ioctl(ptm.Fd(), syscall.TIOCGRANTPT, 0)
}

func Unlockpt(ptm *os.File) error {
	return nil
}

func Ptsname(ptm *os.File) (string, error) {
	var arg ptmget
	err := ioctl.Ioctl(ptm.Fd(), TIOCPTSNAME, uintptr(unsafe.Pointer(&arg)))
	if err != nil {
		return "", err
	}
	return ByteSliceToString(arg.Sn[:]), nil
}

// See https://github.com/golang/go/issues/66871
const TIOCPTSNAME = 0x40287448

type ptmget struct {
	Cfd int32
	Sfd int32
	Cn  [16]byte
	Sn  [16]byte
}
