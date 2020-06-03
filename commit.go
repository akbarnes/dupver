package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"path"
	// "time"
	// "strings"
	"github.com/BurntSushi/toml"
	"archive/tar"
)

type commit struct {
	TarFileName string
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


func ReadTarFileIndex(filePath string) ([]string, workDirConfig) {
	tarFile, _ := os.Open(filePath)
	files, myConfig := ReadTarIndex(tarFile)
	tarFile.Close()

	return files, myConfig
}


func ReadTarIndex(tarFile *os.File) ([]string, workDirConfig) {
	files := []string{}
	var myConfig workDirConfig
	var baseFolder string
	var configPath string


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

		if i == 0 {
			baseFolder = hdr.Name
			myConfig.WorkDirName = baseFolder
			configPath = path.Join(baseFolder,".dupver","config.toml")
			fmt.Printf("Base folder: %s\nConfig path: %s\n", baseFolder, configPath)
		}


		if hdr.Name == configPath {
			fmt.Printf("Matched config path %s\n", configPath)
			if _, err := toml.DecodeReader(tr, &myConfig); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Read config\nworkdir name: %s\nrepo path: %s\n", myConfig.WorkDirName, myConfig.RepositoryPath)
		}

		i++
		fmt.Printf("File %d: %s\n", i, hdr.Name)
		files = append(files, hdr.Name)	
	}

	return files, myConfig
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

