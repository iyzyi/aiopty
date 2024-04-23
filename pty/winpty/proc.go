//go:build windows
// +build windows

package winpty

import (
	"fmt"
	"path/filepath"
	"syscall"
)

var (
	winptyDllName   = "winpty.dll"
	winptyAgentName = "winpty-agent.exe"
)

// include by winpty_constants.h
const _WINPTY_SPAWN_FLAG_AUTO_SHUTDOWN = 1

// include by winpty.h
var (
	modWinPty *syscall.LazyDLL

	// Error handling.
	winpty_error_code *syscall.LazyProc
	winpty_error_msg  *syscall.LazyProc
	winpty_error_free *syscall.LazyProc

	// Configuration of a new agent.
	winpty_config_new               *syscall.LazyProc
	winpty_config_free              *syscall.LazyProc
	winpty_config_set_initial_size  *syscall.LazyProc
	winpty_config_set_mouse_mode    *syscall.LazyProc
	winpty_config_set_agent_timeout *syscall.LazyProc

	// Start the agent.
	winpty_open          *syscall.LazyProc
	winpty_agent_process *syscall.LazyProc

	// I/O pipes.
	winpty_conin_name  *syscall.LazyProc
	winpty_conout_name *syscall.LazyProc
	winpty_conerr_name *syscall.LazyProc

	// winpty agent RPC call: process creation.
	winpty_spawn_config_new  *syscall.LazyProc
	winpty_spawn_config_free *syscall.LazyProc
	winpty_spawn             *syscall.LazyProc

	// winpty agent RPC calls: everything else
	winpty_set_size                 *syscall.LazyProc
	winpty_get_console_process_list *syscall.LazyProc
	winpty_free                     *syscall.LazyProc
)

func loadWinPty() (err error) {
	if modWinPty != nil {
		return
	}

	// If the golang version is at least 1.16, winpty library files in ./bin/arch will be automatically
	// embedded and released; Otherwise, they need to be manually placed in the directory where the
	// current program is located. When placing them manually, ensure that the architecture matches.
	dir, err := extractWinPty()
	if err != nil {
		return
	}

	modWinPty = syscall.NewLazyDLL(filepath.Join(dir, winptyDllName))

	// Error handling.
	winpty_error_code = modWinPty.NewProc("winpty_error_code")
	winpty_error_msg = modWinPty.NewProc("winpty_error_msg")
	winpty_error_free = modWinPty.NewProc("winpty_error_free")

	// Configuration of a new agent.
	winpty_config_new = modWinPty.NewProc("winpty_config_new")
	winpty_config_free = modWinPty.NewProc("winpty_config_free")
	winpty_config_set_initial_size = modWinPty.NewProc("winpty_config_set_initial_size")
	winpty_config_set_mouse_mode = modWinPty.NewProc("winpty_config_set_mouse_mode")
	winpty_config_set_agent_timeout = modWinPty.NewProc("winpty_config_set_agent_timeout")

	// Start the agent.
	winpty_open = modWinPty.NewProc("winpty_open")
	winpty_agent_process = modWinPty.NewProc("winpty_agent_process")

	// I/O pipes.
	winpty_conin_name = modWinPty.NewProc("winpty_conin_name")
	winpty_conout_name = modWinPty.NewProc("winpty_conout_name")
	winpty_conerr_name = modWinPty.NewProc("winpty_conerr_name")

	// winpty agent RPC call: process creation.
	winpty_spawn_config_new = modWinPty.NewProc("winpty_spawn_config_new")
	winpty_spawn_config_free = modWinPty.NewProc("winpty_spawn_config_free")
	winpty_spawn = modWinPty.NewProc("winpty_spawn")

	// winpty agent RPC calls: everything else
	winpty_set_size = modWinPty.NewProc("winpty_set_size")
	winpty_get_console_process_list = modWinPty.NewProc("winpty_get_console_process_list")
	winpty_free = modWinPty.NewProc("winpty_free")

	// check if lib available
	if !isLibAvailable() {
		return fmt.Errorf("load winpty failed")
	}
	return
}

func isLibAvailable() bool {
	return winpty_error_code != nil && winpty_error_code.Find() == nil
}
