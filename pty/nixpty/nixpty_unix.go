//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package nixpty

import (
	"fmt"
	"github.com/iyzyi/aiopty/pty/common"
	"github.com/iyzyi/aiopty/utils/ioctl"
	"github.com/iyzyi/aiopty/utils/log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

func openWithOptions(opt *common.Options) (p *NixPty, err error) {
	err = common.InitOptions(opt)
	if err != nil {
		return
	}
	p = &NixPty{opt: opt}

	var tty *os.File
	p.pty, tty, err = open()
	if err != nil {
		return
	}
	defer tty.Close()

	// set block mode for pty
	// Note: The original code compiled and ran fine with go1.18, but after compiling with go1.22.0, the program
	// crashes on darwin. Debugging revealed that the issue was due to an error EAGAIN (read /dev/ptmx: resource
	// temporarily unavailable) occurring during io.Copy. Eventually, it was determined that the error was caused
	// by the non-blocking nature of the read operation on /dev/ptmx. In fact, if we ignore errors inside io.Copy,
	// the program continues to run normally. (But we can't do that.)
	err = syscall.SetNonblock(int(p.pty.Fd()), false)
	if err != nil {
		return nil, fmt.Errorf("failed to set block mode for pty: %v", err)
	}

	p.setSize(p.opt.Size)

	cmd := &exec.Cmd{
		Path:   p.opt.Path,
		Args:   p.opt.Args,
		Env:    p.opt.Env,
		Dir:    p.opt.Dir,
		Stdin:  tty,
		Stdout: tty,
		Stderr: tty,
		SysProcAttr: &syscall.SysProcAttr{
			Setsid:  true,
			Setctty: true,
		},
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	log.Debug("Start NixPty")
	return
}

func (p *NixPty) setSize(size *common.WinSize) (err error) {
	s := &struct{ Row, Col, Xpixel, Ypixel uint16 }{
		Row:    size.Rows,
		Col:    size.Cols,
		Xpixel: 0,
		Ypixel: 0,
	}
	return ioctl.Ioctl(p.pty.Fd(), uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(s)))
}

func (p *NixPty) close() (err error) {
	return p.pty.Close()
}

func (p *NixPty) read(b []byte) (n int, err error) {
	return p.pty.Read(b)
}

func (p *NixPty) write(b []byte) (n int, err error) {
	return p.pty.Write(b)
}
