package dupver

import (
	"fmt"
	"path/filepath"
	"sort"
)

func LogAllSnapshots() {
	snapshotGlob := filepath.Join(".dupver", "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)
	snaps := []Snapshot{}

	for _, snapshotPath := range snapshotPaths {
		snaps = append(snaps, ReadSnapshotJson(snapshotPath))
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].SnapshotTime < snaps[j].SnapshotTime
	})

	for i, snap := range snaps {
		// Time: 2021/05/08 08:57:46
		// Message: specify workdir path explicitly
		fmt.Printf("%3d) Time: %s\n", i+1, snap.SnapshotTime)

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
	snap := ReadSnapshotJson(snapshotPath)

	snapFiles := snap.ReadFilesList()

	for fileName, _ := range snapFiles {
		fmt.Println(fileName)
	}
}
