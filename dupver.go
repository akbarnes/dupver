package main


import (
	"flag"
    "fmt"
	"io"
	"crypto/sha256"
	"os"
	"github.com/restic/chunker"
	"compress/gzip"
	"archive/tar"
	"log"
)
 	

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func main() {
	filePtr := flag.String("file", "ACTIVSg70k.RAW", "an int")
	backupPtr := flag.Bool("backup", false, "Back up specified file")
	// restorePtr := flag.Bool("restore", false, "Restore specified file")
	msgPtr := flag.String("message", "", "commit message")
	
	flag.Parse()
	
	filePath := *filePtr
	msg := *msgPtr
	// filePath = "ACTIVSg70k.RAW"	

	if (*backupPtr == true) {
		fmt.Println("Backing up ", filePath)

		f0, _ := os.Open(filePath)
		f, _ := gzip.NewReader(f0)

		// os.MkdirAll("data/tree")
		os.Mkdir("./data", 0777)
		treePath := fmt.Sprintf("data/versions.toml")
		h, _ := os.Create(treePath)
		fmt.Fprintf(h, "[versions.2020-05-29]\n")
		fmt.Fprintf(h, "message=\"%s\"\n", msg)
		fmt.Fprintf(h, "archive=\"%s\"\n", filePath)
		fmt.Fprintf(h, "files = [\n")


		// Open and iterate through the files in the archive.
		tr := tar.NewReader(f)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break // End of archive
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n,", hdr.Name)
			fmt.Fprintf(h, "  \"%s\",\n", hdr.Name)
		}

		fmt.Fprint(h, "]\n")

		f.Close()
		f0.Close()
		f0, _ = os.Open(filePath)
		f, _ = gzip.NewReader(f0)
		
		// create a chunker
		mychunker := chunker.New(f, chunker.Pol(0x3DA3358B4DC173))

		// reuse this buffer
		buf := make([]byte, 8*1024*1024)



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
			g0, _ := os.Create(chunkPath)
			g := gzip.NewWriter(g0)
			g.Write(chunk.Data)
			g.Close()
			g0.Close()
		}

		f.Close()
		f0.Close()
		fmt.Fprintf(h, "]\n")
		h.Close()
	}
}
