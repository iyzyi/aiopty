package native

import (
	"bytes"
)

// Note: This variable is used to obtain the ptsname string through a system call.
// The maximum length of ptsname is defined by different macros on various systems.
// For convenience, I have set a unified, larger predefined length.
const ptsnameLen = 2048

// ByteSliceToString returns a string form of the text represented by the slice s,
// with a terminating NUL and any bytes after the NUL removed.
func ByteSliceToString(s []byte) string {
	if i := bytes.IndexByte(s, 0); i != -1 {
		s = s[:i]
	}
	return string(s)
}
