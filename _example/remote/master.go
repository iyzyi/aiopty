package main

import (
	"encoding/binary"
	"fmt"
	"github.com/iyzyi/aiopty/term"
	"io"
	"net"
	"os"
	"sync"
)

type Master struct {
	netConn  net.Conn
	encConn  io.ReadWriter
	conn     io.ReadWriter
	args     *cliArgs
	term     *term.Term
	sendLock sync.Mutex
}

func NewMaster(netConn net.Conn, args *cliArgs) (m *Master, err error) {
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

	m = &Master{
		netConn: netConn,
		encConn: encConn,
		conn:    conn,
		args:    args,
	}
	return
}

func (m *Master) Start() (err error) {
	// enable terminal
	onSizeChange := func(cols, rows uint16) {
		sizeUint32 := uint32(cols) + (uint32(rows) << 16)
		sizeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(sizeBytes, sizeUint32)
		sendPacket(m.conn, SETSIZE, sizeBytes, &m.sendLock)
	}
	m.term, err = term.Open(os.Stdin, os.Stdout, onSizeChange)
	if err != nil {
		return
	}
	defer m.term.Close()

	// data exchange
	exit := make(chan struct{}, 2)
	go handleSend(m.conn, m.term, &m.sendLock, false, exit)
	go handleRecv(m.conn, m.masterRecv, exit)
	<-exit

	return
}

func (m *Master) masterRecv(_type packetType, data []byte) (err error) {
	switch _type {
	case DATA:
		_, err = m.term.Write(data)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("error packet type")
		return
	}
	return
}
