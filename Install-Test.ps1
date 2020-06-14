go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go repo.go workdir.go snapshot.go pack.go randstring.go util.go  
echo "Done building"

$InstallFolder = "$HOME\AppData\Local\Executables"
if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
if (test-path $InstallFolder\dupver.exe) { del $InstallFolder\dupver.exe }
copy dupver.exe $InstallFolder
echo "Copied dupver.exe to $InstallFolder"

$ParentDir = "$HOME\Documents\Books"
copy Test-Dupver.ps1 $ParentDir
echo "Copied Test-Dupver.ps1 to $ParentDir"