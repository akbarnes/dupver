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


func PackFile(filePath string, repoPath string, mypoly chunker.Pol, verbosity int) ([]string, map[string]string) {
	f, _ := os.Open(filePath)
	chunkIDs, chunkPacks := WritePacks(f, repoPath, mypoly, verbosity)
	f.Close()
	return chunkIDs, chunkPacks
}


// func WritePacks(f *os.File, repoPath string, poly int) map[string]string {
func WritePacks(f *os.File, repoPath string, poly chunker.Pol, verbosity int) ([]string, map[string]string) {
	const maxPackSize uint = 104857600 // 100 MB
	mychunker := chunker.New(f, chunker.Pol(poly))
	buf := make([]byte, 8*1024*1024) // reuse this buffer
	chunkIDs := []string{}
	chunkPacks := ReadTrees(repoPath)
	newChunkPacks := make(map[string]string)	
	var curPackSize  uint 
	stillReadingInput := true

	totalDataSize := 0
	totalPackNum := 0
	totalChunkNum := 0
	dupChunkNum := 0

	for stillReadingInput {
		packId := RandHexString(PACK_ID_LEN)
		packFolderPath := path.Join(repoPath, "packs", packId[0:2])
		os.MkdirAll(packFolderPath, 0777)
		packPath := path.Join(packFolderPath, packId + ".zip")	

		totalPackNum++		

		if verbosity >= 2 {
			fmt.Printf("Creating pack file %3d: %s\n", totalPackNum, packPath)	
		} else if verbosity == 1 {
			fmt.Printf("Creating pack number: %3d, ID: %s\n", totalPackNum, packId[0:16])	
		}

		zipFile, err := os.Create(packPath)
		check(err)
		zipWriter := zip.NewWriter(zipFile)

		i := 0
		curPackSize = 0


		for curPackSize < maxPackSize { // white chunks to pack
			chunk, err := mychunker.Next(buf)
			if err == io.EOF {
				// fmt.Printf("Reached end of input file, stop chunking\n")
				stillReadingInput = false	
   				break
			}
	
			check(err)		
			
			i++
			chunkId := fmt.Sprintf("%064x", sha256.Sum256(chunk.Data))
			chunkIDs = append(chunkIDs, chunkId)
			curPackSize += chunk.Length

			totalDataSize += int(chunk.Length)
			totalChunkNum++

			if _, ok := chunkPacks[chunkId]; ok {
				if verbosity >= 2 {
					fmt.Printf("Skipping Chunk ID %s already in pack %s\n", chunkId[0:16], chunkPacks[chunkId][0:16])
				}

				dupChunkNum++
			} else {	
				if verbosity >= 2 {
					fmt.Printf("Chunk %d: chunk size %d kB, total size %d kB, ", i, chunk.Length/1024, curPackSize/1024)
					fmt.Printf("chunk ID: %s\n",chunkId[0:16])
				}
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


		if verbosity >= 2 {
			if stillReadingInput {
				fmt.Printf("Pack size %d exceeds max size %d\n", curPackSize, maxPackSize)		
			}

			fmt.Printf("Reached end of input, closing zip file\n")
		}

		zipWriter.Close()
		zipFile.Close()
	}

	if verbosity >= 1 {
		fmt.Printf("%0.2f MB raw data stored\n", float64(totalDataSize)/1e6)
		fmt.Printf("%d total chunks, %d duplicate chunks\n", totalChunkNum, dupChunkNum)
		fmt.Printf("%d packs stored, %0.2f chunks/pack\n", totalPackNum, float64(totalChunkNum)/float64(totalPackNum))
	}

	return chunkIDs, newChunkPacks 
}


func UnpackFile(filePath string, repoPath string, chunkIds []string, verbosity int) {
	chunkPacks := ReadTrees(repoPath)

	f, _ := os.Create(filePath)
	ReadPacks(f, repoPath, chunkIds, chunkPacks, verbosity)
	f.Close()
}


func ReadPacks(tarFile *os.File, repoPath string, chunkIds []string, chunkPacks map[string]string, verbosity int) {
	for i, chunkId := range chunkIds {
		packId := chunkPacks[chunkId]
		packPath := path.Join(repoPath, "packs", packId[0:2], packId + ".zip")

		if verbosity >= 2 {
			fmt.Printf("Reading chunk %d %s \n from pack %s\n", i, chunkId, packPath)
		}

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

