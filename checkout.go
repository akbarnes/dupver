package dupver

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckoutSnapshot(commitId string, outputFolder string) {
    snap, err := MatchSnapshot(commitId)

    if err != nil {
        fmt.Println("No matching snapshot paths")
        os.Exit(1)
    }


	fmt.Printf("Checking out %d\n", snap.SnapshotId[0:9])
	snap.Checkout(outputFolder)
}

func (snap Snapshot) Checkout(outputFolder string) {
	os.Mkdir(outputFolder, 0777)
	snapFiles := snap.ReadFilesList()
	packs := ReadTrees()

	for fileName, fileProps := range snapFiles {
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
