package main

import (
	"os"
	"io"
	"fmt"
    "path"
	"crypto/sha256"
	"github.com/restic/chunker"
	"compress/gzip"
	// "archive/zip"
	// "github.com/vmihailenco/msgpack/v5"
)


func ChunkFile(filePath string, repoPath string, mypoly int) []string {
	f, _ := os.Open(filePath)
	chunks := WriteChunks(f, repoPath, mypoly)
	f.Close()
	return chunks
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


func UnchunkFile(filePath string, repoPath string, chunks []string) {
	f, _ := os.Create(filePath)
	ReadChunks(f, repoPath, chunks)
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