package pty

import (
	"fmt"
	"github.com/iyzyi/aiopty/pty/common"
	"github.com/iyzyi/aiopty/pty/conpty"
	"github.com/iyzyi/aiopty/pty/nixpty"
	"github.com/iyzyi/aiopty/pty/winpty"
	"github.com/iyzyi/aiopty/utils/log"
	"runtime"
)

// Options contains the necessary information to run Pty
type Options struct {
	// Path is the path of the command to run.
	// This is the only field that must be set to a non-zero value.
	// If Path is a file name and Dir is a zero value, search for the executable
	// file named Path in the directories specified by the PATH environment variable.
	// If Path is relative, it is evaluated relative to Dir.
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, we will use {Path} as Args.
	Args []string

	// Dir specifies the working directory of the command.
	Dir string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is empty or nil, the new process uses the current process's environment.
	Env []string

	// Size is used to set the initial pty window size.
	Size *WinSize

	// Type is used to determine which type of PTY (nixpty, conpty, or wintpy) to use.
	// By default, it is set to AUTO for automatic selection.
	Type PtyType
}

type WinSize struct {
	Cols uint16
	Rows uint16
}

type PtyType string

var (
	AUTO   PtyType = ""
	NIXPTY PtyType = "nixpty"
	CONPTY PtyType = "conpty"
	WINPTY PtyType = "winpty"
)

type Pty struct {
	opt *Options
	pty PtyApp
}

type PtyApp interface {
	SetSize(size *common.WinSize) (err error)
	Close() (err error)
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

// Open create a pty using path as the command to run.
func Open(path string) (*Pty, error) {
	return OpenWithOptions(&Options{Path: path})
}

// OpenWithOptions create a pty with Options.
func OpenWithOptions(opt *Options) (p *Pty, err error) {
	p = &Pty{opt: opt}

	if opt.Type == AUTO {
		os := runtime.GOOS
		if os != "windows" {
			opt.Type = NIXPTY
		} else {
			if conpty.IsSupported() {
				opt.Type = CONPTY
			} else {
				opt.Type = WINPTY
			}
		}
	}

	var size *common.WinSize
	if opt.Size == nil {
		size = nil
	} else {
		size = &common.WinSize{
			Cols: opt.Size.Cols,
			Rows: opt.Size.Rows,
		}
	}

	_opt := &common.Options{
		Path: opt.Path,
		Args: opt.Args,
		Dir:  opt.Dir,
		Env:  opt.Env,
		Size: size,
	}

	err = common.InitOptions(_opt)
	if err != nil {
		return
	}

	opt.Path = _opt.Path
	opt.Args = _opt.Args
	opt.Dir = _opt.Dir
	opt.Env = _opt.Env
	opt.Size = &WinSize{
		Cols: _opt.Size.Cols,
		Rows: _opt.Size.Rows,
	}

	switch opt.Type {
	case NIXPTY:
		p.pty, err = nixpty.OpenWithOptions(_opt)
	case CONPTY:
		p.pty, err = conpty.OpenWithOptions(_opt)
	case WINPTY:
		p.pty, err = winpty.OpenWithOptions(_opt)
	default:
		return nil, fmt.Errorf("error pty type: %v", opt.Type)
	}

	if err != nil {
		return nil, err
	}

	log.Debug("Path: %v", opt.Path)
	log.Debug("Args: %v", opt.Args)
	log.Debug("Dir: %v", opt.Dir)
	log.Debug("Env: %v", opt.Env)
	log.Debug("Size: %v", opt.Size)
	log.Debug("Type: %v", opt.Type)

	return p, err
}

// SetSize is used to set the pty windows size.
func (p *Pty) SetSize(size *WinSize) error {
	return p.pty.SetSize(&common.WinSize{
		Cols: size.Cols,
		Rows: size.Rows,
	})
}

// Close Pty.
func (p *Pty) Close() error {
	return p.pty.Close()
}

// Read from Pty.
func (p *Pty) Read(b []byte) (n int, err error) {
	return p.pty.Read(b)
}

// Write to Pty.
func (p *Pty) Write(b []byte) (n int, err error) {
	return p.pty.Write(b)
}
