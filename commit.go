package main

import (
	"os"
	"io"
	"log"
	"fmt"
	// "time"
	// "strings"
	"github.com/BurntSushi/toml"
	"archive/tar"
)

type commit struct {
	ID string
	Message string
	Time string
	Files []string
	Chunks []string
}

type commitHistory struct {
	Commits []commit
}


// func PrintCommitHeader(commitFile *os.File, msg string, filePath string) {
// 	if len(msg) == 0 {
// 		msg =  strings.Replace(filePath[0:len(filePath)-4], ".\\", "", -1)
// 	}

// 	fmt.Fprintf(commitFile, "message=\"%s\"\n", msg)
// 	t := time.Now()
// 	fmt.Fprintf(commitFile, "time=\"%s\"\n", t.Format("2006-01-02 15:04:05"))
// }


func PrintTarFileIndex(filePath string, commitFile *os.File) {
	tarFile, _ := os.Open(filePath)
	PrintTarIndex(tarFile, commitFile)
	tarFile.Close()
}


func PrintTarIndex(tarFile *os.File, commitFile *os.File) {
	fmt.Fprintf(commitFile, "files = [\n")


	// Open and iterate through the files in the archive.
	tr := tar.NewReader(tarFile)
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


func ReadSnapshot(snapshotPath string) (commit) {
	var mySnapshot commit
	f, _ := os.Open(snapshotPath)

	if _, err := toml.DecodeReader(f, &mySnapshot); err != nil {
		log.Fatal(err)
	}

	f.Close()
	return mySnapshot
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


func PrintSnapshot(mySnapshot commit, maxFiles int) {			
	fmt.Printf("Time: %s\n", mySnapshot.Time)
	fmt.Printf("ID: %s\n", mySnapshot.ID[0:8])

	if len(mySnapshot.Message) > 0 {
		fmt.Printf("Message: %s\n", mySnapshot.Message)
	}

	fmt.Printf("Files:\n")
	for j, file := range mySnapshot.Files {
		fmt.Printf("  %d: %s\n", j + 1, file)

		if j > maxFiles && maxFiles > 0 {
			fmt.Printf("  ...\n  Skipping %d more files\n", len(mySnapshot.Files) - maxFiles)
			break
		}
	}
}

