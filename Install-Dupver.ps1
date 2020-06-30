$InstallDupver = $true


if ($InstallDupver) {
    go install
} else {
    if (-not (test-path bin)) { mkdir bin }
    go build -o bin/dupver
}

echo "Done building"
