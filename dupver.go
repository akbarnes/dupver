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
	"github.com/BurntSushi/toml"
)

type commit struct {
	Message string
	Time string
	Files []string
	Chunks []string
}

type commitHistory struct {
	Commits []commit
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func main() {
	filePtr := flag.String("file", "", "an int")
	backupPtr := flag.Bool("backup", false, "Back up specified file")
	restorePtr := flag.Bool("restore", false, "Restore specified file")
	listPtr := flag.Bool("list", false, "List revisions")
	revisionPtr := flag.Int("revision", -1, "Restore specified revision (default is last)")
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
		fmt.Fprintf(h, "[[commits]]\n")
		// fmt.Fprintf(h, "key=\"2020-05-29\"\n")
		fmt.Fprintf(h, "message=\"%s\"\n", msg)
		fmt.Fprintf(h, "time=\"2020-05-29 5:32pm\"\n")
		// fmt.Fprintf(h, "archive=\"%s\"\n", filePath)
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

			chunkPath := fmt.Sprintf("%s/%02x.gz", chunkFolder, myHash)
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
	} else if *restorePtr == true {
		fmt.Printf("Restoring\n")
		var history commitHistory
		f, _ := os.Open("data/versions.toml")

		if _, err := toml.DecodeReader(f, &history); err != nil {
			log.Fatal(err)
		}

		f.Close()


		fmt.Printf("Number of commits %d\n", len(history.Commits))
		rev := len(history.Commits) - 1

		if *revisionPtr >= 0 {
			rev = *revisionPtr
		}
		
		if (true || len(filePath) == 0) {
			filePath = fmt.Sprintf("snapshot%d.tgz", rev + 1)
		}

		g0, _ := os.Create(filePath)
		g := gzip.NewWriter(g0)

		b := make([]byte, 1024)

		for i, hash := range history.Commits[rev].Chunks {
			chunkPath := fmt.Sprintf("data/%s/%s.gz", hash[0:2], hash)
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

		g.Close()
		g0.Close()

		fmt.Printf("Writing to %s\n", filePath)
	} else if *listPtr == true {
		var history commitHistory
		f, _ := os.Open("data/versions.toml")

		if _, err := toml.DecodeReader(f, &history); err != nil {
			log.Fatal(err)
		}

		f.Close()

		// print a specific revision
		if *revisionPtr >= 0 {
			rev := *revisionPtr - 1
			commit := history.Commits[rev]
			
			fmt.Printf("Revision %d\n", rev + 1)
			fmt.Printf("Time: %s\n", commit.Time)

			if len(commit.Message) > 0 {
				fmt.Printf("Message: %s\n", commit.Message)
			}

			fmt.Printf("Files:\n")
			for j, file := range commit.Files {
				fmt.Printf("  %d: %s\n", j + 1, file)
			}
			fmt.Printf("Chunks: \n")

			for j, hash := range history.Commits[rev].Chunks {
				chunkPath := fmt.Sprintf("data/%s/%s.gz", hash[0:2], hash)
				fmt.Printf("  Chunk %d: %s\n", j + 1, chunkPath)
			}
		} else {
			fmt.Printf("Commit History\n")

			for i, commit := range history.Commits {
				fmt.Printf("Revision %d\n", i + 1)
				fmt.Printf("Time: %s\n", commit.Time)

				if len(commit.Message) > 0 {
					fmt.Printf("Message: %s\n", commit.Message)
				}

				fmt.Printf("Files:\n")
				for j, file := range commit.Files {
					fmt.Printf("  %d: %s\n", j + 1, file)
				}
				fmt.Printf("\n")
			}
		}
	}
}
