package main

import (
	"encoding/binary"
	"github.com/iyzyi/aiopty/utils/log"
	"io"
	"net"
	"sync"
	"time"
)

func createClientConn(addr string) (conn net.Conn, err error) {
	for {
		conn, err = net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			log.Debug("Successfully established TCP connection between %s and %s", conn.LocalAddr(), conn.RemoteAddr())
			return
		} else {
			log.Debug("Failed to connect to %v: %v", addr, err)
			time.Sleep(time.Second)
		}
	}
}

func createServerConn(addr string) (conn net.Conn, err error) {
	// listen
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to listen on %v: %v", addr, err)
		return
	}
	log.Debug("TCP Server is listening on %s...", addr)

	// wait for one netConn
	for {
		conn, err = listener.Accept()
		if err == nil {
			break
		}
	}
	listener.Close()

	log.Debug("Successfully established TCP connection between %s and %s", conn.LocalAddr(), conn.RemoteAddr())
	return
}

type packetType byte

var (
	CMD     packetType = 'c'
	DATA    packetType = 'd'
	SETSIZE packetType = 's'
)

type onRecvFunc func(_type packetType, data []byte) error

func handleRecv(r io.Reader, onRecv onRecvFunc, exit chan struct{}) {
	for {
		lenBytes := make([]byte, 4)
		_, err := io.ReadFull(r, lenBytes)
		if err != nil {
			if err == errCryptoKey {
				log.Error("%v", errCryptoKey)
			}
			break
		}
		length := binary.BigEndian.Uint32(lenBytes)

		if length > 65535 {
			log.Error("illegal package length: %v", length)
			break
		}

		data := make([]byte, length)
		_, err = io.ReadFull(r, data)
		if err != nil {
			break
		}

		_type := packetType(data[0])
		data = data[1:]
		err = onRecv(_type, data)
		if err != nil {
			return
		}
	}
	exit <- struct{}{}
}

func handleSend(w io.Writer, r io.Reader, lock *sync.Mutex, isSlave bool, exit chan struct{}) {
	for {
		data := make([]byte, 4096)
		n, err := r.Read(data)
		if err != nil {
			break
		}

		if n > 0 {
			if isSlave {
				err = sendPacket(w, DATA, data[:n], lock)
				log.Debug("[send:DATA] len: %d", len(data[:n]))
			} else {
				err = sendPacket(w, CMD, data[:n], lock)
			}
			if err != nil {
				break
			}
		}
	}
	exit <- struct{}{}
}

func sendPacket(w io.Writer, _type packetType, data []byte, lock *sync.Mutex) (err error) {
	lock.Lock()
	defer lock.Unlock()

	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(data)+1))

	_, err = w.Write(lenBytes)
	if err != nil {
		return
	}

	_, err = w.Write([]byte{byte(_type)})
	if err != nil {
		return
	}

	_, err = w.Write(data)
	if err != nil {
		return
	}
	return
}
