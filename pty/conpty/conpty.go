package conpty

import (
	"github.com/iyzyi/aiopty/pty/common"
	"os"
	"unsafe"
)

type ConPty struct {
	opt           *common.Options
	pseudoConsole unsafe.Pointer // *syscall.Handle on windows
	pipeIn        *os.File
	pipeOut       *os.File
	process       *os.Process
	isClosed      bool
	exit          chan struct{}
}

// Open create a ConPty using path as the command to run.
func Open(path string) (*ConPty, error) {
	return openWithOptions(&common.Options{Path: path})
}

// OpenWithOptions create a ConPty with Options.
func OpenWithOptions(opt *common.Options) (*ConPty, error) {
	return openWithOptions(opt)
}

// SetSize is used to set the ConPty windows size.
func (p *ConPty) SetSize(size *common.WinSize) (err error) {
	return p.setSize(size)
}

// Close ConPty.
func (p *ConPty) Close() (err error) {
	return p.close()
}

// Read from ConPty.
func (p *ConPty) Read(b []byte) (n int, err error) {
	return p.read(b)
}

// Write to ConPty.
func (p *ConPty) Write(b []byte) (n int, err error) {
	return p.write(b)
}

// IsSupported indicates whether ConPty is supported in the current environment.
func IsSupported() bool {
	return isSupported()
}
