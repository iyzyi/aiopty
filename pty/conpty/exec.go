//go:build windows
// +build windows

package conpty

import (
	"errors"
	"syscall"
	"unsafe"
)

// Some syscall on windows.
// Copied from Golang 1.21 (with slight modifications), should be universally compatible across all golang versions.

// newProcThreadAttributeList allocates new PROC_THREAD_ATTRIBUTE_LIST, with
// the requested maximum number of attributes, which must be cleaned up by
// deleteProcThreadAttributeList.
func newProcThreadAttributeList(maxAttrCount uint32) (*_PROC_THREAD_ATTRIBUTE_LIST, error) {
	var size uintptr
	err := initializeProcThreadAttributeList(nil, maxAttrCount, 0, &size)
	if err != syscall.ERROR_INSUFFICIENT_BUFFER {
		if err == nil {
			//return nil, errorspkg.New("unable to query buffer size from InitializeProcThreadAttributeList")
			return nil, errors.New("unable to query buffer size from InitializeProcThreadAttributeList")
		}
		return nil, err
	}
	// size is guaranteed to be ≥1 by initializeProcThreadAttributeList.
	al := (*_PROC_THREAD_ATTRIBUTE_LIST)(unsafe.Pointer(&make([]byte, size)[0]))
	err = initializeProcThreadAttributeList(al, maxAttrCount, 0, &size)
	if err != nil {
		return nil, err
	}
	return al, nil
}

// syscall functions
var (
	//modkernel32                           = syscall.NewLazyDLL(sysdll.Add("kernel32.dll"))
	//modntdll                              = syscall.NewLazyDLL(sysdll.Add("ntdll.dll"))
	modkernel32                           = syscall.NewLazyDLL("kernel32.dll")
	modntdll                              = syscall.NewLazyDLL("ntdll.dll")
	procInitializeProcThreadAttributeList = modkernel32.NewProc("InitializeProcThreadAttributeList")
	procUpdateProcThreadAttribute         = modkernel32.NewProc("UpdateProcThreadAttribute")
	procDeleteProcThreadAttributeList     = modkernel32.NewProc("DeleteProcThreadAttributeList")
	procRtlGetNtVersionNumbers            = modntdll.NewProc("RtlGetNtVersionNumbers")
)

func initializeProcThreadAttributeList(attrlist *_PROC_THREAD_ATTRIBUTE_LIST, attrcount uint32, flags uint32, size *uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(procInitializeProcThreadAttributeList.Addr(), 4, uintptr(unsafe.Pointer(attrlist)), uintptr(attrcount), uintptr(flags), uintptr(unsafe.Pointer(size)), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func updateProcThreadAttribute(attrlist *_PROC_THREAD_ATTRIBUTE_LIST, flags uint32, attr uintptr, value unsafe.Pointer, size uintptr, prevvalue unsafe.Pointer, returnedsize *uintptr) (err error) {
	r1, _, e1 := syscall.Syscall9(procUpdateProcThreadAttribute.Addr(), 7, uintptr(unsafe.Pointer(attrlist)), uintptr(flags), uintptr(attr), uintptr(value), uintptr(size), uintptr(prevvalue), uintptr(unsafe.Pointer(returnedsize)), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func deleteProcThreadAttributeList(attrlist *_PROC_THREAD_ATTRIBUTE_LIST) {
	syscall.Syscall(procDeleteProcThreadAttributeList.Addr(), 1, uintptr(unsafe.Pointer(attrlist)), 0, 0)
	return
}

func rtlGetNtVersionNumbers(majorVersion *uint32, minorVersion *uint32, buildNumber *uint32) {
	syscall.Syscall(procRtlGetNtVersionNumbers.Addr(), 3, uintptr(unsafe.Pointer(majorVersion)), uintptr(unsafe.Pointer(minorVersion)), uintptr(unsafe.Pointer(buildNumber)))
	return
}

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}
