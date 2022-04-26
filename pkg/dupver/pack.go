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
	packFolderPath := filepath.Join(".dupver", "packs", packID[0:2])
	os.MkdirAll(packFolderPath, 0777)
	packPath := filepath.Join(packFolderPath, packID+".zip")

	if DebugMode {
		fmt.Fprintf(os.Stderr, "Creating pack: %s\n", packID[0:16])
	}

	// TODO: only create pack file if we need to save stuff - set to nil initially
	packFile, err := os.Create(packPath)

	if err != nil {
		fmt.Errorf("Error creating pack file %s: %w", packPath, err)
	}

	return packFile, nil
}

func WriteChunkToPack(zipWriter *zip.Writer, chunkID string, chunk chunker.Chunk, compressionLevel uint16) error {
	var header zip.FileHeader
	header.Name = chunkID
	header.Method = compressionLevel

	writer, err := zipWriter.CreateHeader(&header)

	if err != nil {
	    return fmt.Errorf("Error creating zip header: %w", err)
	}

	if _, err := writer.Write(chunk.Data); err != nil {
	    fmt.Errorf("Error writing chunk %s to zip file: %w", chunkID, err)
	}

	return nil
}

func ExtractChunkFromPack(outFile *os.File, chunkID string, packID string) error {
	packFolderPath := path.Join(".dupver", "packs", packID[0:2])
	packPath := path.Join(packFolderPath, packID+".zip")
	packFile, err := zip.OpenReader(packPath)

	if err != nil {
	    return fmt.Errorf("Error extracting pack %s[%s]: %w", packID, chunkID, err)
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
			    return fmt.Errorf("Error opening chunk %s: %w", chunkID, err)
			}

			_, err = io.Copy(outFile, chunkFile)

			if err != nil {
				return fmt.Errorf("Error reading chunk %s: %w", chunkID, err)
			}

			chunkFile.Close()
            return nil
		}
	}

	return errors.New(fmt.Sprintf("Couldn't find chunk %s in pack", chunkID))
}
