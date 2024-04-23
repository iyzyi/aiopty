package nixpty

import (
	"github.com/iyzyi/aiopty/pty/common"
	"os"
)

type NixPty struct {
	opt *common.Options
	pty *os.File
}

// Open create a NixPty using path as the command to run.
func Open(path string) (*NixPty, error) {
	return openWithOptions(&common.Options{Path: path})
}

// OpenWithOptions create a NixPty with Options.
func OpenWithOptions(opt *common.Options) (*NixPty, error) {
	return openWithOptions(opt)
}

// SetSize is used to set the NixPty windows size.
func (p *NixPty) SetSize(size *common.WinSize) (err error) {
	return p.setSize(size)
}

// Close NixPty.
func (p *NixPty) Close() (err error) {
	return p.close()
}

// Read from NixPty.
func (p *NixPty) Read(b []byte) (n int, err error) {
	return p.read(b)
}

// Write to NixPty.
func (p *NixPty) Write(b []byte) (n int, err error) {
	return p.write(b)
}
