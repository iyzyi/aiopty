//go:build windows && !go1.16
// +build windows,!go1.16

package winpty

import (
	"fmt"
	"os"
	"path/filepath"
)

var errNotExistsWinPtyLib = fmt.Errorf("winpty lib files does not exist")

// For Go versions lower than 1.16, this function will not perform the extraction action;
// it will only check if the winpty library file exists in the directory. Therefore, it is
// necessary to manually place the winpty library file from ./bin/arch into the directory
// where the current program is located, according to the corresponding architecture.
func extractWinPty() (dir string, err error) {
	//execPath, err := os.Executable()
	//if err != nil {
	//	return
	//}
	//execDir := filepath.Dir(execPath)
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	winptyDllPath := filepath.Join(execDir, winptyDllName)
	winptyAgentPath := filepath.Join(execDir, winptyAgentName)

	if _, err = os.Stat(winptyDllPath); os.IsNotExist(err) {
		err = errNotExistsWinPtyLib
		return
	}

	if _, err = os.Stat(winptyAgentPath); os.IsNotExist(err) {
		err = errNotExistsWinPtyLib
		return
	}

	return execDir, err
}
