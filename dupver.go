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


func PrintCommitHeader(commitFile *os.File, msg string, filePath string) {
	fmt.Fprintf(commitFile, "[[commits]]\n")

	if len(msg) == 0 {
		msg =  strings.Replace(filePath[0:len(filePath)-4], ".\\", "", -1)
	}

	fmt.Fprintf(commitFile, "message=\"%s\"\n", msg)
	t := time.Now()
	fmt.Fprintf(commitFile, "time=\"%s\"\n", t.Format("2006-01-02 15:04:05"))
}


func PrintTarIndex(filePath string, commitFile *os.File) {
	f, _ := os.Open(filePath)
	PrintFileList(f, commitFile)
	f.Close()
}


func PrintTGZIndex(filePath string, commitFile *os.File) {
	f0, _ := os.Open(filePath)
	f, _ := gzip.NewReader(f0)		
	PrintGZFileList(f, commitFile)
	f.Close()
	f0.Close()
}


func PrintFileList(f *os.File, commitFile *os.File) {
	fmt.Fprintf(commitFile, "files = [\n")


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

		i++
		fmt.Printf("File %d: %s\n", i, hdr.Name)
		fmt.Fprintf(commitFile, "  \"%s\",\n", hdr.Name)
	}

	fmt.Fprint(commitFile, "]\n")
}


func PrintGZFileList(f *gzip.Reader, commitFile *os.File) {
	fmt.Fprintf(commitFile, "files = [\n")


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

		i++
		fmt.Printf("File %d: %s\n", i, hdr.Name)
		fmt.Fprintf(commitFile, "  \"%s\",\n", hdr.Name)
	}

	fmt.Fprint(commitFile, "]\n")
}


func PackTar(filePath string, commitFile *os.File, mypoly int) {
	f, _ := os.Open(filePath)
	WritePacks(f, commitFile, mypoly)
	f.Close()
}


func PackTGZ(filePath string, commitFile *os.File, mypoly int) {
	f0, _ := os.Open(filePath)
	f, _ := gzip.NewReader(f0)
	WritePacksFromGZ(f, commitFile, mypoly)
	f.Close()
	f0.Close()
}


func WritePacks(f *os.File, commitFile *os.File, poly int) {
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
	g, _ := os.Create(filePath)
	WriteChunks(g, chunks)
	g.Close()
}


func UnpackTGZ(filePath string, chunks []string) {
	g0, _ := os.Create(filePath)
	g := gzip.NewWriter(g0)
	WriteGZChunks(g, chunks)
	g.Close()
	g0.Close()
}


func WriteChunks(g *os.File, chunks []string) {
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


func GetRevIndex(revision int, numCommits int) int {
	revIndex := numCommits - 1
	
	if revision > 0 {
		revIndex = revision - 1
	} else if revision < 0 {
		revIndex = numCommits + revision
	}

	return revIndex
}


func PrintRevision(history commitHistory, revIndex int, maxFiles int) {
	commit := history.Commits[revIndex]
				
	fmt.Printf("Revision %d\n", revIndex + 1)
	fmt.Printf("Time: %s\n", commit.Time)

	if len(commit.Message) > 0 {
		fmt.Printf("Message: %s\n", commit.Message)
	}

	fmt.Printf("Files:\n")
	for j, file := range commit.Files {
		fmt.Printf("  %d: %s\n", j + 1, file)

		if j > maxFiles && maxFiles > 0 {
			fmt.Printf("  ...\n  Skipping %d more files\n", len(commit.Files) - maxFiles)
			break
		}
	}
}


func main() {
	// constants
	mypoly := 0x3DA3358B4DC173
	commitLogPath := fmt.Sprintf(".dupver/versions.toml")

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
	

	if *initPtr {
		os.Mkdir("./.dupver", 0777)
		f, _ := os.Create(commitLogPath)
		f.Close()
	} else if *backupPtr {
		fmt.Println("Backing up ", filePath)
		commitFile, _ := os.OpenFile(commitLogPath, os.O_APPEND|os.O_WRONLY, 0600)
		PrintCommitHeader(commitFile, msg, filePath)
		PrintTarIndex(filePath, commitFile)
		PackTar(filePath, commitFile, mypoly)
		commitFile.Close()
	} else if *restorePtr {
		fmt.Printf("Restoring\n")
		var history commitHistory
		f, _ := os.Open(commitLogPath)

		if _, err := toml.DecodeReader(f, &history); err != nil {
			log.Fatal(err)
		}

		f.Close()

		fmt.Printf("Number of commits %d\n", len(history.Commits))
		revIndex := GetRevIndex(revision, len(history.Commits))
		fmt.Printf("Restoring commit %d\n", revIndex)
		
		if (true || len(filePath) == 0) {
			filePath = fmt.Sprintf("snapshot%d.tar", revIndex + 1)
		}

		fmt.Printf("Writing to %s\n", filePath)
		UnpackTar(filePath, history.Commits[revIndex].Chunks) 
	} else if *listPtr {
		var history commitHistory
		f, _ := os.Open(".dupver/versions.toml")

		if _, err := toml.DecodeReader(f, &history); err != nil {
			log.Fatal(err)
		}

		f.Close()

		// print a specific revision
		if revision != 0 {
			revIndex := GetRevIndex(revision, len(history.Commits))
			PrintRevision(history, revIndex, 0)
		} else {
			fmt.Printf("Commit History\n")

			for i:=0; i < len(history.Commits); i++ {
				PrintRevision(history, i, 10)
			}
		}
	} else {
		fmt.Println("No command specified, exiting")
		fmt.Println("For available commands run: dupver -help")
	}
}
