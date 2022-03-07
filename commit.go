package dupver

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/restic/chunker"
)

// CreateZipFile creates a pack zip file with a random 
// pack ID given a specified pack ID length.
// It returns the pack id, file descriptor, and zip writer
func CreateZipFile(packIDLen int) (string, *os.File, *zip.Writer) {
	packID := RandHexString(packIDLen)
	packFile, err := CreatePackFile(packID)

	if err != nil {
		panic(fmt.Sprintf("Error creating pack file %s\n", packID))
	}

	zipWriter := zip.NewWriter(packFile)
    return packID, packFile, zipWriter
}


func CommitSnapshot(message string, filters []string, archiveTypes []string, archiveTool string, poly chunker.Pol, maxPackBytes int64, compressionLevel uint16) {
	buf := make([]byte, 8*1024*1024) // reuse this buffer

	if DebugMode {
		fmt.Fprintf(os.Stderr, "Start reading head...\n")
	}

	headFiles := ReadHead().ReadFilesHash()

	if DebugMode {
		fmt.Fprintf(os.Stderr, "Done reading head\n")
		fmt.Fprintf(os.Stderr, "Start reading trees...\n")
	}

	existingPacks := ReadTrees()

	if DebugMode {
		fmt.Fprintf(os.Stderr, "Done reading trees\n")
	}

	snap := CreateSnapshot(message)
	files := []SnapshotFile{}
	packs := map[string]string{}

	dupverDir := filepath.Join(".dupver")

	if err := os.MkdirAll(dupverDir, 0777); err != nil {
		panic(fmt.Sprintf("Error creating dupver folder %s\n", dupverDir))
	}

    packID, packFile, zipWriter := CreateZipFile(PackIdLen)
	var packBytesRemaining int64 = maxPackBytes
    committedFilesCount := 0

	var VersionFile = func(fileName string, info os.FileInfo, err error) error {
		fileName = ToForwardSlashes(strings.TrimSuffix(fileName, "\n"))

		if ExcludedFile(fileName, info, filters) {
			return nil
		}

		props, err := os.Stat(fileName)

		if err != nil {
			if VerboseMode {
				fmt.Fprintf(os.Stderr, "Can't stat file %s, skipping\n", fileName)
			}

			return nil
		}

		modTime := props.ModTime().UTC().Format("2006-01-02T15-04-05")
		modLocalTime := props.ModTime().Format("2006-01-02T15-04-05")
		file := SnapshotFile{Name: fileName, ModTime: modTime, ModLocalTime: modLocalTime, Size: props.Size(), IsArchive: false}
		file.ChunkIds = []string{}

		// TODO: fix this. Currently not reading in filechunks from head
		if headFile, ok := headFiles[fileName]; ok && modTime == headFile.ModTime {
			if DebugMode {
				fmt.Fprintf(os.Stderr, "Skipping %s\n", fileName)
			}

			files = append(files, headFiles[fileName])
			// snap.AddFileChunkIds(head, fileName)
			return nil
		}

        archiveFileName := fileName

        if ArchiveFile(fileName, info, archiveTypes) {
            // 7z x -oTempFolder fileName
            // 7z a TempFile.7z TempFolder + FileSep + *
            archiveFileName, err = PreprocessArchive(fileName, archiveTool)

            if err != nil {
                if !QuietMode {
                    fmt.Fprintf(os.Stderr, "Error preprocessing archive %s, skipping: %v\n", fileName, err)
                }

                return nil
            }

            file.IsArchive = true
        }

		in, err := os.Open(archiveFileName)
        committedFilesCount++

		if err != nil {
			if !QuietMode {
				fmt.Fprintf(os.Stderr, "Can't open file %s for reading, skipping\n", archiveFileName)
			}

			return nil
		}

		defer in.Close()
		myChunker := chunker.New(in, chunker.Pol(poly))

		if QuietMode {
			fmt.Println(fileName)
		} else {
            if file.IsArchive {
			    fmt.Fprintf(os.Stderr, "Archive %s\n", fileName)
            } else {
			    fmt.Fprintf(os.Stderr, "Regular %s\n", fileName)
            }
		}

		readingFile := true
        readOk := true

		for readingFile {
			chunk, err := myChunker.Next(buf)

			if err == io.EOF {
				readingFile = false
			} else if err != nil {
				// TODO: Should I return an error instead of quitting here? Is there anythig to do?
                if VerboseMode {
				    fmt.Fprintf(os.Stderr, "Error reading chunk from  file %s\n", fileName)
                }

                readOk = false
                break
			}

			chunkID := fmt.Sprintf("%064x", sha256.Sum256(chunk.Data))

			if chunk.Length > 0 {
				file.ChunkIds = append(file.ChunkIds, chunkID)

				if _, ok := existingPacks[chunkID]; ok {
					if DebugMode {
						fmt.Fprintf(os.Stderr, "Skipping Chunk ID %s already in pack %s\n", chunkID[0:16], existingPacks[chunkID][0:16])
					}

					continue
				}

				if DebugMode {
					fmt.Fprintf(os.Stderr, "Chunk %s: chunk size %d kB, pack %s\n", chunkID[0:16], chunk.Length/1024, packID[0:16])
				}

				packs[chunkID] = packID
				existingPacks[chunkID] = packID

				// save zip data
				WriteChunkToPack(zipWriter, chunkID, chunk, compressionLevel)
				packBytesRemaining -= int64(chunk.Length)
			}

			if packBytesRemaining <= 0 {
				if err := zipWriter.Close(); err != nil {
					// TODO: Should I return an error instead of quitting here? Is there anythig to do?
					panic(fmt.Sprintf("Error closing zipwriter for pack %s\n", packID))
				}

				if err := packFile.Close(); err != nil {
					// TODO: Should I return an error instead of quitting here? Is there anythig to do?
					panic(fmt.Sprintf("Error closing file for pack %s\n", packID))
				}

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

        if readOk {
		    files = append(files, file)
        }

		return nil
	}

	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)
	if err := filepath.Walk(WorkingDirectory, VersionFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error committing:\n")
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	if err := zipWriter.Close(); err != nil {
		panic(fmt.Sprintf("Error closing zipwriter for pack %s\n", packID))
	}

	if err := packFile.Close(); err != nil {
		panic(fmt.Sprintf("Error closing file for pack %s\n", packID))
	}

    if committedFilesCount > 0 || ForceMode {
        snap.Write()
        snap.WriteFiles(files)
        snap.WriteTree(packs)
        snap.WriteHead()
    } else {
        fmt.Fprintf(os.Stderr, "No modified files, skipping commit")
    }
}
