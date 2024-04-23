package common

import (
	"os"
	"os/exec"
	"path/filepath"
)

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
}

type WinSize struct {
	Cols uint16
	Rows uint16
}

func InitOptions(opt *Options) (err error) {
	if opt.Dir == "" {
		var _path string
		_path, err = extendPath(opt.Path)
		if err != nil {
			err = nil
		} else {
			opt.Path = _path
		}
	}

	if len(opt.Args) == 0 {
		opt.Args = []string{opt.Path}
	}

	if len(opt.Env) == 0 {
		opt.Env = os.Environ()
	}

	if opt.Size == nil {
		opt.Size = &WinSize{
			Cols: 120,
			Rows: 30,
		}
	}
	return
}

// If path is a file name, search for the executable file named path in the
// directories specified by the PATH environment variable.
func extendPath(path string) (string, error) {
	if filepath.Base(path) == path {
		lp, err := exec.LookPath(path)
		if err != nil {
			return "", err
		}
		return lp, err
	}
	return path, nil
}
