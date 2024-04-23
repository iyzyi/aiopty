//go:build !windows
// +build !windows

package winpty

import (
	"github.com/iyzyi/aiopty/pty/common"
)

func openWithOptions(opt *common.Options) (*WinPty, error) {
	return nil, errUnsupported
}

func (p *WinPty) setSize(size *common.WinSize) (err error) {
	return errUnsupported
}

func (p *WinPty) close() (err error) {
	return errUnsupported
}

func (p *WinPty) read(b []byte) (n int, err error) {
	return 0, errUnsupported
}

func (p *WinPty) write(b []byte) (n int, err error) {
	return 0, errUnsupported
}
