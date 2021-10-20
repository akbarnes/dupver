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

	status := make(map[string]string)

	for fileName, _ := range snap.FileModTimes {
		status[fileName] = "-"
	}

	// workingDirectory, err := os.Getwd()
	// Check(err)
	workingDirectory := "."
	head := ReadHead()

	var DiffFile = func(fileName string, info os.FileInfo, err error) error {
		fileName = strings.TrimSuffix(fileName, "\n")

		if ExcludedFile(fileName, info, filters) {
			return nil
		}

		props, err := os.Stat(fileName)

		if err != nil {
			if VerboseMode {
				fmt.Printf("Skipping unreadable file %s\n", fileName)
			}

			return nil
		}

		modTime := props.ModTime().Format("2006-01-02T15-04-05")

		if headModTime, ok := head.FileModTimes[fileName]; ok {
			if modTime == headModTime {
				status[fileName] = "="
			} else {
				status[fileName] = "M"
			}
		} else {
			status[fileName] = "+"
		}

		return nil
	}

	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)
	filepath.Walk(workingDirectory, DiffFile)

	for fileName, fileStatus := range status {
		if fileStatus == "=" && !VerboseMode {
			continue
		}

		fmt.Printf("%s %s\n", fileStatus, fileName)
	}
}
