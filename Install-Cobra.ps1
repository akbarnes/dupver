$GetLibs = $false
$InstallDupver = $true

if ($GetLibs) {
    go get -u -v github.com/spf13/cobra/cobra
    go get github.com/restic/chunker
    go get github.com/BurntSushi/toml
}

go build -o dupver.exe main.go
echo "Done building"

$InstallFolder = "$HOME\AppData\Local\Executables"

if (-not (test-path $InstallFolder)) { mkdir $Installfolder }
copy dupver.exe $InstallFolder
echo "Copied to $InstallFolder"

$ParentDir = "$HOME\Documents"
copy Test-Dupver.ps1 $ParentDir
echo "Copied Test-Dupver.ps1 to $ParentDir"