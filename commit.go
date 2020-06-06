package main

import (
	"os"
	"io"
	"bufio"
	"log"
	"fmt"
	"path"
	// "time"
	"crypto/sha256"
	// "strings"
	"github.com/BurntSushi/toml"
	"encoding/json"
	"archive/tar"
)

type commit struct {
	TarFileName string
	ID string
	Message string
	Time string
	Files []fileInfo
	Chunks []string
}

type commitHistory struct {
	Commits []commit
}

type fileInfo struct {
	Path string
	ModTime string
	Size int64
	Hash string
	// Permissions int
}


func ReadTarFileIndex(filePath string) ([]fileInfo, workDirConfig) {
	tarFile, _ := os.Open(filePath)
	files, myConfig := ReadTarIndex(tarFile)
	tarFile.Close()

	return files, myConfig
}


func ReadTarIndex(tarFile *os.File) ([]fileInfo, workDirConfig) {
	files := []fileInfo{}
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
			// fmt.Printf("Base folder: %s\nConfig path: %s\n", baseFolder, configPath)
		}


		if hdr.Name == configPath {
			fmt.Printf("Matched config path %s\n", configPath)
			if _, err := toml.DecodeReader(tr, &myConfig); err != nil {
				log.Fatal(err)
			}

			// fmt.Printf("Read config\nworkdir name: %s\nrepo path: %s\n", myConfig.WorkDirName, myConfig.RepositoryPath)
		}

		var myFileInfo fileInfo

	    bytes := make([]byte, hdr.Size)

	    bufr := bufio.NewReader(tr)
	    _, err = bufr.Read(bytes)

		// Name              |   256B | unlimited | unlimited
		// Linkname          |   100B | unlimited | unlimited
		// Size              | uint33 | unlimited |    uint89
		// Mode              | uint21 |    uint21 |    uint57
		// Uid/Gid           | uint21 | unlimited |    uint57
		// Uname/Gname       |    32B | unlimited |       32B
		// ModTime           | uint33 | unlimited |     int89
		// AccessTime        |    n/a | unlimited |     int89
		// ChangeTime        |    n/a | unlimited |     int89
		// Devmajor/Devminor | uint21 |    uint21 |    uint57

	    myFileInfo.Path = hdr.Name
	    myFileInfo.Size = hdr.Size
		myFileInfo.Hash = fmt.Sprintf("%02x", sha256.Sum256(bytes))
		myFileInfo.ModTime = hdr.ModTime.Format("2006/01/02 15:04:05")

		i++
		fmt.Printf("File %d: %s\n", i, hdr.Name)
		files = append(files, myFileInfo)
	}

	return files, myConfig
}


func WriteSnapshot(snapshotPath string, mySnapshot commit) {
	snapshotFile, _ := os.Create(snapshotPath)
	myEncoder := json.NewEncoder(snapshotFile)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(mySnapshot)
	snapshotFile.Close()
}


func ReadSnapshot(snapshotPath string) (commit) {
	var mySnapshot commit
	f, _ := os.Open(snapshotPath)
	myDecoder := json.NewDecoder(f)


	if err := myDecoder.Decode(&mySnapshot); err != nil {
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
		fmt.Printf("  %d: %s\n", j + 1, file.Path)

		if j > maxFiles && maxFiles > 0 {
			fmt.Printf("  ...\n  Skipping %d more files\n", len(mySnapshot.Files) - maxFiles)
			break
		}
	}
}

