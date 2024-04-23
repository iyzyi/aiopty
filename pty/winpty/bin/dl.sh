#!/bin/bash

url="https://github.com/rprichard/winpty/releases/download/0.4.3/"
ia32="winpty-0.4.3-msys2-2.7.0-ia32"
x64="winpty-0.4.3-msys2-2.7.0-x64"
dll="winpty.dll"
agent="winpty-agent.exe"

dl() {
    local label=$1
    local arch=$2

    wget "$url$label.tar.gz"
    if [ $? -ne 0 ]; then
        return 1      
    fi

    mkdir temp
    tar -xzvf $label.tar.gz -C temp
    mkdir $arch
    mv temp/$label/bin/$dll $arch/
    mv temp/$label/bin/$agent $arch/
    rm -rf temp *.tar.gz
}

cd $(dirname "$0")
rm -rf ia32 x64 temp *.tar.gz

dl $ia32 ia32
dl $x64 x64