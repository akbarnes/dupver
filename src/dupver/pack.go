package dupver

import (
	"os"
	"io"
	"io/ioutil"
	"fmt"
    "path"
	"crypto/sha256"
	"github.com/restic/chunker"
	// "compress/gzip"
	"archive/zip"
	// "github.com/vmihailenco/msgpack/v5"
	"log"

	"github.com/akbarnes/dupver/src/fancyprint"
)


type packTree struct {
	supersedesPackID string
	packIndexes []packIndex
}

type packIndex struct {
	ID string
	ChunkIDs []string
}

// Pack a file to the repository given a file path
func PackFile(filePath string, repoPath string, mypoly chunker.Pol) ([]string, map[string]string) {
	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not open file when packing %s", filePath))
	}
	chunkIDs, chunkPacks := WritePacks(f, repoPath, mypoly)
	f.Close()
	return chunkIDs, chunkPacks
}

// Pack a file stream to the repository
// func WritePacks(f *os.File, repoPath string, poly int) map[string]string {
func WritePacks(f *os.File, repoPath string, poly chunker.Pol) ([]string, map[string]string) {
	const maxPackSize uint = 104857600 // 100 MB
	mychunker := chunker.New(f, chunker.Pol(poly))
	buf := make([]byte, 8*1024*1024) // reuse this buffer
	chunkIDs := []string{}
	chunkPacks := ReadTrees(repoPath)
	newChunkPacks := make(map[string]string)	
	var curPackSize  uint 
	stillReadingInput := true

	totalDataSize := 0
	dupDataSize := 0

	newPackNum := 0
	totalChunkNum := 0
	dupChunkNum := 0

	for stillReadingInput {
		packId := RandHexString(PACK_ID_LEN)
		packFolderPath := path.Join(repoPath, "packs", packId[0:2])
		os.MkdirAll(packFolderPath, 0777)
		packPath := path.Join(packFolderPath, packId + ".zip")	

		newPackNum++		

		fancyprint.Debugf("Creating pack file %3d: %s\n", newPackNum, packPath)	
		
		if fancyprint.Verbosity <= fancyprint.NoticeLevel {
			fmt.Printf("Creating pack number: %3d, ID: %s\n", newPackNum, packId[0:16])	
		}

		zipFile, err := os.Create(packPath)
		
		if err != nil {
			panic(fmt.Sprintf("Error creating zip file %s", packPath))
		}
		zipWriter := zip.NewWriter(zipFile)

		i := 0
		curPackSize = 0


		for curPackSize < maxPackSize { // white chunks to pack
			chunk, err := mychunker.Next(buf)
			if err == io.EOF {
				// fmt.Printf("Reached end of input file, stop chunking\n")
				stillReadingInput = false	
   				break
			} else if err != nil {
				panic("Error chunking input file")
			}
	
			i++
			chunkId := fmt.Sprintf("%064x", sha256.Sum256(chunk.Data))
			chunkIDs = append(chunkIDs, chunkId)

			totalDataSize += int(chunk.Length)
			totalChunkNum++

			if _, ok := chunkPacks[chunkId]; ok {
				fancyprint.Infof("Skipping Chunk ID %s already in pack %s\n", chunkId[0:16], chunkPacks[chunkId][0:16])
				dupChunkNum++
				dupDataSize += int(chunk.Length)
			} else {	
				fancyprint.Infof("Chunk %d: chunk size %d kB, total size %d kB, ", i, chunk.Length/1024, curPackSize/1024)
				fancyprint.Infof("chunk ID: %s\n",chunkId[0:16])
				chunkPacks[chunkId] = packId
				newChunkPacks[chunkId] = packId

				var header zip.FileHeader
				header.Name = chunkId
				header.Method = zip.Deflate
			
				writer, err := zipWriter.CreateHeader(&header)
				
				if err != nil {
					panic(fmt.Sprintf("Error creating zip file header for %s", packPath))
				}

				writer.Write(chunk.Data)	
				curPackSize += chunk.Length
			}		
		}	


		if stillReadingInput {
			fancyprint.Info("Pack size %d exceeds max size %d\n", curPackSize, maxPackSize)		
		}

		fancyprint.Info("Reached end of input, closing zip file\n")
		zipWriter.Close()
		zipFile.Close()
	}

	newChunkNum := totalChunkNum - dupChunkNum
	newDataSize := totalDataSize - dupDataSize

	newMb := float64(newDataSize)/1e6
	dupMb := float64(dupDataSize)/1e6
	totalMb := float64(totalDataSize)/1e6

	fancyprint.Noticef("%0.2f new, %0.2f duplicate, %0.2f total MB raw data stored\n", newMb, dupMb, totalMb)
	fancyprint.Noticef("%d new, %d duplicate, %d total chunks\n", newChunkNum, dupChunkNum, totalChunkNum)
	fancyprint.Noticef("%d packs stored, %0.2f chunks/pack\n", newPackNum, float64(newChunkNum)/float64(newPackNum))

	return chunkIDs, newChunkPacks 
}

// Unpack a file from the repository to a specified file path
func UnpackFile(filePath string, repoPath string, chunkIds []string, opts Options) {
	chunkPacks := ReadTrees(repoPath)

	f, err := os.Create(filePath)

	if err != nil {
		panic(fmt.Sprintf("Could not create output file %s while unpacking", filePath))
	}

	ReadPacks(f, repoPath, chunkIds, chunkPacks, opts)
	f.Close()
}

// Unpack a file from the repository to an output stream
// TODO: change name to something other than read (UnpackData?)
func ReadPacks(tarFile *os.File, repoPath string, chunkIds []string, chunkPacks map[string]string, opts Options) {
	for i, chunkId := range chunkIds {
		packId := chunkPacks[chunkId]
		packPath := path.Join(repoPath, "packs", packId[0:2], packId + ".zip")
		fancyprint.Infof("Reading chunk %d %s \n from pack %s\n", i, chunkId, packPath)

		// From https://golangcode.com/unzip-files-in-go/
		r, err := zip.OpenReader(packPath)
		
		if err != nil {
			panic(fmt.Sprintf("Error opening zip file %s", packPath))
		}
	
		for _, f := range r.File {
			h := f.FileHeader
			if h.Name == chunkId {
				rc, err := f.Open()
				
				if err != nil {
					panic(fmt.Sprintf("Error opening pack/chunk %s/%s", packPath, h.Name))
				}

				if _, err := io.Copy(tarFile, rc); err != nil {
					// fmt.Fprintf(tarFile, "Pack %s, chunk %s, csize %d, usize %d\n", packId, h.Name, h.CompressedSize, h.UncompressedSize)
					panic(fmt.Sprintf("Error reading from pack/chunk %s/%s", packPath, h.Name))
				}

				rc.Close()
			}
		}

		r.Close()			
	}
}

// Read a chunk from a specified repository given a chunk ID and chunk-pack mapping
func LoadChunk(repoPath string, chunkId string, chunkPacks map[string]string, opts Options) []byte {
	packId := chunkPacks[chunkId]

	// if opts.Verbosity >= 2 {
	// 	fmt.Printf("Reading chunk %s \n from pack %s\n", chunkId, packId)
	// }

	packPath := path.Join(repoPath, "packs", packId[0:2], packId + ".zip")
	data := []byte{}

	fancyprint.Infof("Reading chunk %s \n from pack file %s\n", chunkId, packPath)

	// From https://golangcode.com/unzip-files-in-go/
	packReader, err := zip.OpenReader(packPath)
	
	if err != nil {
		panic(fmt.Sprintf("Error opening zip file %s", packPath))
	}

	for _, f := range packReader.File {
		h := f.FileHeader
		if h.Name == chunkId {
			chunkReader, err := f.Open()
			
			if err != nil {
				// TODO: return err
				panic(fmt.Sprintf("Error opening chunk %s from pack file", h.Name, packPath))
			}

			{
				var err error
				// if _, err := io.Copy(tarFile, rc); err != nil {

				if data, err = ioutil.ReadAll(chunkReader); err != nil {
					// fmt.Fprintf(tarFile, "Pack %s, chunk %s, csize %d, usize %d\n", packId, h.Name, h.CompressedSize, h.UncompressedSize)
					// TODO: return err
					panic(fmt.Sprintf("Error opening chunk %s from pack file", h.Name, packPath))
				}
			}

			chunkReader.Close()
		}
	}

	packReader.Close()		
	return data	
}

