package dupver

import (
	"fmt"
	"path/filepath"
)

func LogAllSnapshots() {
	snapshotGlob := filepath.Join(".dupver", "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)

	for i, snapshotPath := range snapshotPaths {
		snap := ReadSnapshotFile(snapshotPath)
		// Time: 2021/05/08 08:57:46
		// Message: specify workdir path explicitly
		fmt.Printf("%3d) Time: %s\n", i+1, snap.Time)

		if len(snap.Message) > 0 {
			fmt.Printf("Message: %s\n\n", snap.Message)
		}
	}
}

func LogSingleSnapshot(snapshotNum int) {
	snapshotGlob := filepath.Join(".dupver", "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)

	snapshotPath := snapshotPaths[snapshotNum-1]
	snap := ReadSnapshotFile(snapshotPath)

	for file, _ := range snap.FileModTimes {
		fmt.Println(file)
	}
}
