package native

import (
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"unsafe"
)

// See https://github.com/DragonFlyBSD/DragonFlyBSD/blob/master/lib/libc/stdlib/ptsname.c
// See https://github.com/DragonFlyBSD/DragonFlyBSD/blob/master/lib/libc/gen/fdevname.c#L60

func Grantpt(ptm *os.File) error {
	return nil
}

func Unlockpt(ptm *os.File) error {
	return nil
}

func Ptsname(ptm *os.File) (string, error) {
	name := make([]byte, ptsnameLen)
	arg := fiodnameArgs{
		name:    uintptr(unsafe.Pointer(&name[0])),
		len:     ptsnameLen,
		padding: 0,
	}
	err := ioctl.Ioctl(ptm.Fd(), FIODNAME, uintptr(unsafe.Pointer(&arg)))
	if err != nil {
		return "", err
	}

	ptmname := "/dev/" + ByteSliceToString(name)
	ptsname := []rune(ptmname)
	ptsname[len("/dev/")+2] = 's'
	return string(ptsname), nil
}

const FIODNAME = 0x80106678

type fiodnameArgs struct {
	name    uintptr
	len     uint32
	padding uint32 // memory alignment
}
