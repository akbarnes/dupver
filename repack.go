package dupver

import (
	"fmt"
	"os"
	"path/filepath"
)

func Repack() {
    // List all the snapshots
    snaps := ReadAllSnapshots()

    // If there is more than 1 snapshot
    if len(snaps) == 0 { 
        if VerboseMode {
            fmt.Fprintf(os.Stderr, "No snapshots, aborting\n")
        }

        return
    }

    // List all of the tree folders
	treesGlob := filepath.Join(".dupver", "trees*")
	treePaths, err := filepath.Glob(treesGlob)

    if err != nil {
        fmt.Fprintf(os.Stderr, "Error listing trees folders, aborting\n")
        return
    }

	existingPacks := ReadTrees()
	packs := map[string]string{}

    // Rename the old tree folder and create a new tree folder
    oldTreesPath := filepath.Join(".dupver", "trees")
    newTreesPath := filepath.Join(".dupver", fmt.Sprintf("trees%d", len(treePaths) - 1))

    if DebugMode {
        fmt.Fprintf(os.Stderr, "%s -> %s\n", oldTreesPath, newTreesPath)
    }

    if err := os.Rename(oldTreesPath, newTreesPath); err != nil { 
        fmt.Fprintf(os.Stderr, "Error renaming trees folder, aborting\n")
        return
    }

	packId := RandHexString(PackIdLen)
	packFile, err := CreatePackFile(packId)

	if err != nil {
		panic(fmt.Sprintf("Error creating pack file %s\n", packFile))
	}

	zipWriter := zip.NewWriter(packFile)
	var packBytesRemaining int64 = maxPackBytes

    // Repack for each snapshot
    for _, snap := range snaps {
	    snapFiles := snap.ReadFilesList()

		outPath := filepath.Join(outputFolder, fileName)
		outFile, err := os.Create(outPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s, skipping\n", outPath)
			continue
		}

		defer outFile.Close()

		for _, chunkId := range fileProps.ChunkIds {
			existingPackId := existingPacks[chunkId]

			if DebugMode {
				fmt.Fprintf(os.Stderr, "Extracting:\n  Pack %s\n  Chunk %s\n  to %s\n\n", existingPackId, chunkId, packId)
			}

		    packs[chunkId] = packId
		    RepackChunkToPack(zipWriter, chunkId, existingPackId, compressionLevel)
		    packBytesRemaining -= int64(chunk.Length)

			if packBytesRemaining <= 0 {
				if err := zipWriter.Close(); err != nil {
					// TODO: Should I return an error instead of quitting here? Is there anythig to do?
					panic(fmt.Sprintf("Error closing zipwriter for pack %s\n", packId))
				}

				if err := packFile.Close(); err != nil {
					// TODO: Should I return an error instead of quitting here? Is there anythig to do?
					panic(fmt.Sprintf("Error closing file for pack %s\n", packId))
				}

				packId = RandHexString(PackIdLen)
				packFile, err = CreatePackFile(packId)

				if err != nil {
					// TODO: Should I return an error instead of quitting here? Is there anythig to do?
					panic(fmt.Sprintf("Error creating pack file %s\n", packFile))
				}

				zipWriter = zip.NewWriter(packFile)
				packBytesRemaining = maxPackBytes
			}
        }
    }

    return
}

func RepackChunk(zipWriter *zip.Writer, chunkId string, packId string, compressionLevel uint16) error {
	packFolderPath := path.Join(".dupver", "packs", packId[0:2])
	packPath := path.Join(packFolderPath, packId+".zip")
	packFile, err := zip.OpenReader(packPath)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error extracting pack %s[%s]\n", packId, chunkId)
		}
		return err
	}

	defer packFile.Close()

	var header zip.FileHeader
	header.Name = chunkId
	header.Method = compressionLevel

	writer, err := zipWriter.CreateHeader(&header)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error creating zip header\n")
		}

		return err
	}


	return RepackChunkFromZipFile(outFile, packFile, chunkId, compressionLevel)
}

func RepackChunkFromZipFile(outFile *os.File, packFile *zip.ReadCloser, chunkId string, compressionLevel uint16) error {
	for _, f := range packFile.File {

		if f.Name == chunkId {
			// fmt.Fprintf(os.Stderr, "Contents of %s:\n", f.Name)
			chunkFile, err := f.Open()

			if err != nil {
				if VerboseMode {
					fmt.Fprintf(os.Stderr, "Error opening chunk %s\n", chunkId)
				}

				return err
			}

			_, err = io.Copy(outFile, chunkFile)

			if err != nil {
				if VerboseMode {
					fmt.Fprintf(os.Stderr, "Error reading chunk %s\n", chunkId)
				}

				return err
			}

			chunkFile.Close()
		}
	}

	return nil
}

