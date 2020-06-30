$InstallDupver = $true
    
go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go get github.com/spf13/cobra


go build main.go workdir.go repo.go snapshot.go pack.go randstring.go util.go


echo "Done building"
