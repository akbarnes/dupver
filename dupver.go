package main


import (
	"flag"
    "fmt"
	"io"
	"crypto/sha256"
	"os"
	"github.com/restic/chunker"
)
 	

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func main() {
	filePtr := flag.String("file", "ACTIVSg70k.RAW", "an int")
	backupPtr := flag.Bool("backup", false, "Back up specified file")
	msgPtr := flag.String("message", "", "commit message")
	
	flag.Parse()
	
	filePath := *filePtr
	msg := *msgPtr
	// filePath = "ACTIVSg70k.RAW"	

	if (*backupPtr == true) {
		fmt.Println("Backing up ", filePath)

		// chunky =  Chunker(rd io.Reader, pol Pol) 
		f, _ := os.Open(filePath)
		
		// generate 32MiB of deterministic pseudo-random data
		// data := getRandom(23, 32*1024*1024)
		os.Mkdir("./data", 0777)

		// create a chunker
		mychunker := chunker.New(f, chunker.Pol(0x3DA3358B4DC173))

		// reuse this buffer
		buf := make([]byte, 8*1024*1024)

		// os.MkdirAll("data/tree")
		treePath := fmt.Sprintf("data/versions.toml")
		h, _ := os.Create(treePath)

		fmt.Fprintf(h, "[versions.2020-05-29]\n")
		fmt.Fprintf(h, "message=\"%s\"\n", msg)
		fmt.Fprintf(h, "file=\"%s\"\n", filePath)
		fmt.Fprintf(h, "chunks = [\n")

		i := 0

		for {
			chunk, err := mychunker.Next(buf)
			if err == io.EOF {
				break
			}

			if err != nil {
				panic(err)
			}
			
			i += 1
			myHash := sha256.Sum256(chunk.Data)
			fmt.Printf("Chunk %d: %d kB, %02x\n", i, chunk.Length/1024, myHash)
			fmt.Fprintf(h, "  \"%02x\",\n", myHash)

			chunkFolder := fmt.Sprintf("data/%02x", myHash[0:1])
			os.MkdirAll(chunkFolder, 0777)

			chunkPath := fmt.Sprintf("%s/%02x", chunkFolder, myHash)
			g, _ := os.Create(chunkPath)
			g.Write(chunk.Data)
			g.Close()
		}

		f.Close()
		fmt.Fprintf(h, "]\n")
		h.Close()
	}
}
