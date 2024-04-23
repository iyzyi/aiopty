# aiopty

(All in One) Pty & Terminal package for Go with an encrypted remote shell as example.

[Eng](https://github.com/iyzyi/aiopty) | 中文

## Features

* 为Windows (WinXP~Win11), Linux, Darwin, Dragonfly, FreeBSD, NetBSD, OpenBSD, Solaris提供统一的Pty和终端支持。
* 使用[ConPty](https://devblogs.microsoft.com/commandline/windows-command-line-introducing-the-windows-pseudo-console-conpty/)和[WinPty](https://github.com/rprichard/winpty)优化Windows上的糟糕体验。
* 将类Unix系统中的Pty行为集成到[NixPty](https://en.wikipedia.org/wiki/Pseudoterminal)中，并分别基于原生Go和CGO提供相关实现。
* 无任何package依赖，适用于绝大多数Go版本。
* 基于此package，开发并发布一个支持正向和反向TCP连接的加密远程shell作为示例。

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

* 如果需要在Go版本低于1.16的环境中使用此package，请注意：
  * 如果需要使用Go模块，请将`go.mod`中的Go版本修改为`1.11`。(由[golang/go#43980](https://github.com/golang/go/issues/43980)导致)
  * 如果需要在Windows上使用WinPty，请手动将`aiopty/pty/winpty/bin/{arch}`中的winpty库文件复制到当前程序所在目录。注意winpty库文件的架构需要和当前程序的架构一致。

## Release

基于此package，我开发并发布了一个支持正向和反向TCP连接的加密远程shell作为示例。可见[Releases](https://github.com/iyzyi/aiopty/releases)。下面介绍如何使用。

```
使用方法:
    1) aiopty master -l/-c ADDRESS [-k KEY] [-d]
    2) aiopty slave -l/-c ADDRESS --cmd CMDLINE [-k KEY] [-t TYPE] [-d]
模式:
    master: 启用终端
    slave: 打开pty
选项:
  -l, --listen ADDRESS
      在ADDRESS上监听。(必须在-l或-c中二选一)
  -c, --connect ADDRESS
      连接到ADDRESS。(必须在-l或-c中二选一)
  -k, --key KEY
      使用KEY加密数据。(可选)
  --cmd CMDLINE
      CMDLINE是要运行的命令。如果有空格，必须用双引号括起来。(用于slave模式，必需)
  -t, --type TYPE
      Pty类型，包括nixpty、conpty、winpty。(用于slave模式，可选)
  -d, --debug (可选)
  -h, --help (可选)
```

例如：

```
aiopty-windows-amd64.exe slave -l 0.0.0.0:50505 --cmd "cmd.exe /c powershell.exe" -k secret
aiopty-windows-amd64.exe master -c 127.0.0.1:50505 -k secret
```

## Contribution

欢迎您的issues和pull requests。

## License

该package基于[MIT许可](https://github.com/iyzyi/aiopty/blob/master/LICENSE)发布。

