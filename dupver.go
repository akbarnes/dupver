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
	"time"
	"strings"
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

func writePacks(f *File, h *File, poly int) {
	// create a chunker
	mychunker := chunker.New(f, chunker.Pol(poly))

	// reuse this buffer
	buf := make([]byte, 8*1024*1024)
	
	fmt.Fprintf(h, "chunks = [\n")

	i = 0

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

		chunkFolder := fmt.Sprintf(".dupver/%02x", myHash[0:1])
		os.MkdirAll(chunkFolder, 0777)

		chunkPath := fmt.Sprintf("%s/%02x.gz", chunkFolder, myHash)
		g0, _ := os.Create(chunkPath)
		g := gzip.NewWriter(g0)
		g.Write(chunk.Data)
		g.Close()
		g0.Close()
	}


	fmt.Fprintf(h, "]\n\n")
}


func main() {
	// constants
	mypoly := 0x3DA3358B4DC173

	initPtr := flag.Bool("init", false, "Initialize the repository")
	backupPtr := flag.Bool("backup", false, "Back up specified file")
	restorePtr := flag.Bool("restore", false, "Restore specified file")
	listPtr := flag.Bool("list", false, "List revisions")

	var filePath string
	var msg string
	var revision int

	flag.StringVar(&filePath, "file", "", "Archive path")
	flag.StringVar(&filePath, "f", "", "Archive path (shorthand)")

	flag.IntVar(&revision, "revision", 0, "Specify revision (default is last)")
	flag.IntVar(&revision, "r", 0, "Specify revision (shorthand)")


	flag.StringVar(&msg, "message", "", "Commit message")
	flag.StringVar(&msg, "m", "", "Commit message (shorthand)")

	
	flag.Parse()
	
	treePath := fmt.Sprintf(".dupver/versions.toml")

	if *initPtr {
		os.Mkdir("./.dupver", 0777)
		f, _ := os.Create(treePath)
		f.Close()
	} else if *backupPtr {
		fmt.Println("Backing up ", filePath)

		f0, _ := os.Open(filePath)
		f, _ := gzip.NewReader(f0)

		h, _ := os.OpenFile(treePath, os.O_APPEND|os.O_WRONLY, 0600)
		fmt.Fprintf(h, "[[commits]]\n")

		if len(msg) == 0 {
			msg =  strings.Replace(filePath[0:len(filePath)-4],".\\","",-1)
		}

		fmt.Fprintf(h, "message=\"%s\"\n", msg)
		t := time.Now()
		fmt.Fprintf(h, "time=\"%s\"\n", t.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(h, "files = [\n")


		// Open and iterate through the files in the archive.
		tr := tar.NewReader(f)
		i := 0
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break // End of archive
			}
			if err != nil {
				log.Fatal(err)
			}

			i += 1
			fmt.Printf("File %d: %s\n", i, hdr.Name)
			fmt.Fprintf(h, "  \"%s\",\n", hdr.Name)
		}

		fmt.Fprint(h, "]\n")

		f.Close()
		f0.Close()
		f0, _ = os.Open(filePath)
		f, _ = gzip.NewReader(f0)
		
		writePacks(f, h, mypoly)

		f.Close()
		f0.Close()
		h.Close()
	} else if *restorePtr {
		fmt.Printf("Restoring\n")
		var history commitHistory
		treePath := ".dupver/versions.toml"
		f, _ := os.Open(treePath)

		if _, err := toml.DecodeReader(f, &history); err != nil {
			log.Fatal(err)
		}

		f.Close()


		fmt.Printf("Number of commits %d\n", len(history.Commits))
		nc := len(history.Commits) 
		rev := nc - 1

		if revision > 0 {
			rev = revision - 1
		} else if revision < 0 {
			rev = nc + revision
		}

		fmt.Printf("Restoring commit %d\n", rev)
		
		if (true || len(filePath) == 0) {
			filePath = fmt.Sprintf("snapshot%d.tgz", rev + 1)
		}

		g0, _ := os.Create(filePath)
		g := gzip.NewWriter(g0)

		b := make([]byte, 1024)

		for i, hash := range history.Commits[rev].Chunks {
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

		g.Close()
		g0.Close()

		fmt.Printf("Writing to %s\n", filePath)
	} else if *listPtr {
		var history commitHistory
		f, _ := os.Open(".dupver/versions.toml")

		if _, err := toml.DecodeReader(f, &history); err != nil {
			log.Fatal(err)
		}

		f.Close()

		// print a specific revision
		if revision != 0 {
			nc := len(history.Commits) 
			rev := nc - 1
	
			if revision > 0 {
				rev = revision - 1
			} else if revision < 0 {
				rev = nc + revision
			}

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
			// fmt.Printf("Chunks: \n")

			// for j, hash := range history.Commits[rev].Chunks {
			// 	chunkPath := fmt.Sprintf(".dupver/%s/%s.gz", hash[0:2], hash)
			// 	fmt.Printf("  Chunk %d: %s\n", j + 1, chunkPath)
			// }
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
	} else {
		fmt.Println("No command specified, exiting")
		fmt.Println("For available commands run: dupver -help")
	}
}
