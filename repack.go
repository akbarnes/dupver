package dupver

import (
	"fmt"
	"os"
    "io"
    "errors"
	"path/filepath"
    "archive/zip"
)

func Repack(maxPackBytes int64, compressionLevel uint16) {
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

	oldPacks := ReadTrees()
	newPacks := map[string]string{}

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

    // TODO: move this ito a function
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

        // for each file in snapshot
	    for fileName, fileProps := range snapFiles {
            for _, chunkId := range fileProps.ChunkIds {
                oldPackId := oldPacks[chunkId]

                if DebugMode {
                    fmt.Fprintf(os.Stderr, "Extracting:\n  File: %s, Pack %s\n  Chunk %s\n  to %s\n\n", fileName, oldPackId, chunkId, packId)
                }

                if chunkSize, err := RepackChunk(zipWriter, chunkId, oldPackId, compressionLevel); err == nil {
                    packBytesRemaining -= int64(chunkSize)
                } else {
                    panic(fmt.Sprintf("Error repacking chunk %v", err))
                }

                newPacks[chunkId] = packId

                if packBytesRemaining <= 0 {
                    if err := zipWriter.Close(); err != nil {
                        // TODO: Should I return an error instead of quitting here? Is there anythig to do?
                        panic(fmt.Sprintf("Error closing zipwriter for pack %s\n", packId))
                    }

                    if err := packFile.Close(); err != nil {
                        // TODO: Should I return an error instead of quitting here? Is there anythig to do?
                        panic(fmt.Sprintf("Error closing file for pack %s\n", packId))
                    }

                    // TODO: move this to a function
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
    }

	snaps[len(snaps)-1].WriteTree(newPacks)
    return
}

func RepackChunk(zipWriter *zip.Writer, chunkId string, packId string, compressionLevel uint16) (uint64, error) {
	packFolderPath := filepath.Join(".dupver", "packs", packId[0:2])
	packPath := filepath.Join(packFolderPath, packId+".zip")
	packFile, err := zip.OpenReader(packPath)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error extracting pack %s[%s]\n", packId, chunkId)
		}
		return 0, err
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

		return 0, err
	}

	return ReZip(writer, packFile, chunkId, compressionLevel)
}

func ReZip(outFile *io.Writer, packFile *zip.ReadCloser, chunkId string, compressionLevel uint16) (uint64, error) {
	for _, f := range packFile.File {
		if f.Name != chunkId {
            continue
        }

        // fmt.Fprintf(os.Stderr, "Contents of %s:\n", f.Name)
        chunkFile, err := f.Open()

        if err != nil {
            if VerboseMode {
                fmt.Fprintf(os.Stderr, "Error opening chunk %s\n", chunkId)
            }

            return 0, err
        }

        _, err = io.Copy(outFile, chunkFile)

        if err != nil {
            if VerboseMode {
                fmt.Fprintf(os.Stderr, "Error reading chunk %s\n", chunkId)
            }

            return 0, err
        }

        chunkFile.Close()
        // TODO: use compressed size rather than uncompressed size
        // return f.CompressedSize64, nil
        return f.UncompressedSize64, nil
	}

	return 0, errors.New(fmt.Sprintf("Chunk %s not found in pack", chunkId))
}

