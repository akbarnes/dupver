$InstallDupver = $true

go get -u -v github.com/spf13/cobra/cobra
go get github.com/restic/chunker
go get github.com/BurntSushi/toml

if ($InstallDupver) {
    go install dupver.go workdir.go repo.go snapshot.go pack.go randstring.go util.go
} else {
    go build dupver.go workdir.go repo.go snapshot.go pack.go randstring.go util.go
}

echo "Done building"
