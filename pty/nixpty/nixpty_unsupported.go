//go:build !(darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris)
// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package nixpty

import (
	"fmt"
	"github.com/iyzyi/aiopty/pty/common"
)

var errUnsupported = fmt.Errorf("unsupported os or arch")

func openWithOptions(opt *common.Options) (*NixPty, error) {
	return nil, errUnsupported
}

func (p *NixPty) setSize(size *common.WinSize) (err error) {
	return errUnsupported
}

func (p *NixPty) close() (err error) {
	return errUnsupported
}

func (p *NixPty) read(b []byte) (n int, err error) {
	return 0, errUnsupported
}

func (p *NixPty) write(b []byte) (n int, err error) {
	return 0, errUnsupported
}
