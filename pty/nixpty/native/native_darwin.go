package native

import (
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"syscall"
	"unsafe"
)

// See https://opensource.apple.com/source/Libc/Libc-1353.60.8/stdlib/grantpt.c.auto.html

func Grantpt(ptm *os.File) error {
	return ioctl.Ioctl(ptm.Fd(), syscall.TIOCPTYGRANT, 0)
}

func Unlockpt(ptm *os.File) error {
	return ioctl.Ioctl(ptm.Fd(), syscall.TIOCPTYUNLK, 0)
}

func Ptsname(ptm *os.File) (string, error) {
	bytes := make([]byte, ptsnameLen)
	err := ioctl.Ioctl(ptm.Fd(), syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&bytes[0])))
	if err != nil {
		return "", err
	}
	return ByteSliceToString(bytes), nil
}
