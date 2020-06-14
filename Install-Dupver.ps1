go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go workdir.go repo.go snapshot.go pack.go randstring.go util.go
echo "Done building"

$InstallFolder = "$HOME\AppData\Local\Executables"
if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
if (test-path $InstallFolder\dupver.exe) { del $InstallFolder\dupver.exe }
copy dupver.exe $InstallFolder
echo "Copied executable to $InstallFolder"
$ParentDir = "$HOME\OneDrive\Documents"
copy Test-Dupver.ps1 $ParentDir
echo "Copied to $ParentDir"