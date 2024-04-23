package main

import (
	"encoding/binary"
	"fmt"
	"github.com/iyzyi/aiopty/pty"
	"github.com/iyzyi/aiopty/utils/log"
	"io"
	"net"
	"strings"
	"sync"
)

type Slave struct {
	netConn  net.Conn
	encConn  io.ReadWriter
	conn     io.ReadWriter
	args     *cliArgs
	pty      *pty.Pty
	sendLock sync.Mutex
}

func NewSlave(netConn net.Conn, args *cliArgs) (s *Slave, err error) {
	var encConn, conn io.ReadWriter
	if args.Key != "" {
		encConn, err = NewCryptoReadWriter(netConn, []byte(args.Key))
		if err != nil {
			return
		}
		conn = encConn
	} else {
		conn = netConn
	}

	s = &Slave{
		netConn: netConn,
		encConn: encConn,
		conn:    conn,
		args:    args,
	}
	return
}

func (s *Slave) Start() (err error) {
	args := strings.Split(s.args.Cmd, " ")

	// open a pty
	s.pty, err = pty.OpenWithOptions(
		&pty.Options{
			Path: args[0],
			Args: args,
			Type: pty.PtyType(s.args.Type)})
	if err != nil {
		log.Error("Failed to create pty: %v", err)
		return
	}
	defer s.pty.Close()

	// data exchange
	exit := make(chan struct{}, 2)
	go handleSend(s.conn, s.pty, &s.sendLock, true, exit)
	go handleRecv(s.conn, s.slaveRecv, exit)
	<-exit

	return
}

func (s *Slave) slaveRecv(_type packetType, data []byte) (err error) {
	switch _type {
	case CMD:
		_, err = s.pty.Write(data)
		if err != nil {
			return
		}
		log.Debug("[recv:CMD] len: %d", len(data))

	case SETSIZE:
		if len(data) != 4 {
			err = fmt.Errorf("error data length for setsize")
			return
		}
		value := binary.BigEndian.Uint32(data)
		size := &pty.WinSize{
			Rows: uint16((value >> 16) & 0xffff),
			Cols: uint16(value & 0xffff),
		}
		s.pty.SetSize(size)
		log.Debug("[recv:SETSIZE] %d X %d", size.Cols, size.Rows)

	default:
		err = fmt.Errorf("error packet type")
		return
	}
	return
}
