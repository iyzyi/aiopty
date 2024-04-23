package main

import (
	"fmt"
	"github.com/iyzyi/aiopty/utils/log"
	"net"
	"os"
)

func main() {
	args, err := cliParse(os.Args)
	if err != nil {
		if err == errPrintUsage {
			usage()
		} else {
			fmt.Println(err)
		}
		return
	}

	if args.Debug {
		log.Level = log.DEBUG
	}

	var conn net.Conn
	if args.Net == "server" {
		conn, err = createServerConn(args.Addr)
		if err != nil {
			return
		}
		defer conn.Close()
	} else if args.Net == "client" {
		conn, err = createClientConn(args.Addr)
		if err != nil {
			return
		}
		defer conn.Close()
	}

	if args.Mode == "master" {
		m, _ := NewMaster(conn, args)
		m.Start()
	} else if args.Mode == "slave" {
		s, _ := NewSlave(conn, args)
		s.Start()
	}
}
