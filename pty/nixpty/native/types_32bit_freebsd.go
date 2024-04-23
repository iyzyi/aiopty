//go:build 386 || arm
// +build 386 arm

package native

const FIODGNAME = 0x80086678

type fiodgnameArg struct {
	len int32
	buf uintptr
}

// ./tool/main.cpp output:
// ********** 64Bit ARCH **********
// [DragonFly] FIODNAME: 0x80106678
// [FreeBSD] FIODGNAME: 0x80106678
// [NetBSD] TIOCPTSNAME: 0x40287448
// [OpenBSD] PTMGET: 0x40287401
// ********** 32Bit ARCH **********
// [DragonFly] FIODNAME: 0x80086678
// [FreeBSD] FIODGNAME: 0x80086678
// [NetBSD] TIOCPTSNAME: 0x40287448
// [OpenBSD] PTMGET: 0x40287401
