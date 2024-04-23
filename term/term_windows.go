package term

import (
	"fmt"
	"github.com/iyzyi/aiopty/term/color"
	"github.com/iyzyi/aiopty/utils/log"
	"syscall"
	"time"
	"unsafe"
)

type fields struct {
	inMode  uint32
	outMode uint32
	inGet   bool
	outGet  bool
	useVT   bool
}

func (t *Term) isTerminal() bool {
	_, err1 := getConsoleMode(t.stdin.Fd())
	_, err2 := getConsoleMode(t.stdout.Fd())
	return err1 == nil && err2 == nil
}

func (t *Term) wrapStdInOut() (err error) {
	err = t.wrapInput()
	if err != nil {
		return
	}
	return t.wrapOutput()
}

func (t *Term) wrapInput() (err error) {
	mode, err := getConsoleMode(t.stdin.Fd())
	if err != nil {
		return
	}
	t.inMode = mode
	t.inGet = true

	raw := mode &^ (ENABLE_ECHO_INPUT | ENABLE_PROCESSED_INPUT | ENABLE_LINE_INPUT)
	vt := raw | ENABLE_VIRTUAL_TERMINAL_INPUT

	err = setConsoleMode(t.stdin.Fd(), vt)
	t.useVT = err == nil

	if !t.useVT {
		err = setConsoleMode(t.stdin.Fd(), raw)
		if err != nil {
			return fmt.Errorf("failed to set VT or RAW input mode: %v", err)
		}
	}

	t.wrapStdin = t.stdin
	return
}

func (t *Term) wrapOutput() (err error) {
	mode, err := getConsoleMode(t.stdout.Fd())
	if err != nil {
		return err
	}
	t.outMode = mode
	t.outGet = true

	if t.useVT {
		vt := mode | ENABLE_VIRTUAL_TERMINAL_PROCESSING
		err = setConsoleMode(t.stdout.Fd(), vt)
		if err != nil {
			return fmt.Errorf("failed to set VT output mode: %v", err)
		}
		t.wrapStdout = t.stdout
		log.Debug("Using Console Virtual Terminal Sequences to handle ANSI escape sequences")
	} else {
		t.wrapStdout = color.NewColorable(t.stdout)
		log.Debug("Using Third Party Package to handle ANSI escape sequences")
	}
	return
}

func (t *Term) restore() (err error) {
	if t.inGet {
		err = setConsoleMode(t.stdin.Fd(), t.inMode)
		if err != nil {
			return
		}
	}
	if t.outGet {
		err = setConsoleMode(t.stdout.Fd(), t.outMode)
		if err != nil {
			return
		}
	}
	return
}

func (t *Term) captureSizeChangeEvent(onSizeChange func(cols, rows uint16)) (onExit func()) {
	ch := make(chan struct{}, 1)
	var prevWidth, prevHeight uint16 = 0, 0
	var curWidth, curHeight uint16 = 0, 0
	var init bool // initial winsize
	var err error

	go func() {
		for {
			select {
			case <-ch:
				break

			default:
				curWidth, curHeight, err = t.getSize()
				if err != nil {
					break
				}
				if curWidth != prevWidth || curHeight != prevHeight || !init {
					prevWidth, prevHeight = curWidth, curHeight
					if onSizeChange == nil {
						break
					}
					onSizeChange(curWidth, curHeight)
					init = true
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	return func() { close(ch) }
}

func (t *Term) getSize() (cols, rows uint16, err error) {
	info, err := getConsoleScreenBufferInfo(t.stdout.Fd())
	if err != nil {
		return
	}
	cols = info.window.right - info.window.left + 1
	if cols == 0 {
		cols += 1
	}
	rows = info.window.bottom - info.window.top + 1
	if rows == 0 {
		rows += 1
	}
	return
}

const (
	// INPUT
	ENABLE_PROCESSED_INPUT        = 0x1
	ENABLE_LINE_INPUT             = 0x2
	ENABLE_ECHO_INPUT             = 0x4
	ENABLE_VIRTUAL_TERMINAL_INPUT = 0x200
	// OUTPUT
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x4
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode             = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

func getConsoleMode(fd uintptr) (mode uint32, err error) {
	r1, _, e1 := procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(&mode)))
	if r1 == 0 {
		err = e1
	}
	return
}

func setConsoleMode(fd uintptr, mode uint32) (err error) {
	r1, _, e1 := procSetConsoleMode.Call(fd, uintptr(mode))
	if r1 == 0 {
		err = e1
	}
	return
}

func getConsoleScreenBufferInfo(fd uintptr) (*consoleScreenBufferInfo, error) {
	var info consoleScreenBufferInfo
	r1, _, e1 := procGetConsoleScreenBufferInfo.Call(fd, uintptr(unsafe.Pointer(&info)))
	if r1 == 0 {
		return nil, e1
	}
	return &info, nil
}

type coord struct {
	x uint16
	y uint16
}

type smallRect struct {
	left   uint16
	top    uint16
	right  uint16
	bottom uint16
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        uint16
	window            smallRect
	maximumWindowSize coord
}
