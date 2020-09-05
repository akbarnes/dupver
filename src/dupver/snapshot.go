package dupver

import (
	"encoding/json"
	"errors"
	"fmt"
	// "log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	)

type Commit struct {
	// TarFileName string
	ID          string
	Message     string
	Time        string
	ParentIDs   []string
	Files       []fileInfo
	ChunkIDs    []string
}

type Head struct {
	BranchName string
	CommitID string // use this for detached head, but do I need this?
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

// CopyAll
// CopyBranch - a bit trickier as need to rename branches to reponame.branch
//              and repo will need to have a unique name
//              stick with names in workdir for now
func CopySnapshot(snapshotId string, opts Options) {
	fromPath := opts.RepoPath
	toPath := opts.DestRepoPath

	if len(fromPath) == 0 {
		fromPath = opts.
	}
	fmt.Println("Copying snapshot %s from %s to %s")
}


func CommitFile(filePath string, parentIds []string, msg string, opts Options) Head {
	var myWorkDirConfig workDirConfig
	t := time.Now()

	var mySnapshot Commit
	var myHead Head
	mySnapshot.ID = RandHexString(SNAPSHOT_ID_LEN)
	mySnapshot.Time = t.Format("2006/01/02 15:04:05")
	mySnapshot = UpdateMessage(mySnapshot, msg, filePath)
	mySnapshot.Files, myWorkDirConfig, myHead = ReadTarFileIndex(filePath, opts.Verbosity)

	if len(myWorkDirConfig.RepoPath) == 0 {
		myWorkDirConfig.RepoPath = myWorkDirConfig.Repos[myWorkDirConfig.DefaultRepo]
	}

	if len(opts.RepoName) > 0 {
		myWorkDirConfig.RepoPath = myWorkDirConfig.Repos[opts.RepoName]
	}

	if len(opts.RepoPath) > 0 {
		myWorkDirConfig.RepoPath = opts.RepoPath
	}

	if opts.Verbosity >= 2 {
		fmt.Printf("Repo config: %s\n", myWorkDirConfig.RepoPath)
	}

	myRepoConfig := ReadRepoConfigFile(path.Join(myWorkDirConfig.RepoPath, "config.toml"))
	
	if opts.Verbosity >= 1 {
		fmt.Println("Head:")
		fmt.Println(myHead)
		fmt.Printf("Branch: %s\nParent commit: %s\n", myHead.BranchName, myHead.CommitID)
	}

	if len(myHead.BranchName) == 0 {
		myHead.BranchName = "main"
	}		

	branchFolder := path.Join(myWorkDirConfig.RepoPath, "branches", myWorkDirConfig.WorkDirName)
	branchPath := path.Join(branchFolder, myHead.BranchName + ".toml")	
	myBranch := ReadBranch(branchPath)

	mySnapshot.ParentIDs = append([]string{myHead.CommitID}, parentIds...)

	chunkIDs, chunkPacks := PackFile(filePath, myWorkDirConfig.RepoPath, myRepoConfig.ChunkerPolynomial, opts.Verbosity)
	mySnapshot.ChunkIDs = chunkIDs

	snapshotFolder := path.Join(myWorkDirConfig.RepoPath, "snapshots", myWorkDirConfig.WorkDirName)
	snapshotBasename := fmt.Sprintf("%s", mySnapshot.ID[0:40])
	os.Mkdir(snapshotFolder, 0777)
	snapshotPath := path.Join(snapshotFolder, snapshotBasename + ".json")
	WriteSnapshot(snapshotPath, mySnapshot)

	// Do I really need to track commit id in head??
	myHead.CommitID = mySnapshot.ID
	myBranch.CommitID = mySnapshot.ID

	WriteBranch(branchPath, myBranch, opts.Verbosity)

	treeFolder := path.Join(myWorkDirConfig.RepoPath, "trees")
	treeBasename := mySnapshot.ID[0:40]
	os.Mkdir(treeFolder, 0777)
	treePath := path.Join(treeFolder, treeBasename+".json")
	WriteTree(treePath, chunkPacks)

	if opts.Verbosity >= 1 {
		fmt.Printf("%s", colorGreen)
		fmt.Printf("Created snapshot %s (%s)\n", mySnapshot.ID[0:16], mySnapshot.ID)
		fmt.Printf("%s", colorReset)
	} else {
		fmt.Println(mySnapshot.ID)
	}

	return myHead
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

func WriteBranch(branchPath string, myBranch Branch, verbosity int) {
	dir := filepath.Dir(branchPath)
	CreateFolder(dir, verbosity)
	f, err := os.Create(branchPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create branch file %s", branchPath))
	}

	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myBranch)
	f.Close()	
}

func ReadBranch(branchPath string) Branch {
	var myBranch Branch
	f, err := os.Open(branchPath)

	if err != nil {
		//panic(fmt.Sprintf("Error: Could not read head file %s", headPath))
		fmt.Printf("No branch file exists, returning default head struct\n")
		return Branch{}
	}

	if _, err := toml.DecodeReader(f, &myBranch); err != nil {
		panic(fmt.Sprintf("Error:could not decode branch file %s", branchPath))
	}

	f.Close()
	return myBranch
}

func GetFullSnapshotId(snapshotId string, cfg workDirConfig) string {
	snapshotPaths := ListSnapshots(cfg)

	for  _,  snapshotPath := range snapshotPaths {
		n := len(snapshotId) - 1
		sid := snapshotPath[n-SNAPSHOT_ID_LEN-5 : n-5]

		if len(sid) < len(snapshotId) {
			n = len(sid) - 1
		}

		if snapshotId[0:n] == sid[0:n] {
			snapshotId = sid
			break
		}
	}

	// TODO: return an error if no match
	return snapshotId
}

func ReadSnapshot(snapshot string, cfg workDirConfig) Commit {
	snapshotsFolder := path.Join(cfg.RepoPath, "snapshots", cfg.WorkDirName)
	snapshotPath := path.Join(snapshotsFolder, snapshot + ".json")
	return ReadSnapshotFile(snapshotPath)
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

func PrintSnapshots(cfg workDirConfig, snapshotId string, maxSnapshots int, opts Options) {
	// fmt.Printf("Verbosity = %d\n", verbosity)
	// print a specific revision
	snapshotCount := 0
	repoPath := cfg.RepoPath
	projectName := cfg.WorkDirName

	if maxSnapshots != 0 && opts.Verbosity >= 1 {
		fmt.Println("Snapshot History")
	}

	for  {	
		snapshotPath := filepath.Join(repoPath, "snapshots", projectName, snapshotId + ".json")
		mySnapshot := ReadSnapshotFile(snapshotPath)
		PrintSnapshot(mySnapshot, 0, opts)
		parents := mySnapshot.ParentIDs

		if len(parents) == 0 || len(parents[0]) == 0 {
			break
		} else {
			snapshotId = parents[0]
		}

		if maxSnapshots > 0 {
			snapshotCount++

			if snapshotCount >= maxSnapshots {
				break
			}
		}
	}
}

func PrintAllSnapshots(cfg workDirConfig, snapshot string, opts Options) {
	// fmt.Printf("Verbosity = %d\n", opts.Verbosity)
	snapshotPaths := ListSnapshots(cfg)
	// print a specific revision
	if len(snapshot) == 0 {
		if opts.Verbosity >= 1 {
			fmt.Println("Snapshot History")
		}

		for _, snapshotPath := range snapshotPaths {
			// fmt.Printf("Path: %s\n", snapshotPath)
			PrintSnapshot(ReadSnapshotFile(snapshotPath), 10, opts)
		}
	} else {
		if opts.Verbosity >= 1 {
			fmt.Println("Snapshot")
		}

		for _, snapshotPath := range snapshotPaths {
			// if i >= 1 {
			// 	fmt.Println("\n")
			// }

			n := len(snapshotPath)
			snapshotId := snapshotPath[n-SNAPSHOT_ID_LEN-5 : n-5]

			if snapshotId[0:8] == snapshot {
				PrintSnapshot(ReadSnapshotFile(snapshotPath), 0, opts)
			}
		}
	}
}

func PrintSnapshot(mySnapshot Commit, maxFiles int, opts Options) {
	if opts.Verbosity <= 0 {
		fmt.Printf("%s\n", mySnapshot.ID)
		return
	}

	if opts.Color {
		fmt.Printf("%s", colorGreen)
	}

	fmt.Printf("ID: %s (%s)", mySnapshot.ID[0:8], mySnapshot.ID)

	if opts.Color {
		fmt.Printf("%s", colorReset)
	}

	fmt.Printf("\n")
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

