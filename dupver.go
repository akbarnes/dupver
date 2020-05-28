package main


import (
    // "bufio"
    "fmt"
	"io"
	"crypto/sha256"
    // "io/ioutil"
	"os"
	"github.com/restic/chunker"
)
 	

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func main() {
	// chunky =  Chunker(rd io.Reader, pol Pol) 
	fmt.Println("Welcome to the playground!")
	f, _ := os.Open("ACTIVSg70k.RAW")

    b1 := make([]byte, 24)
    n1, _ := f.Read(b1)

	fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))
	
	// generate 32MiB of deterministic pseudo-random data
	// data := getRandom(23, 32*1024*1024)

	// create a chunker
	mychunker := chunker.New(f, chunker.Pol(0x3DA3358B4DC173))

	// reuse this buffer
	buf := make([]byte, 8*1024*1024)

	for i := 0; i < 5; i++ {
		chunk, err := mychunker.Next(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		myHash := sha256.Sum256(chunk.Data)
		fmt.Printf("%d %02x\n", chunk.Length, myHash)
		chunkPath := fmt.Sprintf("chunks/%d-%02x.dat", i, myHash[0:5])
		g, _ := os.Create(chunkPath)
		g.Write(chunk.Data)
		g.Close()
	}

	f.Close()
}
