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
        if QuietMode {
            fmt.Println(snap.SnapshotId)
        } else { 
            if len(snap.SnapshotLocalTime) > 0 {
                fmt.Printf("%d) Local Time: %s\n", i+1, snap.SnapshotLocalTime)
            } else {
                fmt.Printf("%d) Time: %s\n", i+1, snap.SnapshotTime)
            }

            fmt.Printf("ID: %s\n", snap.SnapshotId[0:9])

            if len(snap.Message) > 0 {
                fmt.Printf("Message: %s\n\n", snap.Message)
            }
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

    i := 1

	for fileName, fileProps := range snapFiles {
        if QuietMode {
		    fmt.Println(fileName)
        } else {
            fmt.Printf("%d) %s\n", i, fileName)
       
            if len(fileProps.ModLocalTime) > 0 {
                fmt.Printf("Local Modified: %s\n", fileProps.ModLocalTime)
            } else {
                fmt.Printf("Modified: %s\n", fileProps.ModTime)
            }

            fmt.Printf("Size: %0.3f MB\n\n", float64(fileProps.Size)/(1024.0*1024.0))
            i += 1
        }
	}
}
