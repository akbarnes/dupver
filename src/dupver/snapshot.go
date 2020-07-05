package dupver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	)

type Commit struct {
	TarFileName string
	ID          string
	Message     string
	Time        string
	ParentIDs   []string
	Files       []fileInfo
	ChunkIDs    []string
}

type Branch struct {
	CommitID string
}

type fileInfo struct {
	Path    string
	ModTime string
	Size    int64
	Hash    string
	// Permissions int
}

const SNAPSHOT_ID_LEN int = 40
const PACK_ID_LEN int = 64
const CHUNK_ID_LEN int = 64
const TREE_ID_LEN int = 40

func CommitFile(filePath string, parentIds []string, msg string, verbosity int) string {
	var myWorkDirConfig workDirConfig
	t := time.Now()

	var mySnapshot Commit
	mySnapshot.ID = RandHexString(SNAPSHOT_ID_LEN)
	mySnapshot.Time = t.Format("2006/01/02 15:04:05")
	mySnapshot.ParentIDs = parentIds
	mySnapshot = UpdateMessage(mySnapshot, msg, filePath)
	mySnapshot.Files, myWorkDirConfig = ReadTarFileIndex(filePath, verbosity)

	if verbosity >= 2 {
		fmt.Printf("Repo config: %s\n", myWorkDirConfig.RepoPath)
	}

	myRepoConfig := ReadRepoConfigFile(path.Join(myWorkDirConfig.RepoPath, "config.toml"))

	chunkIDs, chunkPacks := PackFile(filePath, myWorkDirConfig.RepoPath, myRepoConfig.ChunkerPolynomial, verbosity)
	mySnapshot.ChunkIDs = chunkIDs

	snapshotFolder := path.Join(myWorkDirConfig.RepoPath, "snapshots", myWorkDirConfig.WorkDirName)
	snapshotBasename := fmt.Sprintf("%s-%s", t.Format("2006-01-02-T15-04-05"), mySnapshot.ID[0:40])
	os.Mkdir(snapshotFolder, 0777)
	snapshotPath := path.Join(snapshotFolder, snapshotBasename + ".json")
	WriteSnapshot(snapshotPath, mySnapshot)

	branchFolder := path.Join(myWorkDirConfig.RepoPath, "branches", myWorkDirConfig.WorkDirName)
	branchPath := path.Join(branchFolder,  "head.json")
	WriteBranch(branchPath, mySnapshot.ID)

	treeFolder := path.Join(myWorkDirConfig.RepoPath, "trees")
	treeBasename := mySnapshot.ID[0:40]
	os.Mkdir(treeFolder, 0777)
	treePath := path.Join(treeFolder, treeBasename+".json")
	WriteTree(treePath, chunkPacks)

	if verbosity >= 1 {
		fmt.Printf("Created snapshot %s (%s)\n", mySnapshot.ID[0:16], mySnapshot.ID)
	} else {
		fmt.Println(mySnapshot.ID)
	}
	return mySnapshot.ID
}

func UpdateMessage(mySnapshot Commit, msg string, filePath string) Commit {
	if len(msg) == 0 {
		msg = strings.Replace(filePath[0:len(filePath)-4], ".\\", "", -1)
	}

	mySnapshot.Message = msg
	return mySnapshot
}


func WriteSnapshot(snapshotPath string, mySnapshot Commit) {
	f, err := os.Create(snapshotPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot file %s", snapshotPath))
	}
	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(mySnapshot)
	f.Close()
}

func WriteBranch(branchPath string, commitId string) {
	myBranch := Branch{CommitID: commitId}

	f, err := os.Create(branchPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create branch file %s", branchPath))
	}
	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(myBranch)
	f.Close()
}

func ReadSnapshot(snapshot string, cfg workDirConfig) Commit {
	snapshotPaths := ListSnapshots(cfg)

	for _, snapshotPath := range snapshotPaths {
		n := len(snapshotPath)
		snapshotId := snapshotPath[n-SNAPSHOT_ID_LEN-5 : n-5]

		if snapshotId[0:len(snapshot)] == snapshot {
			return ReadSnapshotFile(snapshotPath)
		}
	}

	log.Fatal(fmt.Sprintf("Error: Could not find snapshot %s in repo", snapshot))
	return Commit{}
}

func ReadSnapshotFile(snapshotPath string) Commit {
	var mySnapshot Commit
	f, err := os.Open(snapshotPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&mySnapshot); err != nil {
		panic(fmt.Sprintf("Error:could not decode snapshot file %s", snapshotPath))
	}

	f.Close()
	return mySnapshot
}

func ReadSnapshotId(snapshotId string, cfg workDirConfig) (Commit, error) {
	snapshotPaths := ListSnapshots(cfg)

	for _, snapshotPath := range snapshotPaths {
		n := len(snapshotPath)
		sid := snapshotPath[n-SNAPSHOT_ID_LEN-5 : n-5]

		if sid[0:8] == snapshotId {
			return ReadSnapshotFile(snapshotPath), nil
		}
	}

	return Commit{}, errors.New(fmt.Sprintf("Could not find snapshot %s", snapshotId))
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
	snapshotsFolder := path.Join(cfg.RepoPath, "snapshots", cfg.WorkDirName)
	snapshotGlob := path.Join(snapshotsFolder, "*.json")
	// fmt.Println(snapshotGlob)
	snapshotPaths, err := filepath.Glob(snapshotGlob)

	if err != nil {
		panic(fmt.Sprintf("Error listing snapshots glob %s", snapshotGlob))
	}
	return snapshotPaths
}

func PrintSnapshots(cfg workDirConfig, snapshot string, verbosity int) {
	// fmt.Printf("Verbosity = %d\n", verbosity)
	snapshotPaths := ListSnapshots(cfg)
	// print a specific revision
	if len(snapshot) == 0 {
		if verbosity >= 1 {
			fmt.Println("Snapshot History")
		}

		for _, snapshotPath := range snapshotPaths {
			// fmt.Printf("Path: %s\n", snapshotPath)
			PrintSnapshot(ReadSnapshotFile(snapshotPath), 10, verbosity)
		}
	} else {
		if verbosity >= 1 {
			fmt.Println("Snapshot")
		}

		for _, snapshotPath := range snapshotPaths {
			// if i >= 1 {
			// 	fmt.Println("\n")
			// }

			n := len(snapshotPath)
			snapshotId := snapshotPath[n-SNAPSHOT_ID_LEN-5 : n-5]

			if snapshotId[0:8] == snapshot {
				PrintSnapshot(ReadSnapshotFile(snapshotPath), 0, verbosity)
			}
		}
	}
}

func PrintSnapshot(mySnapshot Commit, maxFiles int, verbosity int) {
	if verbosity <= 0 {
		fmt.Printf("%s\n", mySnapshot.ID)
		return
	}

	fmt.Printf("%sID: %s (%s)%s\n", colorYellow, mySnapshot.ID[0:8], mySnapshot.ID, colorReset)
	fmt.Printf("Time: %s\n", mySnapshot.Time)

	if len(mySnapshot.Message) > 0 {
		fmt.Printf("Message: %s\n", mySnapshot.Message)
	}

	fmt.Printf("\n")

	// fmt.Printf("Files:\n")
	// for i, file := range mySnapshot.Files {
	// 	fmt.Printf("  %d: %s\n", i+1, file.Path)

	// 	if i > maxFiles && maxFiles > 0 {
	// 		fmt.Printf("  ...\n  Skipping %d more files\n", len(mySnapshot.Files)-maxFiles)
	// 		break
	// 	}
	// }
}

