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
