// Code generated by extract.py. DO NOT EDIT.

package export

const (
	TCGETS = 0x402c7413
	TCSETS = 0x802c7414
	IXON   = 0x200
	ECHONL = 0x10
	ICANON = 0x100
	ISIG   = 0x80
	IEXTEN = 0x400
	CSIZE  = 0x300
	PARENB = 0x1000
	CS8    = 0x300
	VMIN   = 0x5
	VTIME  = 0x7
)

type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [19]uint8
	Line   uint8
	Ispeed uint32
	Ospeed uint32
}
