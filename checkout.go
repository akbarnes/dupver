package dupver

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
)

func CheckoutSnapshot(commitId string, outputFolder string, filter string) {
    snap, err := MatchSnapshot(commitId)

    if err != nil {
        fmt.Println("No matching snapshot paths")
        os.Exit(1)
    }

	fmt.Printf("Checking out %s\n", snap.SnapshotId[0:9])
    snap.Checkout(outputFolder, filter)
}

func (snap Snapshot) Checkout(outputFolder string, filter string) {
	os.MkdirAll(outputFolder, 0777)
	snapFiles := snap.ReadFilesList()
	packs := ReadTrees()

	for fileName, fileProps := range snapFiles {
        matched, err := doublestar.PathMatch(filter, fileName)

        if err != nil && VerboseMode {
            fmt.Printf("Error matching %s\n", filter)
        }

        if !matched {
            if VerboseMode {
                fmt.Printf("Skipping file %s\n", fileName)
            }

            continue
        }

		fileDir := filepath.Dir(fileName)
		outDir := outputFolder

		if fileDir != "." {
			outDir = filepath.Join(outputFolder, fileDir)
			fmt.Printf("Creating folder %s\n", outDir)
			os.MkdirAll(outDir, 0777)
		}

		outPath := filepath.Join(outputFolder, fileName)
		outFile, err := os.Create(outPath)

		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error creating %s, skipping\n", outPath)
			fmt.Printf("Error creating %s, skipping\n", outPath)
			continue
		}

		defer outFile.Close()

		for _, chunkId := range fileProps.ChunkIds {
			packId := packs[chunkId]

			if VerboseMode {
				fmt.Printf("Extracting:\n  Pack %s\n  Chunk %s\n  to %s\n\n", packId, chunkId, outPath)
			}

			ExtractChunkFromPack(outFile, chunkId, packId)
		}

		fmt.Printf("Restored %s to %s\n", fileName, outPath)
	}
}
