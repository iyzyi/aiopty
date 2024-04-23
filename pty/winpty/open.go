//go:build windows
// +build windows

package winpty

import (
	"fmt"
	"github.com/iyzyi/aiopty/pty/common"
	"os"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func newAgentConfig(flags uint32, size *common.WinSize) (agentConfig uintptr, err error) {
	var errPtr uintptr
	defer winpty_error_free.Call(errPtr)

	if runtime.GOARCH == "amd64" {
		agentConfig, _, _ = winpty_config_new.Call(uintptr(flags), uintptr(unsafe.Pointer(&errPtr)))
	} else if runtime.GOARCH == "386" {
		// In the C++ source code of winpty, the first parameter of this function is explicitly declared as UINT64,
		// so an extra 0 needs to be added for memory alignment on the 386 architecture
		agentConfig, _, _ = winpty_config_new.Call(uintptr(flags), 0, uintptr(unsafe.Pointer(&errPtr)))
	} else {
		return 0, errUnsupported
	}

	if agentConfig == 0 {
		return 0, fmt.Errorf("failed to create config: %s", getErrorMsg(errPtr))
	}

	winpty_config_set_initial_size.Call(agentConfig, uintptr(size.Cols), uintptr(size.Rows))

	return agentConfig, nil
}

func startAgent(agentConfig uintptr) (pty uintptr, err error) {
	defer winpty_config_free.Call(agentConfig)

	var errPtr uintptr
	defer winpty_error_free.Call(errPtr)

	pty, _, _ = winpty_open.Call(agentConfig, uintptr(unsafe.Pointer(&errPtr)))
	if pty == 0 {
		return 0, fmt.Errorf("failed to start agent: %s", getErrorMsg(errPtr))
	}
	return
}

func getPipe(pty uintptr) (conin, conout *os.File, err error) {
	coninName, _, _ := winpty_conin_name.Call(pty)
	coninHandle, err := syscall.CreateFile((*uint16)(unsafe.Pointer(coninName)), syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get conin handle: %s", err)
	}
	conin = os.NewFile(uintptr(coninHandle), "|0")

	conoutName, _, _ := winpty_conout_name.Call(pty)
	conoutHandle, err := syscall.CreateFile((*uint16)(unsafe.Pointer(conoutName)), syscall.GENERIC_READ, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get conout handle: %s", err)
	}
	conout = os.NewFile(uintptr(conoutHandle), "|1")

	return
}

func newSpawnConfig(flags uint32, appname, cmdline, cwd string, env []string) (spawnConfig uintptr, err error) {
	var errPtr uintptr
	defer winpty_error_free.Call(errPtr)

	_appname, err := syscall.UTF16PtrFromString(appname)
	if err != nil {
		return 0, fmt.Errorf("failed to convert appname string")
	}

	_cmdline, err := syscall.UTF16PtrFromString(cmdline)
	if err != nil {
		return 0, fmt.Errorf("failed to convert cmdline string")
	}

	_cwd, err := syscall.UTF16PtrFromString(cwd)
	if err != nil {
		return 0, fmt.Errorf("failed to convert cwd string")
	}

	_env, err := createEnvBlock(env)
	if err != nil {
		return 0, fmt.Errorf("failed to convert env string array")
	}

	if runtime.GOARCH == "amd64" {
		spawnConfig, _, _ = winpty_spawn_config_new.Call(
			uintptr(flags),
			uintptr(unsafe.Pointer(_appname)), uintptr(unsafe.Pointer(_cmdline)),
			uintptr(unsafe.Pointer(_cwd)), uintptr(unsafe.Pointer(_env)),
			uintptr(unsafe.Pointer(&errPtr)),
		)
	} else if runtime.GOARCH == "386" {
		// In the C++ source code of winpty, the first parameter of this function is explicitly declared as UINT64,
		// so an extra 0 needs to be added for memory alignment on the 386 architecture
		spawnConfig, _, _ = winpty_spawn_config_new.Call(
			uintptr(flags), uintptr(0),
			uintptr(unsafe.Pointer(_appname)), uintptr(unsafe.Pointer(_cmdline)),
			uintptr(unsafe.Pointer(_cwd)), uintptr(unsafe.Pointer(_env)),
			uintptr(unsafe.Pointer(&errPtr)),
		)
	} else {
		return 0, errUnsupported
	}

	if spawnConfig == 0 {
		return 0, fmt.Errorf("failed to create spawn config: %s", getErrorMsg(errPtr))
	}
	return
}

func spawnProcess(pty uintptr, spawnConfig uintptr) (process uintptr, err error) {
	defer winpty_spawn_config_free.Call(spawnConfig)

	var (
		processHandle    uintptr // PROCESS_INFORMATION.hProcess: A handle to the newly created process.
		threadHandle     uintptr // PROCESS_INFORMATION.hThread: A handle to the primary thread of the newly created process.
		errCreateProcess uintptr // If the agent's CreateProcess call failed, then *errCreateProcess is set to GetLastError()
		errPtr           uintptr // On failure, the function returns FALSE, and if errPtr is non-NULL, then *errPtr is set to an error object
	)
	defer winpty_error_free.Call(errPtr)

	res, _, _ := winpty_spawn.Call(
		pty, spawnConfig,
		uintptr(unsafe.Pointer(&processHandle)), uintptr(unsafe.Pointer(&threadHandle)),
		uintptr(unsafe.Pointer(&errCreateProcess)), uintptr(unsafe.Pointer(&errPtr)))

	if res == 0 {
		return 0, fmt.Errorf("failed to spawn process, err=%s, GetLastError=%v", getErrorMsg(errPtr), errCreateProcess)
	}
	return processHandle, nil
}

func getErrorMsg(ptr uintptr) string {
	msg, _, _ := winpty_error_msg.Call(ptr)
	if msg == 0 {
		return "unknown error"
	}
	return syscall.UTF16ToString((*[1024]uint16)(unsafe.Pointer(msg))[:])
}

// createEnvBlock converts an array of environment strings into
// the representation required by CreateProcess: a sequence of NUL
// terminated strings followed by a nil.
// Last bytes are two UCS-2 NULs, or four NUL bytes.
// If any string contains a NUL, it returns (nil, EINVAL).
func createEnvBlock(envv []string) (*uint16, error) {
	if len(envv) == 0 {
		return &utf16.Encode([]rune("\x00\x00"))[0], nil
	}
	length := 0
	for _, s := range envv {
		//if bytealg.IndexByteString(s, 0) != -1 {
		//	return nil, EINVAL
		//}
		length += len(s) + 1
	}
	length += 1

	b := make([]byte, length)
	i := 0
	for _, s := range envv {
		l := len(s)
		copy(b[i:i+l], []byte(s))
		copy(b[i+l:i+l+1], []byte{0})
		i = i + l + 1
	}
	copy(b[i:i+1], []byte{0})

	return &utf16.Encode([]rune(string(b)))[0], nil
}
