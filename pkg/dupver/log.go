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
            fmt.Println(snap.SnapshotID)
        } else {
            if len(snap.SnapshotLocalTime) > 0 {
                fmt.Printf("%d) Local Time: %s\n", i+1, snap.SnapshotLocalTime)
            } else {
                fmt.Printf("%d) Time: %s\n", i+1, snap.SnapshotTime)
            }

            fmt.Printf("ID: %s\n", snap.SnapshotID[0:9])

            if len(snap.Username) > 0 {
                fmt.Printf("Username: %s\n", snap.Username)
            }

            if len(snap.Message) > 0 {
                fmt.Printf("Message: %s\n", snap.Message)
            }

            fmt.Printf("\n")
        }
	}
}

func MatchSnapshot(commitID string) (Snapshot, error) {
	snapshotGlob := filepath.Join(WorkingDirectory, ".dupver", "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)

    for _, snapshotPath := range snapshotPaths {
        if strings.Contains(snapshotPath, commitID) {
	        return ReadSnapshotJson(snapshotPath), nil
        }
    }

    return Snapshot{}, errors.New("No matching snapshots")
}

func LogSingleSnapshot(commitID string) {
    var snap Snapshot
    var err error

    if commitID == "last" || commitID == "latest" {
        snap = ReadHead()
    } else {
        snap, err = MatchSnapshot(commitID)

        if err != nil {
            fmt.Fprintf(os.Stderr, "No matching snapshot paths\n")
            os.Exit(1)
        }
    }

	snapFiles := snap.ReadFilesHash()

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
            i++
        }
	}
}
