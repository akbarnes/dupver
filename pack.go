package main

import (
	"os"
	"io"
	"fmt"
    "path"
	"crypto/sha256"
	"github.com/restic/chunker"
	"compress/gzip"
	"archive/zip"
	// "github.com/vmihailenco/msgpack/v5"
)

type packIndex struct {
	ID string
	ChunkIDs []string
}


func ChunkFile(filePath string, repoPath string, mypoly int) []string {
	f, _ := os.Open(filePath)
	chunks := WriteChunks(f, repoPath, mypoly)
	f.Close()
	return chunks
}

func PackFile(filePath string, repoPath string, mypoly int) []packIndex {
	f, _ := os.Open(filePath)
	packIndexes := WritePacks(f, repoPath, mypoly)
	f.Close()
	return packIndexes
}



func WriteChunks(f *os.File, repoPath string, poly int) []string {
	const minPackSize int = 524288000

	// create a chunker
	chunks := []string{}
	mychunker := chunker.New(f, chunker.Pol(poly))

	// reuse this buffer
	buf := make([]byte, 8*1024*1024)
	
	i := 0

	for {
		chunk, err := mychunker.Next(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
		
		i++
		myHash := sha256.Sum256(chunk.Data)
		fmt.Printf("Chunk %d: %d kB, %02x\n", i, chunk.Length/1024, myHash)
		chunks = append(chunks, fmt.Sprintf("%02x", myHash))

		chunkFolder := fmt.Sprintf("%02x", myHash[0:1])
		chunkFolderPath := path.Join(repoPath, "packs", chunkFolder)
		os.MkdirAll(chunkFolderPath, 0777)

		chunkFilename := fmt.Sprintf("%064x.gz", myHash)
		chunkPath := path.Join(chunkFolderPath, chunkFilename)

		if _, err := os.Stat(chunkPath); err == nil {
			// path/to/whatever exists
			fmt.Printf("Duplicate chunk file %s exists\n", chunkPath)
		} else {
			g0, chunkPathErr := os.Create(chunkPath)
			check(chunkPathErr)
			g := gzip.NewWriter(g0)
			g.Write(chunk.Data)
			g.Close()
			g0.Close()
		}
	}


	return chunks
}


func WritePacks(f *os.File, repoPath string, poly int) []packIndex {
	const maxPackSize uint = 104857600 // 100 MB
	mychunker := chunker.New(f, chunker.Pol(poly))
	buf := make([]byte, 8*1024*1024) // reuse this buffer
	packIndexes := []packIndex{}
	chunkPack := make(map[string]string)	
	var curPackSize  uint 
	stillReadingInput := true

	for stillReadingInput {
		packId := RandHexString(PACK_ID_LEN)
		myPackIndex := packIndex{ID: packId}
		packFolderPath := path.Join(repoPath, "packs", packId[0:2])
		os.MkdirAll(packFolderPath, 0777)
		packPath := path.Join(packFolderPath, packId + ".zip")	
		fmt.Printf("Creating pack file %s\n", packPath)		

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
			myPackIndex.ChunkIDs = append(myPackIndex.ChunkIDs, chunkId)
			chunkPack[chunkId] = packId
			curPackSize += chunk.Length

			if _, ok := chunkPack[chunkId]; ok {
				//do something here
				fmt.Printf("Skipping Chunk ID %s\n  Already in pack %s\n",chunkId, chunkPack[chunkId])
			} else {	
				fmt.Printf("Chunk %d: chunk size %d kB, total size %d kB, ", i, chunk.Length/1024, curPackSize/1024)
				fmt.Printf("chunk ID: %s\n",chunkId)

				var header zip.FileHeader
				header.Name = chunkId
				header.Method = zip.Deflate
			
				writer, err := zipWriter.CreateHeader(&header)
				check(err)
				writer.Write(chunk.Data)	
			}		
		}	

		packIndexes = append(packIndexes, myPackIndex)
		if stillReadingInput {
			fmt.Printf("Pack size %d exceeds max size %d\n", curPackSize, maxPackSize)		
		} else {
			fmt.Printf("Reached EOF of input\n")		
		}

		fmt.Printf("Closing zip file\n")
		zipWriter.Close()
		zipFile.Close()
	}

	return packIndexes
}

func UnchunkFile(filePath string, repoPath string, chunks []string) {
	f, _ := os.Create(filePath)
	ReadChunks(f, repoPath, chunks)
	f.Close()
}


func UnpackFile(filePath string, repoPath string, packIndexes []packIndex) {
	f, _ := os.Create(filePath)
	// ReadChunks(f, repoPath, chunks)
	f.Close()
}


func ReadChunks(tarFile *os.File, repoPath string, chunks []string) {
	b := make([]byte, 1024)

	for i, hash := range chunks {
		chunkPath := path.Join(repoPath, "packs", hash[0:2], hash + ".gz")
		fmt.Printf("Reading %d %s\n", i, chunkPath)

		f0, err := os.Open(chunkPath)
		check(err)
		f, _ := gzip.NewReader(f0)

		for {
			n, _ := f.Read(b)
			tarFile.Write(b[0:n])

			if n == 0 {
				break
			}
		}

		f.Close()
		f0.Close()			
	}
}



