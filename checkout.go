package dupver

import (
	"fmt"
	"os"
    "time"
    "os/exec"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
)

func CheckoutSnapshot(commitId string, outputFolder string, filter string) {
    var snap Snapshot
    var err error

    if commitId == "last" {
        snap = ReadHead()
    } else {
        snap, err = MatchSnapshot(commitId)

        if err != nil {
            fmt.Fprintf(os.Stderr, "No matching snapshot paths\n")
            os.Exit(1)
        }
    }


	fmt.Fprintf(os.Stderr, "Checking out %s\n", snap.SnapshotId[0:9])
    snap.Checkout(outputFolder, filter)
}

func DiffToolSnapshot(diffTool string) {
    snap := ReadHead()
	fmt.Fprintf(os.Stderr, "Comparing %s\n", snap.SnapshotId[0:9])
    snap.DiffTool(diffTool)
}

func (snap Snapshot) DiffTool(diffTool string) {
    home, err := os.UserHomeDir()
    Check(err)
    tempFolder := filepath.Join(home, ".dupver", "temp", RandHexString(24))
    snap.Checkout(tempFolder, "*")
    cmd := exec.Command(diffTool, tempFolder, ".")
    cmd.Run()
}

func DiffToolSnapshotFile(fileName string, diffTool string) {
    snap := ReadHead()
	fmt.Fprintf(os.Stderr, "Comparing %s/%s\n", snap.SnapshotId[0:9], fileName)
    snap.DiffToolFile(fileName, diffTool)
}

func (snap Snapshot) DiffToolFile(fileName string, diffTool string) {
    home, err := os.UserHomeDir()
    Check(err)
    tempFolder := filepath.Join(home, ".dupver", "temp", RandHexString(24))
    snap.Checkout(tempFolder, fileName)
    tempFile := filepath.Join(tempFolder, fileName)
    cmd := exec.Command(diffTool, tempFile, fileName) 
    cmd.Run()
}

func (snap Snapshot) Checkout(outputFolder string, filter string) {
	os.MkdirAll(outputFolder, 0777)
	snapFiles := snap.ReadFilesList()
	packs := ReadTrees()

	for fileName, fileProps := range snapFiles {
        matched, err := doublestar.PathMatch(filter, fileName)

        if err != nil && VerboseMode {
            fmt.Fprintf(os.Stderr, "Error matching %s\n", filter)
        }

        if !matched {
            if DebugMode {
                fmt.Fprintf(os.Stderr, "Skipping file %s\n", fileName)
            }

            continue
        }

		fileDir := filepath.Dir(fileName)
		outDir := outputFolder

		if fileDir != "." {
			outDir = filepath.Join(outputFolder, fileDir)
			fmt.Fprintf(os.Stderr, "Creating folder %s\n", outDir)
			os.MkdirAll(outDir, 0777)
		}

		outPath := filepath.Join(outputFolder, fileName)
		outFile, err := os.Create(outPath)

		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error creating %s, skipping\n", outPath)
			fmt.Fprintf(os.Stderr, "Error creating %s, skipping\n", outPath)
			continue
		}

		defer outFile.Close()

		for _, chunkId := range fileProps.ChunkIds {
			packId := packs[chunkId]

			if DebugMode {
				fmt.Fprintf(os.Stderr, "Extracting:\n  Pack %s\n  Chunk %s\n  to %s\n\n", packId, chunkId, outPath)
			}

			ExtractChunkFromPack(outFile, chunkId, packId)
		}

        mtime, err := time.Parse("2006-01-02T15-04-05", fileProps.ModTime)

        if err == nil {
            os.Chtimes(outPath, mtime, mtime)
        } else {
            fmt.Fprintf(os.Stderr, "Error parsing time %s for file %s, not setting", fileProps.ModTime, fileName)
        }

        if VerboseMode {
		    fmt.Printf("Restored %s to %s\n", fileName, outPath)
        } else {
            fmt.Println(fileName)
        }
	}
}
