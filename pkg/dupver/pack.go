package dupver

import (
	"fmt"
	"io"
	"os"
	"path"
    "errors"
	"path/filepath"
	"archive/zip"

	"github.com/restic/chunker"
)

const PackSize int64 = 500 * 1024 * 1024

func CreatePackFile(packID string) (*os.File, error) {
	packFolderPath := filepath.Join(WorkingDirectory, ".dupver", "packs", packID[0:2])
	os.MkdirAll(packFolderPath, 0777)
	packPath := filepath.Join(packFolderPath, packID+".zip")

	if DebugMode {
		fmt.Fprintf(os.Stderr, "Creating pack: %s\n", packID[0:16])
	}

	// TODO: only create pack file if we need to save stuff - set to nil initially
	packFile, err := os.Create(packPath)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error creating pack file %s", packPath)
		}

		return nil, err
	}

	return packFile, nil
}

func WriteChunkToPack(zipWriter *zip.Writer, chunkID string, chunk chunker.Chunk, compressionLevel uint16) error {
	var header zip.FileHeader
	header.Name = chunkID
	header.Method = compressionLevel

	writer, err := zipWriter.CreateHeader(&header)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error creating zip header\n")
		}

		return err
	}

	if _, err := writer.Write(chunk.Data); err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error writing chunk %s to zip file\n", chunkID)
		}

		return err
	}

	return nil
}

func ExtractChunkFromPack(outFile *os.File, chunkID string, packID string) error {
	packFolderPath := path.Join(WorkingDirectory, ".dupver", "packs", packID[0:2])
	packPath := path.Join(packFolderPath, packID+".zip")
	packFile, err := zip.OpenReader(packPath)

	if err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error extracting pack %s[%s]\n", packID, chunkID)
		}
		return err
	}

	defer packFile.Close()
	return ExtractChunkFromZipFile(outFile, packFile, chunkID)
}

func ExtractChunkFromZipFile(outFile *os.File, packFile *zip.ReadCloser, chunkID string) error {
	for _, f := range packFile.File {

		if f.Name == chunkID {
			// fmt.Fprintf(os.Stderr, "Contents of %s:\n", f.Name)
			chunkFile, err := f.Open()

			if err != nil {
				if VerboseMode {
					fmt.Fprintf(os.Stderr, "Error opening chunk %s\n", chunkID)
				}

				return err
			}

			_, err = io.Copy(outFile, chunkFile)

			if err != nil {
				if VerboseMode {
					fmt.Fprintf(os.Stderr, "Error reading chunk %s\n", chunkID)
				}

				return err
			}

			chunkFile.Close()
            return nil
		}
	}

	return errors.New(fmt.Sprintf("Couldn't find chunk %s in pack", chunkID))
}
