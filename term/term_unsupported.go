//go:build !(darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows)
// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package term

import "fmt"

var errUnsupported = fmt.Errorf("unsupported os or arch")

type fields struct{}

func (t *Term) isTerminal() bool {
	return false
}

func (t *Term) wrapStdInOut() (err error) {
	return errUnsupported
}

func (t *Term) restore() (err error) {
	return errUnsupported
}

func (t *Term) captureSizeChangeEvent(onSizeChange func(cols, rows uint16)) (onExit func()) {
	return nil
}
