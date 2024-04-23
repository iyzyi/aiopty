package conpty

import (
	"fmt"
	"github.com/iyzyi/aiopty/pty/common"
	"github.com/iyzyi/aiopty/utils/log"
	"os"
	"syscall"
	"unsafe"
)

func openWithOptions(opt *common.Options) (p *ConPty, err error) {
	if !IsSupported() {
		return nil, fmt.Errorf("unsupported windows version")
	}

	err = common.InitOptions(opt)
	if err != nil {
		return
	}

	p = &ConPty{
		opt:           opt,
		pseudoConsole: unsafe.Pointer(new(syscall.Handle)),
	}

	// create pipe
	var ptyIn, ptyOut *os.File
	ptyIn, p.pipeOut, err = os.Pipe()
	// Note: We can close the handles to the PTY-end of the pipes after createPseudoConsole because
	// the handles are dup'ed into the ConHost and will be released when the ConPty is destroyed.
	defer ptyIn.Close()
	if err != nil {
		return
	}
	p.pipeIn, ptyOut, err = os.Pipe()
	defer ptyOut.Close()
	if err != nil {
		p.pipeOut.Close()
		return
	}

	// create conpty
	err = createPseudoConsole(packWinSize(p.opt.Size), syscall.Handle(ptyIn.Fd()), syscall.Handle(ptyOut.Fd()), p.getPseudoConsole())
	if err != nil {
		return
	}

	attr := &syscall.ProcAttr{
		Dir:   p.opt.Dir,
		Env:   p.opt.Env,
		Files: make([]uintptr, 3),
		Sys:   nil,
	}

	pid, _, err := startProcessWithConPty(p.opt.Path, p.opt.Args, attr, p.getPseudoConsole())
	if err != nil {
		return
	}

	// Tests revealed that when the terminal corresponding to ConPty exits, the read & write pipes of ConPty
	// are not closed, causing both io.Copy operations to be blocked. Therefore, once we detect that the
	// subprocess launched by ConPty has exited, we will close that ConPty to terminate the io.Copy operations.
	go func() {
		p.process, err = os.FindProcess(pid)
		if err != nil {
			return
		}
		p.process.Wait()
		p.close()
	}()

	log.Debug("Start ConPty")
	return
}

func (p *ConPty) setSize(size *common.WinSize) (err error) {
	return resizePseudoConsole(*p.getPseudoConsole(), packWinSize(size))
}

func (p *ConPty) close() (err error) {
	if p.isClosed {
		return
	}

	err = closePseudoConsole(*p.getPseudoConsole())

	p.pipeIn.Close()
	p.pipeOut.Close()

	p.isClosed = true
	return
}

func (p *ConPty) read(b []byte) (n int, err error) {
	return p.pipeIn.Read(b)
}

func (p *ConPty) write(b []byte) (n int, err error) {
	return p.pipeOut.Write(b)
}

func (p *ConPty) getPseudoConsole() *syscall.Handle {
	return (*syscall.Handle)(p.pseudoConsole)
}

func packWinSize(size *common.WinSize) (s uintptr) {
	return uintptr(size.Cols) + (uintptr(size.Rows) << 16)
}

func isSupported() bool {
	return procCreatePseudoConsole.Find() == nil &&
		procResizePseudoConsole.Find() == nil &&
		procClosePseudoConsole.Find() == nil
}

// syscall functions
var (
	modKernel32             = syscall.NewLazyDLL("kernel32.dll")
	procCreatePseudoConsole = modKernel32.NewProc("CreatePseudoConsole")
	procResizePseudoConsole = modKernel32.NewProc("ResizePseudoConsole")
	procClosePseudoConsole  = modKernel32.NewProc("ClosePseudoConsole")
)

func createPseudoConsole(size uintptr, ptyIn syscall.Handle, ptyOut syscall.Handle, pseudoConsole *syscall.Handle) (err error) {
	r1, _, e1 := procCreatePseudoConsole.Call(size, uintptr(ptyIn), uintptr(ptyOut), 0, uintptr(unsafe.Pointer(pseudoConsole)))
	if r1 != 0 {
		err = e1
	}
	return
}

func resizePseudoConsole(pseudoConsole syscall.Handle, size uintptr) (err error) {
	r1, _, e1 := procResizePseudoConsole.Call(uintptr(pseudoConsole), size)
	if r1 != 0 {
		err = e1
	}
	return
}

func closePseudoConsole(pseudoConsole syscall.Handle) (err error) {
	r1, _, e1 := procClosePseudoConsole.Call(uintptr(pseudoConsole))
	if r1 == 0 {
		err = e1
	}
	return
}
