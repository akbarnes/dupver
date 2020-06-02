package main

import (
	"os"
	"io"
	"fmt"
    "path"
	"crypto/sha256"
	"github.com/restic/chunker"
	"compress/gzip"
)

func PackFile(filePath string, repoPath string, commitFile *os.File, mypoly int) {
	f, _ := os.Open(filePath)
	WritePacks(f, repoPath, commitFile, mypoly)
	f.Close()
}


func WritePacks(f *os.File, repoPath string, commitFile *os.File, poly int) {
	// create a chunker
	mychunker := chunker.New(f, chunker.Pol(poly))

	// reuse this buffer
	buf := make([]byte, 8*1024*1024)
	
	fmt.Fprintf(commitFile, "chunks = [\n")

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
		fmt.Fprintf(commitFile, "  \"%02x\",\n", myHash)

		chunkFolder := fmt.Sprintf("%02x", myHash[0:1])
		chunkFolderPath := path.Join(repoPath, "packs", chunkFolder)
		os.MkdirAll(chunkFolderPath, 0777)

		chunkFilename := fmt.Sprintf("%064x.gz", myHash)
		chunkPath := path.Join(chunkFolderPath, chunkFilename)
		g0, chunkPathErr := os.Create(chunkPath)
        check(chunkPathErr)
		g := gzip.NewWriter(g0)
		g.Write(chunk.Data)
		g.Close()
		g0.Close()
	}


	fmt.Fprintf(commitFile, "]\n\n")
}



func UnpackTar(filePath string, chunks []string) {
	tarFile, _ := os.Create(filePath)
	WriteChunks(tarFile, chunks)
	tarFile.Close()
}


func WriteChunks(tarFile *os.File, chunks []string) {
	b := make([]byte, 1024)

	for i, hash := range chunks {
		chunkPath := fmt.Sprintf(".dupver/%s/%s.gz", hash[0:2], hash)
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



