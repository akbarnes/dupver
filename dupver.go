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
	f, err := os.Open("scoop_apps.txt")
	check(err)

    b1 := make([]byte, 24)
    n1, err := f.Read(b1)
    check(err)
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

		fmt.Printf("%d %02x\n", chunk.Length, sha256.Sum256(chunk.Data))
	}
}
