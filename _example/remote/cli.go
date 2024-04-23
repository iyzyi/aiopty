package main

import (
	"fmt"
	"strings"
)

type cliArgs struct {
	Mode  string
	Net   string
	Addr  string
	Key   string
	Cmd   string
	Type  string
	Debug bool
}

func cliParse(args []string) (*cliArgs, error) {
	res := &cliArgs{}

	if len(args) == 1 {
		return nil, errPrintUsage
	}

	switch args[1] {
	case "-h", "--help":
		return nil, errPrintUsage
	case "master", "slave":
		res.Mode = args[1]
	default:
		return nil, errInvalidMode
	}

	args = args[2:]
	ptr := 0
	for {
		if ptr == len(args) {
			break
		}

		switch args[ptr] {
		case "-l", "--listen":
			if ptr+1 < len(args) {
				res.Net = "server"
				res.Addr = args[ptr+1]
			} else {
				return nil, errMissingNet
			}
			ptr += 2

		case "-c", "--connect":
			if ptr+1 < len(args) {
				res.Net = "client"
				res.Addr = args[ptr+1]
			} else {
				return nil, errMissingNet
			}
			ptr += 2

		case "-k", "--key":
			if ptr+1 < len(args) {
				res.Key = args[ptr+1]
			} else {
				return nil, errPrintUsage
			}
			ptr += 2

		case "--cmd":
			if ptr+1 < len(args) {
				res.Cmd = args[ptr+1]
			} else {
				return nil, errMissingCmd
			}
			ptr += 2

		case "-t", "--type":
			if ptr+1 < len(args) {
				res.Type = args[ptr+1]
			} else {
				return nil, errPrintUsage
			}
			ptr += 2

		case "-d", "--debug":
			res.Debug = true
			ptr += 1

		default:
			return nil, errPrintUsage
		}
	}

	if res.Net == "" {
		return nil, errMissingNet
	}
	if res.Mode == "slave" {
		_args := strings.Split(res.Cmd, " ")
		if res.Cmd == "" || len(_args) == 0 {
			return nil, errMissingCmd
		}
		if res.Type != "" && res.Type != "nixpty" && res.Type != "conpty" && res.Type != "winpty" {
			return nil, errInvalidType
		}
	}

	return res, nil
}

func usage() {
	fmt.Print(
		"aiopty example: remote shell (https://github.com/iyzyi/aiopty)\n" +
			"Usage: \n" +
			"    1) aiopty master -l/-c ADDRESS [-k KEY] [-d]\n" +
			"    2) aiopty slave -l/-c ADDRESS --cmd CMDLINE [-k KEY] [-t TYPE] [-d]\n" +
			"Mode:\n" +
			"    master: enable terminal\n" +
			"    slave: open a pty\n" +
			"Options:\n" +
			"  -l, --listen ADDRESS\n" +
			"      Listen on ADDRESS using tcp. (must choose either -l or -c)\n" +
			"  -c, --connect ADDRESS\n" +
			"      Connect to ADDRESS using tcp. (must choose either -l or -c)\n" +
			"  -k, --key KEY\n" +
			"      Encrypt data with the KEY. (optional)\n" +
			"  --cmd CMDLINE\n" +
			"      CMDLINE is the command to run. If there are spaces, they must be enclosed\n" +
			"      in double quotes. (for slave mode, required)\n" +
			"  -t, --type TYPE\n" +
			"      Pty TYPE, including nixpty, conpty, winpty. (for slave mode, optional)\n" +
			"  -d, --debug (optional)\n" +
			"  -h, --help (optional)\n",
	)
}

var (
	errPrintUsage  = fmt.Errorf("")
	errInvalidMode = fmt.Errorf("Invalid mode. Must choose a working mode in [master/slave].")
	errMissingNet  = fmt.Errorf("Missing tcp type. Must choose a tcp type in [-l/-c], followed by the network address to listen or to connect.")
	errMissingCmd  = fmt.Errorf("Missing cmdline.")
	errInvalidType = fmt.Errorf("Invalid pty type. Should choose a pty type in [nixpty, conpty, winpty], or do not use the -t option.")
)
