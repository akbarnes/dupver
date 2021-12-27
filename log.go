package dupver

import (
	"fmt"
    "strings"
    "errors"
    "os"
	"path/filepath"
)

func LogAllSnapshots() {
	for i, snap := range ReadAllSnapshots() {
		// Time: 2021/05/08 08:57:46
		// Message: specify workdir path explicitly
		fmt.Fprintf(os.Stderr, "%d) Time: %s\n", i+1, snap.SnapshotTime)
        fmt.Printf("%s\n", snap.SnapshotId[0:9])

		if len(snap.Message) > 0 {
			fmt.Fprintf(os.Stderr, "Message: %s\n\n", snap.Message)
		}
	}
}

func MatchSnapshot(commitId string) (Snapshot, error) {
	snapshotGlob := filepath.Join(".dupver", "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)

    for _, snapshotPath := range snapshotPaths {
        if strings.Contains(snapshotPath, commitId) {
	        return ReadSnapshotJson(snapshotPath), nil
        }
    }

    return Snapshot{}, errors.New("No matching snapshots")
}

func LogSingleSnapshot(commitId string) {
    snap, err := MatchSnapshot(commitId)

    if err != nil {
        fmt.Fprintf(os.Stderr, "No matching snapshot paths")
        os.Exit(1)
    }


	snapFiles := snap.ReadFilesList()

	for fileName, _ := range snapFiles {
		fmt.Println(fileName)
	}
}
