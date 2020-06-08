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


type pack struct {
	Chunks [][]byte
	Hashes []string
}

func ChunkFile(filePath string, repoPath string, mypoly int) []string {
	f, _ := os.Open(filePath)
	chunks := WriteChunks(f, repoPath, mypoly)
	f.Close()
	return chunks
}

func PackFile(filePath string, repoPath string, mypoly int) ([]string, []string) {
	f, _ := os.Open(filePath)
	packs, chunks := WritePacks(f, repoPath, mypoly)
	f.Close()
	return packs, chunks
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


func WritePacks(f *os.File, repoPath string, poly int) ([]string, []string) {
	const maxPackSize uint = 104857600 // 100 MB


	// create a chunker
	packIds := []string{}
	chunkIds := []string{}
	mychunker := chunker.New(f, chunker.Pol(poly))

	// reuse this buffer
	buf := make([]byte, 8*1024*1024)
	
	i := 0	
	var packId string
	var g0 *os.File
	var g *gzip.Writer
	var packPathErr error	
	
	var curPackSize  uint 
	curPackSize = 0
	

	for {
		chunk, err := mychunker.Next(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		if curPackSize == 0 {
			packId = RandHexString(PACK_ID_LEN)
			packIds = append(packIds, packId)

			packFolder := fmt.Sprintf("%s", packId[0:2])
			packFolderPath := path.Join(repoPath, "packs", packFolder)
			os.MkdirAll(packFolderPath, 0777)
	
			packFilename := fmt.Sprintf("%s.gz", packId)
			packPath := path.Join(packFolderPath, packFilename)	
			fmt.Printf("Creating pack file %s\n", packPath)		
			g0, packPathErr = os.Create(packPath)
			check(packPathErr)
			g = gzip.NewWriter(g0)
		}
		
		i++
		chunkId := sha256.Sum256(chunk.Data)
		curPackSize += chunk.Length
		fmt.Printf("Chunk %d: chunk size %d kB, total size %d kB\n", i, chunk.Length/1024, curPackSize/1024)
		fmt.Printf("  Pack ID: %s\n  Chunk ID: %02x\n", packId, chunkId)
		chunkIds = append(chunkIds, fmt.Sprintf("%064x", chunkId))
		g.Write(chunk.Data)

		if curPackSize >= maxPackSize {
			fmt.Println("Closing pack file")
			g.Close()
			g0.Close()
			curPackSize = 0
		}
	}


	return packIds, chunkIds
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



