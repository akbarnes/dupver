go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go commit.go pack.go config.go randstring.go util.go
echo "Done building"

$InstallFolder = "$HOME\AppData\Local\Executables"
if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
if (test-path $InstallFolder\dupver.exe) { del $InstallFolder\dupver.exe }
move dupver.exe $InstallFolder

$ParentDir = "$HOME\Data\Static\Electric\Transmission Atlas"
copy Test-Dupver.ps1 $ParentDir
echo "Copied to $ParentDir"