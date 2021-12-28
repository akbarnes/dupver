package dupver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DiffSnapshot(snapId string, filters []string) {
	var snap Snapshot

	if len(snapId) > 0 {
		snap = ReadSnapshot(snapId)
	} else {
		snap = ReadHead()
	}

	snap.Diff(filters)
}

func (snap Snapshot) Diff(filters []string) {
	status := make(map[string]string)
	snapFiles := snap.ReadFilesList()

	for fileName, _ := range snapFiles {
		status[fileName] = "-"
	}

	// workingDirectory, err := os.Getwd()
	// Check(err)
	workingDirectory := "."

	var DiffFile = func(fileName string, info os.FileInfo, err error) error {
		fileName = strings.TrimSuffix(fileName, "\n")

		if ExcludedFile(fileName, info, filters) {
			return nil
		}

		props, err := os.Stat(fileName)

		if err != nil {
			if DebugMode {
				fmt.Fprintf(os.Stderr, "Skipping unreadable file %s\n", fileName)
			}

			return nil
		}

		modTime := props.ModTime().UTC().Format("2006-01-02T15-04-05")

		if snapFile, ok := snapFiles[fileName]; ok {
			if modTime == snapFile.ModTime {
				status[fileName] = "="
			} else {
                if DebugMode {
                    fmt.Printf("%s -> %s: %s\n", snapFile.ModTime, modTime, fileName)
                }

				status[fileName] = "M"
			}
		} else {
			status[fileName] = "+"
		}

		return nil
	}

	// fmt.Fprintf(os.Stderr, "No changes detected in %s for commit %s\n", workDir, snapshot.ID)
	filepath.Walk(workingDirectory, DiffFile)

	for fileName, fileStatus := range status {
		if fileStatus == "=" && !VerboseMode {
			continue
		}

		fmt.Printf("%s %s\n", fileStatus, fileName)
	}
}
