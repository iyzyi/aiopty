//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package term

import (
	"github.com/iyzyi/aiopty/term/export"
	"github.com/iyzyi/aiopty/utils/ioctl"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

type fields struct {
	termios *export.Termios
}

func (t *Term) isTerminal() bool {
	_, err := ioctlGetTermios(t.stdin.Fd())
	return err == nil
}

func (t *Term) wrapStdInOut() (err error) {
	termios, err := ioctlGetTermios(t.stdin.Fd())
	if err != nil {
		return
	}
	t.termios = termios

	// See https://github.com/golang/term/blob/5b15d269ba1f54e8da86c8aa5574253aea0c2198/term_unix.go#L22
	// See https://github.com/freebsd/freebsd-src/blob/1bd4f769caf8ffda35477e3c0b2c92348cf2fd5d/lib/libc/gen/termios.c#L163
	raw := *termios
	raw.Iflag &^= export.IGNBRK | export.BRKINT | export.PARMRK | export.ISTRIP | export.INLCR | export.IGNCR | export.ICRNL | export.IXON
	raw.Oflag &^= export.OPOST
	raw.Lflag &^= export.ECHO | export.ECHONL | export.ICANON | export.ISIG | export.IEXTEN
	raw.Cflag &^= export.CSIZE | export.PARENB
	raw.Cflag |= export.CS8
	raw.Cc[export.VMIN] = 1
	raw.Cc[export.VTIME] = 0

	err = ioctlSetTermios(t.stdin.Fd(), &raw)
	if err != nil {
		return
	}

	t.wrapStdin = t.stdin
	t.wrapStdout = t.stdout
	return
}

func (t *Term) restore() (err error) {
	if t.termios != nil {
		err = ioctlSetTermios(t.stdin.Fd(), t.termios)
		if err != nil {
			return
		}
	}
	return
}

func (t *Term) captureSizeChangeEvent(onSizeChange func(cols, rows uint16)) (onExit func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			cols, rows, err := t.getSize()
			if err != nil {
				return
			}
			onSizeChange(cols, rows)
		}
	}()
	ch <- syscall.SIGWINCH // initial winsize
	return func() { signal.Stop(ch); close(ch) }
}

func (t *Term) getSize() (cols, rows uint16, err error) {
	s, err := ioctlGetWinSize(t.stdin.Fd())
	if err != nil {
		return 0, 0, err
	}
	if s.Col == 0 {
		s.Col += 1
	}
	if s.Row == 0 {
		s.Row += 1
	}
	return s.Col, s.Row, err
}

func ioctlGetTermios(fd uintptr) (t *export.Termios, err error) {
	t = &export.Termios{}
	err = ioctl.Ioctl(fd, reqGetTermios, uintptr(unsafe.Pointer(t)))
	return
}

func ioctlSetTermios(fd uintptr, t *export.Termios) (err error) {
	return ioctl.Ioctl(fd, reqSetTermios, uintptr(unsafe.Pointer(t)))
}

func ioctlGetWinSize(fd uintptr) (s *Winsize, err error) {
	s = &Winsize{}
	err = ioctl.Ioctl(fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(s)))
	return
}

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}
