# aiopty

<p>
    <a href="https://github.com/iyzyi/aiopty/releases"><img src="https://img.shields.io/github/release/iyzyi/aiopty.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/iyzyi/aiopty"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
    <a href="https://goreportcard.com/report/github.com/iyzyi/aiopty"><img src="https://goreportcard.com/badge/github.com/iyzyi/aiopty" alt="Go Report Card"></a>
</p>

(All in One) Pty & Terminal package for Go with an encrypted remote shell as example.

Eng | [中文](https://github.com/iyzyi/aiopty/blob/master/README_CN.md)

## Features

* Provide unified Pty and Terminal support for Windows (WinXP~Win11), Linux, Darwin, Dragonfly, FreeBSD, NetBSD, OpenBSD, Solaris.
* Optimize the poor experience on Windows using [ConPty](https://devblogs.microsoft.com/commandline/windows-command-line-introducing-the-windows-pseudo-console-conpty/) and [WinPty](https://github.com/rprichard/winpty).
* Integrate Pty behavior on Unix-like systems into [NixPty](https://en.wikipedia.org/wiki/Pseudoterminal), and respectively provide implementations based on native Go and CGO.
* No package dependencies, compatible with the majority of Go versions.
* Based on this package, develop and release an encrypted remote shell as example that supports both forward and reverse TCP connections.

## Example

```go
package main

import (
	"github.com/iyzyi/aiopty/pty"
	"github.com/iyzyi/aiopty/term"
	"github.com/iyzyi/aiopty/utils/log"
	"io"
	"os"
)

func main() {
	// open a pty with options
	opt := &pty.Options{
		Path: "cmd.exe",
		Args: []string{"cmd.exe", "/c", "powershell.exe"},
		Dir:  "",
		Env:  nil,
		Size: &pty.WinSize{
			Cols: 120,
			Rows: 30,
		},
		Type: "",
	}
	p, err := pty.OpenWithOptions(opt)

	// You can also open a pty simply like this:
	// p, err := pty.Open(path)

	if err != nil {
		log.Error("Failed to create pty: %v", err)
		return
	}
	defer p.Close()

	// When the terminal window size changes, synchronize the size of the pty
	onSizeChange := func(cols, rows uint16) {
		size := &pty.WinSize{
			Cols: cols,
			Rows: rows,
		}
		p.SetSize(size)
	}

	// enable terminal
	t, err := term.Open(os.Stdin, os.Stdout, onSizeChange)
	if err != nil {
		log.Error("Failed to enable terminal: %v", err)
		return
	}
	defer t.Close()

	// start data exchange between terminal and pty
	exit := make(chan struct{}, 2)
	go func() { io.Copy(p, t); exit <- struct{}{} }()
	go func() { io.Copy(t, p); exit <- struct{}{} }()
	<-exit
}
```

## Note

* If you need to use this package in a Go version lower than `1.16`, please note that:
  * If you need to use Go modules, modify the Go version in the `go.mod` file to `1.11`. (due to [golang/go#43980](https://github.com/golang/go/issues/43980))
  * If you need to use WinPty on Windows, please manually copy the winpty library files from `aiopty/pty/winpty/bin/{arch}` to the directory where your current program resides. Please note that the architecture of the winpty library files needs to match the architecture of the current program.

## Release

Based on this package, I develop and release an encrypted remote shell as example that supports both forward and reverse TCP connections. See [Releases](https://github.com/iyzyi/aiopty/releases). Here is how to use it.

```
Usage:
    1) aiopty master -l/-c ADDRESS [-k KEY] [-d]
    2) aiopty slave -l/-c ADDRESS --cmd CMDLINE [-k KEY] [-t TYPE] [-d]
Mode:
    master: enable terminal
    slave: open a pty
Options:
  -l, --listen ADDRESS
      Listen on ADDRESS using tcp. (must choose either -l or -c)
  -c, --connect ADDRESS
      Connect to ADDRESS using tcp. (must choose either -l or -c)
  -k, --key KEY
      Encrypt data with the KEY. (optional)
  --cmd CMDLINE
      CMDLINE is the command to run. If there are spaces, they must be enclosed
      in double quotes. (for slave mode, required)
  -t, --type TYPE
      Pty TYPE, including nixpty, conpty, winpty. (for slave mode, optional)
  -d, --debug (optional)
  -h, --help (optional)
```

for example:

```
aiopty-windows-amd64.exe slave -l 0.0.0.0:50505 --cmd "cmd.exe /c powershell.exe" -k secret
aiopty-windows-amd64.exe master -c 127.0.0.1:50505 -k secret
```

## Contribution

Welcome to your issues and pull requests.

## License

This package is released under the [MIT License](https://github.com/iyzyi/aiopty/blob/master/LICENSE).

