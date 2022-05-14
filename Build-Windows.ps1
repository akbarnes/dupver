$MainPath = Join-Path "cmd" "dupver" "main.go"
$BinPath = Join-Path "bin" "dupver-win-x64.exe"
go build -o $BinPath $MainPath

