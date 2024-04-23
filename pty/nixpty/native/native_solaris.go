package native

import (
	"fmt"
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"unsafe"
)

// See https://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libc/port/gen/pt.c?r=7d8deab2

func Grantpt(ptm *os.File) error {
	var ptown ptOwn
	ptown.ptoRuid = int32(os.Getuid())

	users, err := user.LookupGroup(DEFAULT_TTY_GROUP)
	if err == nil {
		var gid int
		gid, err = strconv.Atoi(users.Gid)
		if err != nil {
			return fmt.Errorf("failed to convert gid: %v", err)
		}
		ptown.ptoRgid = int32(gid)
	} else {
		ptown.ptoRgid = int32(os.Getgid())
	}

	arg := strioctl{
		icCmd:     OWNERPT,
		icTimeout: 0,
		icLen:     int32(unsafe.Sizeof(ptown)),
		icDp:      uintptr(unsafe.Pointer(&ptown)),
	}
	return ioctl.Ioctl(ptm.Fd(), I_STR, uintptr(unsafe.Pointer(&arg)))
}

func Unlockpt(ptm *os.File) error {
	arg := strioctl{
		icCmd:     UNLKPT,
		icTimeout: 0,
		icLen:     0,
		icDp:      0,
	}
	return ioctl.Ioctl(ptm.Fd(), I_STR, uintptr(unsafe.Pointer(&arg)))
}

func Ptsname(ptm *os.File) (ptsname string, err error) {
	arg := strioctl{
		icCmd:     ISPTM,
		icTimeout: 0,
		icLen:     0,
		icDp:      0,
	}
	err = ioctl.Ioctl(ptm.Fd(), I_STR, uintptr(unsafe.Pointer(&arg)))
	if err != nil {
		return
	}

	// See https://www.cnblogs.com/zongzi10010/p/11945545.html
	// See https://blog.csdn.net/zhoulaowu/article/details/14224429
	var stat syscall.Stat_t
	syscall.Fstat(int(ptm.Fd()), &stat)
	ptsname = PTSNAME + strconv.FormatUint(minor(stat.Rdev), 10)
	return
}

// See https://src.illumos.org/source/xref/illumos-gate/usr/src/uts/common/sys/stropts.h?r=b4203d75#288
type strioctl struct {
	icCmd     int32   // command
	icTimeout int32   // timeout value
	icLen     int32   // length of data
	icDp      uintptr // pointer to data
}

type ptOwn struct {
	ptoRuid int32
	ptoRgid int32
}

const I_STR = (('S' << 8) | 010)
const OWNERPT = (('P' << 8) | 5) // set owner/group for subsidiary
const UNLKPT = (('P' << 8) | 2)  // unlock manager/subsidiary pair
const ISPTM = (('P' << 8) | 1)   // query for manager

const DEFAULT_TTY_GROUP = "tty"
const PTSNAME = "/dev/pts/"

// See https://src.illumos.org/source/xref/illumos-gate/usr/src/contrib/ast/src/lib/libast/features/fs?r=b30d1939#123
func minor(x uint64) uint64 {
	return x & 0377
}
