package dupver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/restic/chunker"
)

const PackSize int64 = 500 * 1024 * 1024

func CreatePackFile(packId string) (*os.File, error) {
	packFolderPath := filepath.Join(".dupver", "packs", packId[0:2])
	os.MkdirAll(packFolderPath, 0777)
	packPath := filepath.Join(packFolderPath, packId+".zip")

	if VerboseMode {
		fmt.Fprintf(os.Stderr, "Creating pack: %s\n", packId[0:16])
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

func WriteChunkToPack(zipWriter *zip.Writer, chunkId string, chunk chunker.Chunk, compressionLevel uint16) error {
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

	if _, err := writer.Write(chunk.Data); err != nil {
		if VerboseMode {
			fmt.Fprintf(os.Stderr, "Error writing chunk %s to zip file\n", chunkId)
		}

		return err
	}

	return nil
}

func ExtractChunkFromPack(outFile *os.File, chunkId string, packId string) error {
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
	return ExtractChunkFromZipFile(outFile, packFile, chunkId)
}

func ExtractChunkFromZipFile(outFile *os.File, packFile *zip.ReadCloser, chunkId string) error {
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
