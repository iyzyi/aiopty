package native

import (
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"syscall"
	"unsafe"
)

// See https://github.com/freebsd/freebsd-src/blob/master/lib/libc/gen/fdevname.c

func Openpt(flags int) (ptm *os.File, err error) {
	r1, _, e1 := syscall.Syscall(syscall.SYS_POSIX_OPENPT, uintptr(flags), 0, 0)
	if e1 != 0 {
		err = e1
		return
	}
	ptm = os.NewFile(uintptr(r1), "ptm")
	return
}

func Unlockpt(ptm *os.File) error {
	return nil
}

func Grantpt(ptm *os.File) error {
	return nil
}

func Ptsname(ptm *os.File) (string, error) {
	name := make([]byte, ptsnameLen)
	arg := fiodgnameArg{
		len:     ptsnameLen,
		padding: 0,
		buf:     uintptr(unsafe.Pointer(&name[0])),
	}
	err := ioctl.Ioctl(ptm.Fd(), FIODGNAME, uintptr(unsafe.Pointer(&arg)))
	if err != nil {
		return "", err
	}
	return "/dev/" + ByteSliceToString(name), nil
}
