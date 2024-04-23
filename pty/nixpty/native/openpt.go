//go:build darwin || dragonfly || linux || netbsd || solaris
// +build darwin dragonfly linux netbsd solaris

package native

import (
	"os"
)

// See https://github.com/coreutils/gnulib/blob/master/lib/posix_openpt.c

func Openpt(flags int) (ptm *os.File, err error) {
	ptm, err = os.OpenFile("/dev/ptmx", flags, 0)
	if err != nil {
		return
	}
	return
}
