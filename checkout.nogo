package dupver

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckoutSnaphot(snapshotNum int, outputFolder string) {
	dupverDir := filepath.Join(WorkingDirectory, ".dupver")

	if len(outputFolder) == 0 {
		outputFolder = fmt.Sprintf("snapshot%03d", snapshotNum)
	}

	fmt.Printf("Checking out %s\n", snapshotNum)

	snapshotGlob := filepath.Join(dupverDir, "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)

	snapshotPath := snapshotPaths[snapshotNum-1]
	fmt.Printf("Reading %s\n", snapshotPath)
	snap := ReadSnapshotFile(snapshotPath)

	os.Mkdir(outputFolder, 0777)

	for file, _ := range snap.FileModTimes {
		fileDir := filepath.Dir(file)
		outDir := outputFolder

		if fileDir != "." {
			outDir = filepath.Join(outputFolder, fileDir)
			fmt.Printf("Creating folder %s\n", outDir)
			os.MkdirAll(outDir, 0777)
		}

		outPath := filepath.Join(outputFolder, file)
		outFile, err := os.Create(outPath)

		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error creating %s, skipping\n", outPath)
			fmt.Printf("Error creating %s, skipping\n", outPath)
			continue
		}

		defer outFile.Close()

		for _, chunkId := range snap.FileChunkIds[file] {
			packId := snap.ChunkPackIds[chunkId]
			ExtractChunkFromPack(outFile, chunkId, packId)
		}

		fmt.Printf("Restored %s to %s\n", file, outPath)
	}
}
