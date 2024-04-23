//go:build !windows
// +build !windows

package conpty

import (
	"fmt"
	"github.com/iyzyi/aiopty/pty/common"
)

var errUnsupported = fmt.Errorf("unsupported os or arch")

func openWithOptions(opt *common.Options) (*ConPty, error) {
	return nil, errUnsupported
}

func (pty *ConPty) setSize(size *common.WinSize) (err error) {
	return errUnsupported
}

func (pty *ConPty) close() (err error) {
	return errUnsupported
}

func (pty *ConPty) read(b []byte) (n int, err error) {
	return 0, errUnsupported
}

func (pty *ConPty) write(b []byte) (n int, err error) {
	return 0, errUnsupported
}

func isSupported() bool {
	return false
}
