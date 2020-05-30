go get github.com/restic/chunker
go get github.com/BurntSushi/toml
go build dupver.go

$install_folder = "$HOME/.local/bin"
mkdir -p $install_folder
copy dupver $install_folder