package main

import (
	"os"
	"io"
	"fmt"
	"crypto/sha256"
	"github.com/restic/chunker"
	"compress/gzip"
)

func PackTar(filePath string, commitFile *os.File, mypoly int) {
	tarFile, _ := os.Open(filePath)
	WritePacks(tarFile, commitFile, mypoly)
	tarFile.Close()
}


func PackTGZ(filePath string, commitFile *os.File, mypoly int) {
	f0, _ := os.Open(filePath)
	f, _ := gzip.NewReader(f0)
	WritePacksFromGZ(f, commitFile, mypoly)
	f.Close()
	f0.Close()
}


func WritePacks(tarFile *os.File, commitFile *os.File, poly int) {
	// create a chunker
	mychunker := chunker.New(tarFile, chunker.Pol(poly))

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

		chunkFolder := fmt.Sprintf(".dupver/%02x", myHash[0:1])
		os.MkdirAll(chunkFolder, 0777)

		chunkPath := fmt.Sprintf("%s/%02x.gz", chunkFolder, myHash)
		g0, _ := os.Create(chunkPath)
		g := gzip.NewWriter(g0)
		g.Write(chunk.Data)
		g.Close()
		g0.Close()
	}


	fmt.Fprintf(commitFile, "]\n\n")
}


func WritePacksFromGZ(f *gzip.Reader, commitFile *os.File, poly int) {
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

		chunkFolder := fmt.Sprintf(".dupver/%02x", myHash[0:1])
		os.MkdirAll(chunkFolder, 0777)

		chunkPath := fmt.Sprintf("%s/%02x.gz", chunkFolder, myHash)
		g0, _ := os.Create(chunkPath)
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


func UnpackTGZ(filePath string, chunks []string) {
	g0, _ := os.Create(filePath)
	g := gzip.NewWriter(g0)
	WriteGZChunks(g, chunks)
	g.Close()
	g0.Close()
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


func WriteGZChunks(g *gzip.Writer, chunks []string) {
	b := make([]byte, 1024)

	for i, hash := range chunks {
		chunkPath := fmt.Sprintf(".dupver/%s/%s.gz", hash[0:2], hash)
		fmt.Printf("Reading %d %s\n", i, chunkPath)

		f0, err := os.Open(chunkPath)
		check(err)
		f, _ := gzip.NewReader(f0)

		for {
			n, _ := f.Read(b)
			g.Write(b[0:n])

			if n == 0 {
				break
			}
		}

		f.Close()
		f0.Close()			
	}
}
