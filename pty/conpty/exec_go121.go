//go:build windows && go1.21
// +build windows,go1.21

// The core code is based on syscall.StartProcess.

// Due to exec.Cmd currently not supporting ConPty ( See https://github.com/golang/go/pull/62710 ),
// we need to manually export copies of the relevant functions and make modifications to them.

// Once we obtain the newly created pseudoConsole by calling createPseudoConsole, we need to enable
// ConPTY support for child process according to the following steps:
// 1) Use StartupInfoEx instead of StartupInfo.
// 2) Create ProcThreadAttributeList using newProcThreadAttributeList, and remember to delete this
//    structure using deleteProcThreadAttributeList at the end.
// 3) Update the pseudoConsole info into the ProcThreadAttributeList using updateProcThreadAttribute.
// 4) Update the flags with _EXTENDED_STARTUPINFO_PRESENT to enable StartupInfoEx.
// 5) CreateProcess with si and flags.

package conpty

import "syscall"

import (
	"runtime"
	"unicode/utf16"
	"unsafe"
)

type _STARTUPINFOEXW struct {
	syscall.StartupInfo
	ProcThreadAttributeList *_PROC_THREAD_ATTRIBUTE_LIST
}

type _PROC_THREAD_ATTRIBUTE_LIST struct {
	_ [1]byte
}

const (
	_PROC_THREAD_ATTRIBUTE_PARENT_PROCESS = 0x00020000
	_PROC_THREAD_ATTRIBUTE_HANDLE_LIST    = 0x00020002
)

const _EXTENDED_STARTUPINFO_PRESENT = 0x00080000

var zeroProcAttr syscall.ProcAttr
var zeroSysProcAttr syscall.SysProcAttr

// func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid int, handle uintptr, err error)
func startProcessWithConPty(argv0 string, argv []string, attr *syscall.ProcAttr, pseudoConsole *syscall.Handle) (pid int, handle uintptr, err error) {
	if len(argv0) == 0 {
		return 0, 0, syscall.EWINDOWS
	}
	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	if len(attr.Files) > 3 {
		return 0, 0, syscall.EWINDOWS
	}
	if len(attr.Files) < 3 {
		return 0, 0, syscall.EINVAL
	}

	if len(attr.Dir) != 0 {
		// StartProcess assumes that argv0 is relative to attr.Dir,
		// because it implies Chdir(attr.Dir) before executing argv0.
		// Windows CreateProcess assumes the opposite: it looks for
		// argv0 relative to the current directory, and, only once the new
		// process is started, it does Chdir(attr.Dir). We are adjusting
		// for that difference here by making argv0 absolute.
		var err error
		argv0, err = joinExeDirAndFName(attr.Dir, argv0)
		if err != nil {
			return 0, 0, err
		}
	}
	argv0p, err := syscall.UTF16PtrFromString(argv0)
	if err != nil {
		return 0, 0, err
	}

	var cmdline string
	// Windows CreateProcess takes the command line as a single string:
	// use attr.CmdLine if set, else build the command line by escaping
	// and joining each argument with spaces
	if sys.CmdLine != "" {
		cmdline = sys.CmdLine
	} else {
		cmdline = makeCmdLine(argv)
	}

	var argvp *uint16
	if len(cmdline) != 0 {
		argvp, err = syscall.UTF16PtrFromString(cmdline)
		if err != nil {
			return 0, 0, err
		}
	}

	var dirp *uint16
	if len(attr.Dir) != 0 {
		dirp, err = syscall.UTF16PtrFromString(attr.Dir)
		if err != nil {
			return 0, 0, err
		}
	}

	var maj, min, build uint32
	rtlGetNtVersionNumbers(&maj, &min, &build)
	isWin7 := maj < 6 || (maj == 6 && min <= 1)
	// NT kernel handles are divisible by 4, with the bottom 3 bits left as
	// a tag. The fully set tag correlates with the types of handles we're
	// concerned about here.  Except, the kernel will interpret some
	// special handle values, like -1, -2, and so forth, so kernelbase.dll
	// checks to see that those bottom three bits are checked, but that top
	// bit is not checked.
	isLegacyWin7ConsoleHandle := func(handle syscall.Handle) bool { return isWin7 && handle&0x10000003 == 3 }

	p, _ := syscall.GetCurrentProcess()
	parentProcess := p
	if sys.ParentProcess != 0 {
		parentProcess = sys.ParentProcess
	}
	fd := make([]syscall.Handle, len(attr.Files))
	for i := range attr.Files {
		if attr.Files[i] > 0 {
			destinationProcessHandle := parentProcess

			// On Windows 7, console handles aren't real handles, and can only be duplicated
			// into the current process, not a parent one, which amounts to the same thing.
			if parentProcess != p && isLegacyWin7ConsoleHandle(syscall.Handle(attr.Files[i])) {
				destinationProcessHandle = p
			}

			err := syscall.DuplicateHandle(p, syscall.Handle(attr.Files[i]), destinationProcessHandle, &fd[i], 0, true, syscall.DUPLICATE_SAME_ACCESS)
			if err != nil {
				return 0, 0, err
			}
			defer syscall.DuplicateHandle(parentProcess, fd[i], 0, nil, 0, false, syscall.DUPLICATE_CLOSE_SOURCE)
		}
	}
	si := new(_STARTUPINFOEXW)
	//si.ProcThreadAttributeList, err = newProcThreadAttributeList(2)
	si.ProcThreadAttributeList, err = newProcThreadAttributeList(3)
	if err != nil {
		return 0, 0, err
	}
	defer deleteProcThreadAttributeList(si.ProcThreadAttributeList)
	si.Cb = uint32(unsafe.Sizeof(*si))
	si.Flags = syscall.STARTF_USESTDHANDLES
	if sys.HideWindow {
		si.Flags |= syscall.STARTF_USESHOWWINDOW
		si.ShowWindow = syscall.SW_HIDE
	}
	if sys.ParentProcess != 0 {
		err = updateProcThreadAttribute(si.ProcThreadAttributeList, 0, _PROC_THREAD_ATTRIBUTE_PARENT_PROCESS, unsafe.Pointer(&sys.ParentProcess), unsafe.Sizeof(sys.ParentProcess), nil, nil)
		if err != nil {
			return 0, 0, err
		}
	}
	si.StdInput = fd[0]
	si.StdOutput = fd[1]
	si.StdErr = fd[2]

	fd = append(fd, sys.AdditionalInheritedHandles...)

	// On Windows 7, console handles aren't real handles, so don't pass them
	// through to PROC_THREAD_ATTRIBUTE_HANDLE_LIST.
	for i := range fd {
		if isLegacyWin7ConsoleHandle(fd[i]) {
			fd[i] = 0
		}
	}

	// The presence of a NULL handle in the list is enough to cause PROC_THREAD_ATTRIBUTE_HANDLE_LIST
	// to treat the entire list as empty, so remove NULL handles.
	j := 0
	for i := range fd {
		if fd[i] != 0 {
			fd[j] = fd[i]
			j++
		}
	}
	fd = fd[:j]

	willInheritHandles := len(fd) > 0 && !sys.NoInheritHandles

	// Do not accidentally inherit more than these handles.
	if willInheritHandles {
		err = updateProcThreadAttribute(si.ProcThreadAttributeList, 0, _PROC_THREAD_ATTRIBUTE_HANDLE_LIST, unsafe.Pointer(&fd[0]), uintptr(len(fd))*unsafe.Sizeof(fd[0]), nil, nil)
		if err != nil {
			return 0, 0, err
		}
	}

	const _PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE uintptr = 22 | 0x00020000
	err = updateProcThreadAttribute(si.ProcThreadAttributeList, 0, _PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE, unsafe.Pointer(*pseudoConsole), unsafe.Sizeof(*pseudoConsole), nil, nil)
	if err != nil {
		return 0, 0, err
	}

	envBlock, err := createEnvBlock(attr.Env)
	if err != nil {
		return 0, 0, err
	}

	pi := new(syscall.ProcessInformation)
	flags := sys.CreationFlags | syscall.CREATE_UNICODE_ENVIRONMENT | _EXTENDED_STARTUPINFO_PRESENT
	if sys.Token != 0 {
		err = syscall.CreateProcessAsUser(sys.Token, argv0p, argvp, sys.ProcessAttributes, sys.ThreadAttributes, willInheritHandles, flags, envBlock, dirp, &si.StartupInfo, pi)
	} else {
		err = syscall.CreateProcess(argv0p, argvp, sys.ProcessAttributes, sys.ThreadAttributes, willInheritHandles, flags, envBlock, dirp, &si.StartupInfo, pi)
	}
	if err != nil {
		return 0, 0, err
	}
	defer syscall.CloseHandle(syscall.Handle(pi.Thread))
	runtime.KeepAlive(fd)
	runtime.KeepAlive(sys)

	return int(pi.ProcessId), uintptr(pi.Process), nil
}

// appendEscapeArg escapes the string s, as per escapeArg,
// appends the result to b, and returns the updated slice.
func appendEscapeArg(b []byte, s string) []byte {
	if len(s) == 0 {
		return append(b, `""`...)
	}

	needsBackslash := false
	hasSpace := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			needsBackslash = true
		case ' ', '\t':
			hasSpace = true
		}
	}

	if !needsBackslash && !hasSpace {
		// No special handling required; normal case.
		return append(b, s...)
	}
	if !needsBackslash {
		// hasSpace is true, so we need to quote the string.
		b = append(b, '"')
		b = append(b, s...)
		return append(b, '"')
	}

	if hasSpace {
		b = append(b, '"')
	}
	slashes := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		default:
			slashes = 0
		case '\\':
			slashes++
		case '"':
			for ; slashes > 0; slashes-- {
				b = append(b, '\\')
			}
			b = append(b, '\\')
		}
		b = append(b, c)
	}
	if hasSpace {
		for ; slashes > 0; slashes-- {
			b = append(b, '\\')
		}
		b = append(b, '"')
	}

	return b
}

// makeCmdLine builds a command line out of args by escaping "special"
// characters and joining the arguments with spaces.
func makeCmdLine(args []string) string {
	var b []byte
	for _, v := range args {
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = appendEscapeArg(b, v)
	}
	return string(b)
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

func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

func normalizeDir(dir string) (name string, err error) {
	ndir, err := syscall.FullPath(dir)
	if err != nil {
		return "", err
	}
	if len(ndir) > 2 && isSlash(ndir[0]) && isSlash(ndir[1]) {
		// dir cannot have \\server\share\path form
		return "", syscall.EINVAL
	}
	return ndir, nil
}

func volToUpper(ch int) int {
	if 'a' <= ch && ch <= 'z' {
		ch += 'A' - 'a'
	}
	return ch
}

func joinExeDirAndFName(dir, p string) (name string, err error) {
	if len(p) == 0 {
		return "", syscall.EINVAL
	}
	if len(p) > 2 && isSlash(p[0]) && isSlash(p[1]) {
		// \\server\share\path form
		return p, nil
	}
	if len(p) > 1 && p[1] == ':' {
		// has drive letter
		if len(p) == 2 {
			return "", syscall.EINVAL
		}
		if isSlash(p[2]) {
			return p, nil
		} else {
			d, err := normalizeDir(dir)
			if err != nil {
				return "", err
			}
			if volToUpper(int(p[0])) == volToUpper(int(d[0])) {
				return syscall.FullPath(d + "\\" + p[2:])
			} else {
				return syscall.FullPath(p)
			}
		}
	} else {
		// no drive letter
		d, err := normalizeDir(dir)
		if err != nil {
			return "", err
		}
		if isSlash(p[0]) {
			return syscall.FullPath(d[:2] + p)
		} else {
			return syscall.FullPath(d + "\\" + p)
		}
	}
}
