package main

import (
	"os"
	"io"
	"fmt"
    "path"
	"crypto/sha256"
	"github.com/restic/chunker"
	// "compress/gzip"
	"archive/zip"
	// "github.com/vmihailenco/msgpack/v5"
)


type packTree struct {
	supersedesPackID string
	packIndexes []packIndex
}

type packIndex struct {
	ID string
	ChunkIDs []string
}


func PackFile(filePath string, repoPath string, mypoly int) ([]string, map[string]string) {
	f, _ := os.Open(filePath)
	chunkIDs, chunkPacks := WritePacks(f, repoPath, mypoly)
	f.Close()
	return chunkIDs, chunkPacks
}


// func PackFile(filePath string, repoPath string, mypoly int) ([]string, []packIndex) {
// 	f, _ := os.Open(filePath)
// 	chunkIDs, packIndexes := WritePacks(f, repoPath, mypoly)
// 	f.Close()
// 	return chunkIDs, packIndexes
// }


// func WritePacks(f *os.File, repoPath string, poly int) map[string]string {
func WritePacks(f *os.File, repoPath string, poly int) ([]string, map[string]string) {
	const maxPackSize uint = 104857600 // 100 MB
	mychunker := chunker.New(f, chunker.Pol(poly))
	buf := make([]byte, 8*1024*1024) // reuse this buffer
	chunkIDs := []string{}
	chunkPacks := ReadTrees(repoPath)
	newChunkPacks := make(map[string]string)	
	var curPackSize  uint 
	stillReadingInput := true

	for stillReadingInput {
		packId := RandHexString(PACK_ID_LEN)
		packFolderPath := path.Join(repoPath, "packs", packId[0:2])
		os.MkdirAll(packFolderPath, 0777)
		packPath := path.Join(packFolderPath, packId + ".zip")	
		fmt.Printf("\nCreating pack file %s\n", packPath)		

		zipFile, err := os.Create(packPath)
		check(err)
		zipWriter := zip.NewWriter(zipFile)

		i := 0
		curPackSize = 0

		for curPackSize < maxPackSize { // white chunks to pack
			chunk, err := mychunker.Next(buf)
			if err == io.EOF {
				fmt.Printf("Reached end of input file, stop chunking\n")
				stillReadingInput = false	
   				break
			}
	
			check(err)		
			
			i++
			chunkId := fmt.Sprintf("%064x", sha256.Sum256(chunk.Data))
			chunkIDs = append(chunkIDs, chunkId)
			curPackSize += chunk.Length

			if _, ok := chunkPacks[chunkId]; ok {
				fmt.Printf("Skipping Chunk ID %s already in pack %s\n", chunkId[0:16], chunkPacks[chunkId][0:16])
			} else {	
				fmt.Printf("Chunk %d: chunk size %d kB, total size %d kB, ", i, chunk.Length/1024, curPackSize/1024)
				fmt.Printf("chunk ID: %s\n",chunkId[0:16])
				chunkPacks[chunkId] = packId
				newChunkPacks[chunkId] = packId

				var header zip.FileHeader
				header.Name = chunkId
				header.Method = zip.Deflate
			
				writer, err := zipWriter.CreateHeader(&header)
				check(err)
				writer.Write(chunk.Data)	
			}		
		}	

		if stillReadingInput {
			fmt.Printf("Pack size %d exceeds max size %d\n", curPackSize, maxPackSize)		
		} else {
			fmt.Printf("Reached EOF of input\n")		
		}

		fmt.Printf("Closing zip file\n")
		zipWriter.Close()
		zipFile.Close()
	}

	return chunkIDs, newChunkPacks 
}


func UnpackFile(filePath string, repoPath string, chunkIds []string) {
	chunkPacks := ReadTrees(repoPath)

	f, _ := os.Create(filePath)
	ReadPacks(f, repoPath, chunkIds, chunkPacks)
	f.Close()
}


func ReadPacks(tarFile *os.File, repoPath string, chunkIds []string, chunkPacks map[string]string) {
	for i, chunkId := range chunkIds {
		packId := chunkPacks[chunkId]
		packPath := path.Join(repoPath, "packs", packId[0:2], packId + ".zip")
		fmt.Printf("Reading chunk %d %s \n from pack %s\n", i, chunkId, packPath)

		// From https://golangcode.com/unzip-files-in-go/
		r, err := zip.OpenReader(packPath)
		check(err)
	
		for _, f := range r.File {
			h := f.FileHeader
			if h.Name == chunkId {
				rc, err := f.Open()
				check(err)
				_, err = io.Copy(tarFile, rc)
				// fmt.Fprintf(tarFile, "Pack %s, chunk %s, csize %d, usize %d\n", packId, h.Name, h.CompressedSize, h.UncompressedSize)
				check(err)
				rc.Close()
			}
		}

		r.Close()			
	}
}

