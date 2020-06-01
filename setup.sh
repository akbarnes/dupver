#!/bin/bash
go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go commit.go pack.go config.go

install_folder="$HOME/.local/bin"
mkdir -p $install_folder
cp dupver $install_folder
