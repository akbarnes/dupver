package dupver

import (
	"encoding/json"
	"errors"
	"fmt"

	// "log"
	"archive/zip"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Commit struct {
	// TarFileName string
	ID        string
	Message   string
	Time      string
	ParentIDs []string
	Files     []fileInfo
	ChunkIDs  []string
}

type Head struct {
	BranchName string
	CommitID   string // use this for detached head, but do I need this?
	Branches   map[string]string
	CommitIDs  map[string]string
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
func CopySnapshot(snapshotId string, sourceRepoPath string, destRepoPath string, opts Options) {
	fmt.Printf("Copying snapshot %s: %s -> %s\n", snapshotId, sourceRepoPath, destRepoPath)
	sourceSnapshotsFolder := filepath.Join(sourceRepoPath, "snapshots", opts.WorkDirName)
	destSnapshotsFolder := filepath.Join(destRepoPath, "snapshots", opts.WorkDirName)
	os.Mkdir(destSnapshotsFolder, 0777)

	sourceSnapshotPath := filepath.Join(sourceSnapshotsFolder, snapshotId+".json")
	destSnapshotPath := filepath.Join(destSnapshotsFolder, snapshotId+".json")

	fmt.Printf("Copying %s -> %s\n", sourceSnapshotPath, destSnapshotPath)
	CopyFile(sourceSnapshotPath, destSnapshotPath) // TODO: check error status
	snapshot := ReadSnapshotFile(sourceSnapshotPath)
	chunkIndex := 0

	// TODO: Move this into CopyChunks in pack.go
	const maxPackSize int = 104857600 // 100 MB
	chunkIDs := []string{}
	sourceChunkPacks := ReadTrees(sourceRepoPath)
	// fmt.Printf("Source chunk packs:\n")
	// fmt.Println(sourceChunkPacks)
	destChunkPacks := ReadTrees(destRepoPath)
	newChunkPacks := make(map[string]string)
	var curPackSize int
	stillReadingInput := true

	totalDataSize := 0
	dupDataSize := 0

	newPackNum := 0
	totalChunkNum := 0
	dupChunkNum := 0

	for stillReadingInput {
		packId := RandHexString(PACK_ID_LEN)
		destPackFolderPath := path.Join(destRepoPath, "packs", packId[0:2])
		os.MkdirAll(destPackFolderPath, 0777)
		destPackPath := path.Join(destPackFolderPath, packId+".zip")

		newPackNum++

		if opts.Verbosity >= 2 {
			fmt.Printf("Creating pack file %3d: %s\n", newPackNum, destPackPath)
		} else if opts.Verbosity == 1 {
			fmt.Printf("Creating pack number: %3d, ID: %s\n", newPackNum, packId[0:16])
		}

		zipFile, err := os.Create(destPackPath)

		if err != nil {
			panic(fmt.Sprintf("Error creating zip file %s", destPackPath))
		}
		zipWriter := zip.NewWriter(zipFile)

		i := 0
		curPackSize = 0

		for curPackSize < maxPackSize { // white chunks to pack
			// chunk, err := mychunker.Next(buf)
			chunkId := snapshot.ChunkIDs[chunkIndex]
			chunk := LoadChunk(sourceRepoPath, chunkId, sourceChunkPacks, opts)
			chunkIndex++

			if chunkIndex >= len(snapshot.ChunkIDs) {
				// fmt.Printf("Reached end of input file, stop chunking\n")
				stillReadingInput = false
				break
			} else if err != nil {
				panic("Error chunking input file")
			}

			i++
			// chunkId := fmt.Sprintf("%064x", sha256.Sum256(chunk.Data))
			chunkIDs = append(chunkIDs, chunkId)

			totalDataSize += int(len(chunk))
			totalChunkNum++

			if _, ok := destChunkPacks[chunkId]; ok {
				if opts.Verbosity >= 2 {
					fmt.Printf("Skipping Chunk ID %s already in pack %s\n", chunkId[0:16], destChunkPacks[chunkId][0:16])
				}

				dupChunkNum++
				dupDataSize += int(len(chunk))
			} else {
				if opts.Verbosity >= 2 {
					fmt.Printf("Chunk %d: chunk size %d kB, total size %d kB, ", i, len(chunk)/1024, curPackSize/1024)
					fmt.Printf("chunk ID: %s\n", chunkId[0:16])
				}
				destChunkPacks[chunkId] = packId
				newChunkPacks[chunkId] = packId

				var header zip.FileHeader
				header.Name = chunkId
				header.Method = zip.Deflate

				writer, err := zipWriter.CreateHeader(&header)

				if err != nil {
					panic(fmt.Sprintf("Error creating zip file header for %s", destPackPath))
				}

				writer.Write(chunk)
				curPackSize += len(chunk)
			}
		}

		if opts.Verbosity >= 2 {
			if stillReadingInput {
				fmt.Printf("Pack size %d exceeds max size %d\n", curPackSize, maxPackSize)
			}

			fmt.Printf("Reached end of input, closing zip file\n")
		}

		zipWriter.Close()
		zipFile.Close()
	}

	if opts.Verbosity >= 1 {
		newChunkNum := totalChunkNum - dupChunkNum
		newDataSize := totalDataSize - dupDataSize

		newMb := float64(newDataSize) / 1e6
		dupMb := float64(dupDataSize) / 1e6
		totalMb := float64(totalDataSize) / 1e6

		fmt.Printf("%0.2f new, %0.2f duplicate, %0.2f total MB raw data stored\n", newMb, dupMb, totalMb)
		fmt.Printf("%d new, %d duplicate, %d total chunks\n", newChunkNum, dupChunkNum, totalChunkNum)
		fmt.Printf("%d packs stored, %0.2f chunks/pack\n", newPackNum, float64(newChunkNum)/float64(newPackNum))
	}

	treeFolder := path.Join(destRepoPath, "trees")
	treeBasename := snapshotId[0:40]
	os.Mkdir(treeFolder, 0777)
	treePath := path.Join(treeFolder, treeBasename+".json")
	WriteTree(treePath, destChunkPacks)

	if opts.Verbosity >= 1 {
		fmt.Printf("%s", colorGreen)
		fmt.Printf("Copied snapshot %s (%s)\n", snapshotId[0:16], snapshotId)
		fmt.Printf("%s", colorReset)
	} else {
		fmt.Println(snapshotId)
	}
}

func CommitFile(filePath string, parentIds []string, msg string, opts Options) {
	var myWorkDirConfig workDirConfig

	t := time.Now()

	var mySnapshot Commit
	// var myHead Head
	mySnapshot.ID = RandHexString(SNAPSHOT_ID_LEN)
	mySnapshot.Time = t.Format("2006/01/02 15:04:05")
	mySnapshot = UpdateMessage(mySnapshot, msg, filePath)
	mySnapshot.Files, myWorkDirConfig, _ = ReadTarFileIndex(filePath, opts.Verbosity)

	if len(opts.RepoName) == 0 {
		opts.RepoName = myWorkDirConfig.DefaultRepo
	}

	if len(opts.RepoPath) == 0 {
		opts.RepoPath = myWorkDirConfig.Repos[opts.RepoName]
	}

	if len(opts.WorkDirName) == 0 {
		opts.WorkDirName = myWorkDirConfig.WorkDirName
	}

	if len(opts.BranchName) == 0 {
		opts.BranchName = myWorkDirConfig.BranchName
	}

	myRepoConfig := ReadRepoConfigFile(path.Join(opts.RepoPath, "config.toml"))
	branchFolder := path.Join(opts.RepoPath, "branches", opts.WorkDirName)
	branchPath := path.Join(branchFolder, myWorkDirConfig.BranchName+".toml")
	myBranch := ReadBranch(branchPath)

	if opts.Verbosity >= 1 {
		fmt.Printf("Branch: %s\nParent commit: %s\n", opts.BranchName, myBranch.CommitID)
	}

	mySnapshot.ParentIDs = append([]string{myBranch.CommitID}, parentIds...)

	chunkIDs, chunkPacks := PackFile(filePath, opts.RepoPath, myRepoConfig.ChunkerPolynomial, opts.Verbosity)
	mySnapshot.ChunkIDs = chunkIDs

	snapshotFolder := path.Join(opts.RepoPath, "snapshots", opts.WorkDirName)
	snapshotBasename := fmt.Sprintf("%s", mySnapshot.ID[0:40])
	os.Mkdir(snapshotFolder, 0777)
	snapshotPath := path.Join(snapshotFolder, snapshotBasename+".json")
	WriteSnapshot(snapshotPath, mySnapshot)

	// Do I really need to track commit id in head??
	myBranch.CommitID = mySnapshot.ID

	WriteBranch(branchPath, myBranch, opts.Verbosity)

	treeFolder := path.Join(opts.RepoPath, "trees")
	treeBasename := mySnapshot.ID[0:40]
	os.Mkdir(treeFolder, 0777)
	treePath := path.Join(treeFolder, treeBasename+".json")
	WriteTree(treePath, chunkPacks)

	if opts.Verbosity >= 1 {
		if opts.Color {
			fmt.Printf("%s", colorGreen)
		}

		fmt.Printf("Created snapshot %s (%s)\n", mySnapshot.ID[0:16], mySnapshot.ID)

		if opts.Color {
			fmt.Printf("%s", colorReset)
		}
	} else {
		fmt.Println(mySnapshot.ID)
	}
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

func CreateTag(tagName string, snapshotId string, opts Options) {
	tagFolder := path.Join(opts.RepoPath, "tags", opts.WorkDirName)
	tagPath := path.Join(tagFolder, tagName+".toml")
	myTag := Branch{CommitID: snapshotId}

	if opts.Verbosity >= 1 {
		fmt.Printf("Tag commit: %s\n", snapshotId)
	}

	WriteBranch(tagPath, myTag, opts.Verbosity)
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
		fmt.Printf("Branch file %s does not exist, returning default head struct\n", branchPath)
		return Branch{}
	}

	if _, err := toml.DecodeReader(f, &myBranch); err != nil {
		panic(fmt.Sprintf("Error:could not decode branch file %s", branchPath))
	}

	f.Close()
	return myBranch
}

func GetFullSnapshotId(snapshotId string, opts Options) string {
	snapshotPaths := ListSnapshots(opts)

	for _, snapshotPath := range snapshotPaths {
		n := len(snapshotId) - 1
		sid := filepath.Base(snapshotPath)
		sid = sid[0 : len(sid)-5]
		fmt.Printf("path: %s\nsid: %s\n", snapshotPath, sid)

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

func ReadSnapshot(snapshot string, opts Options) Commit {
	snapshotsFolder := filepath.Join(opts.RepoPath, "snapshots", opts.WorkDirName)
	snapshotPath := filepath.Join(snapshotsFolder, snapshot+".json")

	if opts.Verbosity >= 2 {
		fmt.Printf("Snapshot path: %s\n", snapshotPath)
	}

	return ReadSnapshotFile(snapshotPath)
}

func ReadSnapshotFile(snapshotPath string) Commit {
	var mySnapshot Commit
	f, err := os.Open(snapshotPath)

	if err != nil {
		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
		return Commit{}
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&mySnapshot); err != nil {
		panic(fmt.Sprintf("Error:could not decode snapshot file %s", snapshotPath))
	}

	f.Close()
	return mySnapshot
}

func ReadSnapshotId(snapshotId string, opts Options) (Commit, error) {
	snapshotPaths := ListSnapshots(opts)

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

func ListSnapshots(opts Options) []string {
	snapshotsFolder := path.Join(opts.RepoPath, "snapshots", opts.WorkDirName)
	snapshotGlob := path.Join(snapshotsFolder, "*.json")
	// fmt.Println(snapshotGlob)
	snapshotPaths, err := filepath.Glob(snapshotGlob)

	if err != nil {
		panic(fmt.Sprintf("Error listing snapshots glob %s", snapshotGlob))
	}
	return snapshotPaths
}

func PrintSnapshots(snapshotId string, maxSnapshots int, opts Options) {
	// fmt.Printf("Verbosity = %d\n", verbosity)
	// print a specific revision
	snapshotCount := 0
	repoPath := opts.RepoPath
	projectName := opts.WorkDirName

	if maxSnapshots != 0 && opts.Verbosity >= 1 {
		fmt.Println("Snapshot History")
	}

	for {
		snapshotPath := filepath.Join(repoPath, "snapshots", projectName, snapshotId+".json")

		if opts.Verbosity >= 2 {
			fmt.Printf("Snapshot path: %s\n\n", snapshotPath)
		}

		mySnapshot := ReadSnapshotFile(snapshotPath)

		if len(mySnapshot.ID) == 0 {
			fmt.Println("No snapshots")
			break
		}

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

func PrintAllSnapshots(snapshot string, opts Options) {
	// fmt.Printf("Verbosity = %d\n", opts.Verbosity)
	snapshotPaths := ListSnapshots(opts)
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
}
