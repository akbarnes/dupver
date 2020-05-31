package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"time"
	"strings"
	"archive/tar"
	"compress/gzip"
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
	tarFile, _ := os.Open(filePath)
	PrintFileList(tarFile, commitFile)
	tarFile.Close()
}


func PrintTGZIndex(filePath string, commitFile *os.File) {
	f0, _ := os.Open(filePath)
	f, _ := gzip.NewReader(f0)		
	PrintGZFileList(f, commitFile)
	f.Close()
	f0.Close()
}


func PrintFileList(tarFile *os.File, commitFile *os.File) {
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