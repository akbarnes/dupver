package dupver

import (
	"fmt"
	"os"
    "time"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
)

// CheckoutSnapshot extracts a snapshot to the working directory
// given its commit ID
// It takes an optional output folder and pattern to match files
// to support partial checkouts
func CheckoutSnapshot(commitID string, outputFolder string, filter string, archiveTool string) {
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


	fmt.Fprintf(os.Stderr, "Checking out %s\n", snap.SnapshotID[0:9])
    snap.Checkout(outputFolder, filter, archiveTool)
}

// snapshot.Checkout extracts a snapshot to the working directory
// It takes an optional output folder and pattern to match files
// to support partial checkouts
func (snap Snapshot) Checkout(outputFolder string, filter string, archiveTool string) {
	os.MkdirAll(outputFolder, 0777)
	snapFiles := snap.ReadFilesHash()
	packs := ReadTrees()

	for fileName, fileProps := range snapFiles {
        fileName = ToForwardSlashes(fileName)
        matched, err := doublestar.Match(filter, fileName)

        if err != nil && VerboseMode {
            fmt.Fprintf(os.Stderr, "Error matching %s\n", filter)
        }

        if !matched {
            if DebugMode {
                fmt.Fprintf(os.Stderr, "Skipping file %s\n", fileName)
            }

            continue
        }

        nativeFileName := ToNativeSeparators(fileName)
		fileDir := filepath.Dir(nativeFileName)
		outDir := ToNativeSeparators(outputFolder)

		if fileDir != "." {
			outDir = filepath.Join(outputFolder, fileDir)
            fileInfo, err := os.Stat(outDir)

            if os.IsNotExist(err) || !fileInfo.IsDir() {
                fmt.Fprintf(os.Stderr, "Creating folder %s\n", outDir)
                os.MkdirAll(outDir, 0777)
            }

		}

		outPath := filepath.Join(outputFolder, nativeFileName)
        archivePath := outPath
        archiveBaseName := ""

        if fileProps.IsArchive {
            archiveBaseName = GenArchiveBaseName()
            archivePath, err = GenTempArchivePath(archiveBaseName)

            if err != nil {
                fmt.Fprintf(os.Stderr, "Error creating output archive path, skipping: %v\n", err)
                continue
            }
        }

		outFile, err := os.Create(archivePath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s, skipping\n", outPath)
			continue
		}

		defer outFile.Close()

		for _, chunkID := range fileProps.ChunkIds {
			packID := packs[chunkID]

			if DebugMode {
				fmt.Fprintf(os.Stderr, "Extracting:\n  Pack %s\n  Chunk %s\n  to %s\n\n", packID, chunkID, outPath)
			}

			if err := ExtractChunkFromPack(outFile, chunkID, packID); err != nil {
                fmt.Fprintf(os.Stderr, "Error extracting:\n  chunk: %s\n  pack: %s\n\n", chunkID, packID)
            }
		}

        if fileProps.IsArchive {
            if err := PostprocessArchive(archiveBaseName, outPath, archiveTool); err != nil {
                fmt.Fprintf(os.Stderr, "Error postprocessing %s to %s\n: %v\n", archivePath, outPath, err)
                continue
            }
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
