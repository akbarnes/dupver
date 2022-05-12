#go get github.com/akbarnes/dupver/cmd/dupver
go build -o dupver.exe cmd/dupver/main.go
move -Force dupver.exe $HOME/go/bin/
