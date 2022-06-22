#go get github.com/akbarnes/dupver/cmd/dupver
# go build -o dupver.exe cmd/dupver/main.go
go install github.com/akbarnes/dupver/cmd/dupver

if ($? -eq $false) {
    Throw "Error building dupver executable"
}

# move -Force dupver.exe $HOME/go/bin/
