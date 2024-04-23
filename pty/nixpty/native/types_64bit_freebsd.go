//go:build amd64 || arm64 || riscv64
// +build amd64 arm64 riscv64

package native

const FIODGNAME = 0x80106678

type fiodgnameArg struct {
	len     int32
	padding uint32 //memory alignment
	buf     uintptr
}
