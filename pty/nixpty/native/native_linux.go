package native

import (
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// See https://codebrowser.dev/glibc/glibc/sysdeps/unix/sysv/linux/ptsname.c.html#54

func Grantpt(ptm *os.File) error {
	return nil
}

func Unlockpt(ptm *os.File) error {
	var arg int32
	// Set (if *argp is nonzero) or remove (if *argp is zero) the
	// lock on the pseudoterminal slave device.
	// See https://man7.org/linux/man-pages/man2/ioctl_tty.2.html
	return ioctl.Ioctl(ptm.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&arg)))
}

func Ptsname(ptm *os.File) (string, error) {
	var arg uint32
	err := ioctl.Ioctl(ptm.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&arg)))
	if err != nil {
		return "", err
	}
	name := "/dev/pts/" + strconv.FormatUint(uint64(arg), 10)
	return name, nil
}
