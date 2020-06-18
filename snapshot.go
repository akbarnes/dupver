package main

import (
	"os"
	"io"
	"bufio"
	"log"
	"fmt"
	"path"
	"path/filepath"
	"time"
	"crypto/sha256"
	"strings"
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
	// ChunkPacks map[string]string
	ChunkIDs []string
	// PackIndexes []packIndex
	Tags []string
}


type fileInfo struct {
	Path string
	ModTime string
	Size int64
	Hash string
	// Permissions int
}

const SNAPSHOT_ID_LEN int = 40 
const PACK_ID_LEN int = 64
const CHUNK_ID_LEN int = 64
const TREE_ID_LEN int = 40


func CommitFile(filePath string, msg string, verbosity int) string {
	var myWorkDirConfig workDirConfig
	t := time.Now()

	var mySnapshot commit
	mySnapshot.ID = RandHexString(SNAPSHOT_ID_LEN)
	mySnapshot.Time = t.Format("2006/01/02 15:04:05")
	mySnapshot.TarFileName = filePath
	// mySnapshot = UpdateTags(mySnapshot, tagName)
	mySnapshot = UpdateMessage(mySnapshot, msg, filePath)		
	mySnapshot.Files, myWorkDirConfig = ReadTarFileIndex(filePath)
	fmt.Printf("Repo config: %s\n", myWorkDirConfig.RepoPath)
	myRepoConfig := ReadRepoConfigFile(path.Join(myWorkDirConfig.RepoPath, "config.toml"))
	
	chunkIDs, chunkPacks := PackFile(filePath, myWorkDirConfig.RepoPath, myRepoConfig.ChunkerPolynomial, verbosity)
	mySnapshot.ChunkIDs = chunkIDs

	snapshotFolder := path.Join(myWorkDirConfig.RepoPath, "snapshots", myWorkDirConfig.WorkDirName)
	snapshotBasename := fmt.Sprintf("%s-%s", t.Format("2006-01-02-T15-04-05"), mySnapshot.ID[0:40])		
	os.Mkdir(snapshotFolder, 0777)
	snapshotPath := path.Join(snapshotFolder, snapshotBasename + ".json")
	WriteSnapshot(snapshotPath, mySnapshot)

	treeFolder := path.Join(myWorkDirConfig.RepoPath, "trees")
	treeBasename := mySnapshot.ID[0:40]
	os.Mkdir(treeFolder, 0777)
	treePath := path.Join(treeFolder, treeBasename + ".json")
	WriteTree(treePath, chunkPacks)

	if verbosity >= 1 {
		fmt.Printf("Created snapshot %s\n", mySnapshot.ID[0:16])
	}

	return mySnapshot.ID
}


func UpdateTags(mySnapshot commit, tagName string) commit {
	if len(tagName) > 0 {
		mySnapshot.Tags = []string{tagName}
	}

	return mySnapshot
}


func UpdateMessage(mySnapshot commit, msg string, filePath string) commit {
	if len(msg) == 0 {
		msg =  strings.Replace(filePath[0:len(filePath)-4], ".\\", "", -1)
	}

	mySnapshot.Message = msg
	return mySnapshot
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
	// var baseFolder string
	var configPath string
	maxFiles := 10


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

		// if i == 0 {
		// 	baseFolder = hdr.Name
		// 	myConfig.WorkDirName = baseFolder
		// 	configPath = path.Join(baseFolder,".dupver","config.toml")
		// 	// fmt.Printf("Base folder: %s\nConfig path: %s\n", baseFolder, configPath)
		// }


		if strings.HasSuffix(hdr.Name, ".dupver/config.toml") {
			fmt.Printf("Matched config path %s\n", configPath)
			if _, err := toml.DecodeReader(tr, &myConfig); err != nil {
				log.Fatal(err)
			}

			// fmt.Printf("Read config\nworkdir name: %s\nrepo path: %s\n", myConfig.WorkDirName, myConfig.RepoPath)
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

		if i <= maxFiles {
			fmt.Printf("File %d: %s\n", i, hdr.Name)
		}

		files = append(files, myFileInfo)
	}

	if i > maxFiles && maxFiles > 0 {
		fmt.Printf("...\nSkipping %d more files\n", i - maxFiles)
	}

	return files, myConfig
}


func WriteSnapshot(snapshotPath string, mySnapshot commit) {
	f, _ := os.Create(snapshotPath)
	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(mySnapshot)
	f.Close()
}


func WriteTree(treePath string, chunkPacks map[string]string) {
	f, _ := os.Create(treePath)
	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(chunkPacks)
	f.Close()
}


func ReadTrees(repoPath string) map[string]string {
	treesGlob := path.Join(repoPath, "trees", "*.json")
	// fmt.Println(treesGlob)
	treePaths, err := filepath.Glob(treesGlob)
	check(err)
	chunkPacks := make(map[string]string)	

	
	for _, treePath := range treePaths {
		treePacks := make(map[string]string)	
		
		f, _ := os.Open(treePath)
		myDecoder := json.NewDecoder(f)
	
		if err := myDecoder.Decode(&treePacks); err != nil {
			log.Fatal(err)
		}	

		// TODO: handle supersedes to allow repacking files			
		for k, v := range treePacks {
			chunkPacks[k] = v
		}

		f.Close()
	}

	return chunkPacks
}


func ReadSnapshot(snapshot string, cfg workDirConfig) commit {
	snapshotPaths := ListSnapshots(cfg)

	for _, snapshotPath := range snapshotPaths {
		n := len(snapshotPath)
		snapshotId := snapshotPath[n-SNAPSHOT_ID_LEN-5:n-5]

		if snapshotId[0:len(snapshot)] == snapshot {
			return ReadSnapshotFile(snapshotPath)
		}
	}

	panic("Could not find snapshot")
}


func ReadSnapshotFile(snapshotPath string) (commit) {
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

func ListSnapshots(cfg workDirConfig) []string {
	snapshotGlob := path.Join(cfg.RepoPath, "snapshots", cfg.WorkDirName, "*.json")
	// fmt.Println(snapshotGlob)
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	check(err)
	return snapshotPaths
}

func PrintSnapshots(snapshotPaths[] string, snapshot string) {
	// print a specific revision
	if len(snapshot) == 0 {
		fmt.Printf("Snapshot History\n")

		for _, snapshotPath := range snapshotPaths {
			fmt.Printf("Path: %s\n", snapshotPath)
			PrintSnapshot(ReadSnapshotFile(snapshotPath), 10)
		}			
	} else {
		fmt.Println("Snapshot")

		for _, snapshotPath := range snapshotPaths {
			n := len(snapshotPath)
			snapshotId := snapshotPath[n-SNAPSHOT_ID_LEN-5:n-5]
		
			if snapshotId[0:8] == snapshot {
				PrintSnapshot(ReadSnapshotFile(snapshotPath), 0)
			}
		}	
	}
}



func PrintSnapshot(mySnapshot commit, maxFiles int) {			
	fmt.Printf("Time: %s\n", mySnapshot.Time)
	fmt.Printf("ID: %s\n", mySnapshot.ID[0:8])

	if len(mySnapshot.Message) > 0 {
		fmt.Printf("Message: %s\n", mySnapshot.Message)
	}

	if len(mySnapshot.Tags) > 0 {
		fmt.Printf("Tags:\n")
		for _,  tag := range mySnapshot.Tags {
			fmt.Printf("  %s\n", tag)
		}
	}

	fmt.Printf("Files:\n")
	for i, file := range mySnapshot.Files {
		fmt.Printf("  %d: %s\n", i + 1, file.Path)

		if i > maxFiles && maxFiles > 0 {
			fmt.Printf("  ...\n  Skipping %d more files\n", len(mySnapshot.Files) - maxFiles)
			break
		}
	}
}

