//go:build windows && go1.16
// +build windows,go1.16

package winpty

import (
	"embed"
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

//go:embed bin/*
var f embed.FS

func extractWinPty() (dir string, err error) {
	var label string
	arch := runtime.GOARCH
	if arch == "386" {
		label = "ia32"
	} else if arch == "amd64" {
		label = "x64"
	} else {
		return "", errUnsupported
	}

	execPath, err := os.Executable()
	if err != nil {
		return
	}
	execDir := filepath.Dir(execPath)

	winptyDllPath := filepath.Join(execDir, winptyDllName)
	winptyAgentPath := filepath.Join(execDir, winptyAgentName)

	var dll []byte
	if _, err = os.Stat(winptyDllPath); errors.Is(err, os.ErrNotExist) {
		dll, err = f.ReadFile(path.Join("bin", label, winptyDllName))
		if err != nil {
			return
		}
		err = os.WriteFile(winptyDllPath, dll, 0700)
		if err != nil {
			return
		}
	}

	var exe []byte
	if _, err = os.Stat(winptyAgentPath); errors.Is(err, os.ErrNotExist) {
		exe, err = f.ReadFile(path.Join("bin", label, winptyAgentName))
		if err != nil {
			return
		}
		err = os.WriteFile(winptyAgentPath, exe, 0700)
		if err != nil {
			return
		}
	}

	return execDir, err
}
