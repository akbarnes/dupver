package dupver

import (
	"fmt"
    "strings"
	"path/filepath"
)

func LogAllSnapshots() {
	for i, snap := range ReadAllSnapshots() {
		// Time: 2021/05/08 08:57:46
		// Message: specify workdir path explicitly
		fmt.Printf("%d) Time: %s\n", i+1, snap.SnapshotTime)
        fmt.Printf("ID: %s\n", snap.SnapshotId[0:9])

		if len(snap.Message) > 0 {
			fmt.Printf("Message: %s\n\n", snap.Message)
		}
	}
}

func LogSingleSnapshot(commitId string) {
	snapshotGlob := filepath.Join(".dupver", "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)

    var snapshotPath string

    for _, snapshotPath := range snapshotPaths {
        if strings.Contains(snapshotPath, commitId) {
            break
        }
    }

	snap := ReadSnapshotJson(snapshotPath)

	snapFiles := snap.ReadFilesList()

	for fileName, _ := range snapFiles {
		fmt.Println(fileName)
	}
}
