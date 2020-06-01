go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go commit.go pack.go config.go randstring.go

$InstallFolder = "$HOME\AppData\Local\Executables"
if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
copy dupver.exe $InstallFolder
