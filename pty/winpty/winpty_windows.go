package winpty

import (
	"fmt"
	"github.com/iyzyi/aiopty/pty/common"
	"github.com/iyzyi/aiopty/utils/log"
	"strings"
	"syscall"
	"unsafe"
)

func openWithOptions(opt *common.Options) (p *WinPty, err error) {
	err = common.InitOptions(opt)
	if err != nil {
		return
	}
	p = &WinPty{opt: opt}

	err = loadWinPty()
	if err != nil {
		return
	}

	agentConfig, err := newAgentConfig(0, p.opt.Size)
	if err != nil {
		return
	}

	p.pty, err = startAgent(agentConfig)
	if err != nil {
		return
	}

	p.conin, p.conout, err = getPipe(p.pty)
	if err != nil {
		return
	}

	cmdline := strings.Join(p.opt.Args, " ")
	spawnConfig, err := newSpawnConfig(_WINPTY_SPAWN_FLAG_AUTO_SHUTDOWN, p.opt.Path, cmdline, p.opt.Dir, p.opt.Env)
	if err != nil {
		return
	}

	p.process, err = spawnProcess(p.pty, spawnConfig)
	if err != nil {
		return
	}

	log.Debug("Start WinPty")
	return
}

func (p *WinPty) setSize(size *common.WinSize) (err error) {
	var errPtr uintptr
	defer winpty_error_free.Call(errPtr)
	res, _, _ := winpty_set_size.Call(p.pty, uintptr(size.Cols), uintptr(size.Rows), uintptr(unsafe.Pointer(&errPtr)))
	if res == 0 {
		return fmt.Errorf("failed to setsize: %v", getErrorMsg(errPtr))
	}
	return
}

func (p *WinPty) close() (err error) {
	if p.isClosed {
		return
	}

	if isLibAvailable() {
		winpty_free.Call(p.pty)
	}

	syscall.CloseHandle(syscall.Handle(p.process))

	p.conin.Close()
	p.conout.Close()

	p.isClosed = true
	return
}

func (p *WinPty) read(b []byte) (n int, err error) {
	return p.conout.Read(b)
}

func (p *WinPty) write(b []byte) (n int, err error) {
	return p.conin.Write(b)
}
