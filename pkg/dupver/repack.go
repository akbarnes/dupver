package dupver

import (
	"fmt"
	"os"
    "io"
    "errors"
	"path/filepath"
    "archive/zip"
)

func Repack(maxPackBytes int64, compressionLevel uint16) error {
    // List all the snapshots
    snaps := ReadAllSnapshots()

    // If there is more than 1 snapshot
    if len(snaps) == 0 { 
        if VerboseMode {
            fmt.Fprintf(os.Stderr, "No snapshots, aborting\n")
        }

        return nil
    }

    // List all of the tree folders
	treesGlob := filepath.Join(WorkingDirectory, ".dupver", "trees*")
	treePaths, err := filepath.Glob(treesGlob)

    if err != nil {
        // fmt.Fprintf(os.Stderr, "Error listing trees folders, aborting\n")
        return errors.New("Error listing trees folders")
    }

	oldPacks := ReadTrees()
	newPacks := map[string]string{}

    // Rename the old tree folder and create a new tree folder
    oldTreesPath := filepath.Join(WorkingDirectory, ".dupver", "trees")
    newTreesPath := filepath.Join(WorkingDirectory, ".dupver", fmt.Sprintf("trees%d", len(treePaths) - 1))

    if DebugMode {
        fmt.Fprintf(os.Stderr, "%s -> %s\n", oldTreesPath, newTreesPath)
    }

    if err := os.Rename(oldTreesPath, newTreesPath); err != nil { 
        // fmt.Fprintf(os.Stderr, "Error renaming trees folder, aborting\n")
        return errors.New("Error renaming trees folder")
    }

    // List all of the packs folders
	packsGlob := filepath.Join(WorkingDirectory, ".dupver", "packs*")
	packPaths, err := filepath.Glob(packsGlob)

    if err != nil {
        // fmt.Fprintf(os.Stderr, "Error listing pack folders, aborting\n")
        return errors.New("Error listing pack folders")
    }

    // Rename the old tree folder and create a new tree folder
    oldPacksPath := filepath.Join(WorkingDirectory, ".dupver", "packs")
    newPacksPath := filepath.Join(WorkingDirectory, ".dupver", fmt.Sprintf("packs%d", len(packPaths) - 1))

    if DebugMode {
        fmt.Fprintf(os.Stderr, "%s -> %s\n", oldPacksPath, newPacksPath)
    }

    if err := os.Rename(oldPacksPath, newPacksPath); err != nil { 
        // fmt.Fprintf(os.Stderr, "Error renaming packs folder, aborting\n")
        return errors.New("Error renaming packs folder, aborting")
    }

    // TODO: move this ito a function
	packID := RandHexString(PackIdLen)
	packFile, err := CreatePackFile(packID)

	if err != nil {
		panic(fmt.Sprintf("Error creating pack file %s\n", packID))
	}

	zipWriter := zip.NewWriter(packFile)
	var packBytesRemaining int64 = maxPackBytes

    // Repack for each snapshot
    for _, snap := range snaps {
	    snapFiles := snap.ReadFilesList()

        // for each file in snapshot
	    for _, fileProps := range snapFiles {
            for _, chunkID := range fileProps.ChunkIds {
                oldPackId := oldPacks[chunkID]

                if DebugMode {
                    fmt.Fprintf(os.Stderr, "Extracting:\n  File: %s, Pack %s\n  Chunk %s\n  to %s\n\n", fileProps.Name, oldPackId, chunkID, packID)
                }

                if chunkSize, err := RepackChunk(zipWriter, chunkID, oldPackId, newPacksPath, compressionLevel); err == nil {
                    packBytesRemaining -= int64(chunkSize)
                } else {
                    panic(fmt.Sprintf("Error repacking chunk %v", err))
                }

                newPacks[chunkID] = packID

                if packBytesRemaining <= 0 {
                    if err := zipWriter.Close(); err != nil {
                        // TODO: Should I return an error instead of quitting here? Is there anythig to do?
                        panic(fmt.Sprintf("Error closing zipwriter for pack %s\n", packID))
                    }

                    if err := packFile.Close(); err != nil {
                        // TODO: Should I return an error instead of quitting here? Is there anythig to do?
                        panic(fmt.Sprintf("Error closing file for pack %s\n", packID))
                    }

                    // TODO: move this to a function
                    packID = RandHexString(PackIdLen)
                    packFile, err = CreatePackFile(packID)

                    if err != nil {
                        // TODO: Should I return an error instead of quitting here? Is there anythig to do?
                        panic(fmt.Sprintf("Error creating pack file %s\n", packID))
                    }

                    zipWriter = zip.NewWriter(packFile)
                    packBytesRemaining = maxPackBytes
                }
            }
        }
    }

	if err := zipWriter.Close(); err != nil {
		panic(fmt.Sprintf("Error closing zipwriter for pack %s\n", packID))
	}

	if err := packFile.Close(); err != nil {
		panic(fmt.Sprintf("Error closing file for pack %s\n", packID))
	}


	snaps[len(snaps)-1].WriteTree(newPacks)
    return nil
}

func RepackChunk(zipWriter *zip.Writer, chunkID string, packID string, oldPacksPath string, compressionLevel uint16) (uint64, error) {
	packFolderPath := filepath.Join(oldPacksPath, packID[0:2])
	packPath := filepath.Join(packFolderPath, packID+".zip")
	packFile, err := zip.OpenReader(packPath)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error extracting pack %s[%s]\n", packID, chunkID)
		}
		return 0, err
	}

	defer packFile.Close()

	var header zip.FileHeader
	header.Name = chunkID
	header.Method = compressionLevel

	writer, err := zipWriter.CreateHeader(&header)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error creating zip header\n")
		}

		return 0, err
	}

	for _, f := range packFile.File {
		if f.Name != chunkID {
            continue
        }

        // fmt.Fprintf(os.Stderr, "Contents of %s:\n", f.Name)
        chunkFile, err := f.Open()

        if err != nil {
            if VerboseMode {
                fmt.Fprintf(os.Stderr, "Error opening chunk %s\n", chunkID)
            }

            return 0, err
        }

        _, err = io.Copy(writer, chunkFile)

        if err != nil {
            if VerboseMode {
                fmt.Fprintf(os.Stderr, "Error reading chunk %s\n", chunkID)
            }

            return 0, err
        }

        chunkFile.Close()
        // TODO: use compressed size rather than uncompressed size
        // return f.CompressedSize64, nil
        return f.UncompressedSize64, nil
	}

	return 0, errors.New(fmt.Sprintf("Chunk %s not found in pack", chunkID))
}

