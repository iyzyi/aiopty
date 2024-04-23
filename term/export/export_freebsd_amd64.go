// Code generated by extract.py. DO NOT EDIT.

package export

const (
	TIOCGETA = 0x402c7413
	TIOCSETA = 0x802c7414
	IGNBRK   = 0x1
	BRKINT   = 0x2
	PARMRK   = 0x8
	ISTRIP   = 0x20
	INLCR    = 0x40
	IGNCR    = 0x80
	ICRNL    = 0x100
	IXON     = 0x200
	OPOST    = 0x1
	ECHO     = 0x8
	ECHONL   = 0x10
	ICANON   = 0x100
	ISIG     = 0x80
	IEXTEN   = 0x400
	CSIZE    = 0x300
	PARENB   = 0x1000
	CS8      = 0x300
	VMIN     = 0x10
	VTIME    = 0x11
)

type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
}
