package dupver

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	// "path"
	"path/filepath"
	"sort"
	"strings"
	"time"
	// "log"

	"github.com/BurntSushi/toml"
	"github.com/akbarnes/dupver/src/fancyprint"
)

type Commit struct {
	ID        string
	Branch    string
	Message   string
	Time      string
	ParentIDs []string
	Files     []fileInfo
	ChunkIDs  []string
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

// Copy a snapshot given a snapshot ID, source repo path and dest repo path
func CopySnapshot(snapshotId string, sourceRepoPath string, destRepoPath string, opts Options) {
	// fancyprint.Noticef("Copying snapshot %s: %s -> %s\n", snapshotId, sourceRepoPath, destRepoPath)
	sourceSnapshotsFolder := filepath.Join(sourceRepoPath, "snapshots", opts.WorkDirName)
	destSnapshotsFolder := filepath.Join(destRepoPath, "snapshots", opts.WorkDirName)
	os.Mkdir(destSnapshotsFolder, 0777)

	sourceSnapshotPath := filepath.Join(sourceSnapshotsFolder, snapshotId+".json")
	destSnapshotPath := filepath.Join(destSnapshotsFolder, snapshotId+".json")

	fancyprint.Noticef("Copying %s\n  to %s\n", sourceSnapshotPath, destSnapshotPath)
	CopyFile(sourceSnapshotPath, destSnapshotPath) // TODO: check error status
	snapshot := ReadSnapshotFile(sourceSnapshotPath)
	chunkIndex := 0

	// TODO: Move this into CopyChunks in pack.go
	// TODO: Move this into repo configuration
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
		destPackFolderPath := filepath.Join(destRepoPath, "packs", packId[0:2])
		os.MkdirAll(destPackFolderPath, 0777)
		destPackPath := filepath.Join(destPackFolderPath, packId+".zip")

		newPackNum++

		fancyprint.Infof("Creating pack file %3d: %s\n", newPackNum, destPackPath)
		if fancyprint.Verbosity <= fancyprint.NoticeLevel {
			fancyprint.Noticef("Creating pack number: %3d, ID: %s\n", newPackNum, packId[0:16])
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

			i++
			// chunkId := fmt.Sprintf("%064x", sha256.Sum256(chunk.Data))
			chunkIDs = append(chunkIDs, chunkId)

			totalDataSize += int(len(chunk))
			totalChunkNum++

			if _, ok := destChunkPacks[chunkId]; ok {
				fancyprint.Infof("Skipping Chunk ID %s already in pack %s\n", chunkId[0:16], destChunkPacks[chunkId][0:16])
				dupChunkNum++
				dupDataSize += int(len(chunk))
			} else {
				fancyprint.Infof("Chunk %d: chunk size %d kB, total size %d kB, ", i, len(chunk)/1024, curPackSize/1024)
				fancyprint.Infof("chunk ID: %s\n", chunkId[0:16])
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

			chunkIndex++

			if chunkIndex >= len(snapshot.ChunkIDs) {
				fancyprint.Infof("Reached end of input file, stop chunking\n")
				stillReadingInput = false
				break
			} 
		}

		if stillReadingInput {
			fancyprint.Infof("Pack size %d exceeds max size %d\n", curPackSize, maxPackSize)
		}

		fancyprint.Info("Reached end of input, closing zip file")

		// TODO: delete or don't create zip file if there weren't any files stored
		zipWriter.Close()
		zipFile.Close()

		if curPackSize == 0 {
			fancyprint.Infof("No new chunks stored so deleting %s\n", destPackPath )
			os.Remove(destPackPath)
		}
	}

	newChunkNum := totalChunkNum - dupChunkNum
	newDataSize := totalDataSize - dupDataSize

	newMb := float64(newDataSize) / 1e6
	dupMb := float64(dupDataSize) / 1e6
	totalMb := float64(totalDataSize) / 1e6

	fancyprint.Noticef("%0.2f new, %0.2f duplicate, %0.2f total MB raw data stored\n", newMb, dupMb, totalMb)
	fancyprint.Noticef("%d new, %d duplicate, %d total chunks\n", newChunkNum, dupChunkNum, totalChunkNum)
	fancyprint.Noticef("%d packs stored, %0.2f chunks/pack\n", newPackNum, float64(newChunkNum)/float64(newPackNum))

	treeFolder := filepath.Join(destRepoPath, "trees")
	treeBasename := snapshotId[0:40]
	os.Mkdir(treeFolder, 0777)
	treePath := filepath.Join(treeFolder, treeBasename+".json")
	WriteTree(treePath, destChunkPacks)

	if fancyprint.Verbosity >= fancyprint.NoticeLevel {
		fancyprint.SetColor(fancyprint.ColorGreen)
		fmt.Printf("Copied snapshot %s (%s)\n", snapshotId[0:16], snapshotId)
		fancyprint.ResetColor()
	} else {
		fmt.Println(snapshotId)
	}
}

// Commit a tar file into the repository. Project working directory name,
// branch and repository path are specified in the .dupver/config.toml
// file within the tar file
func CommitFile(filePath string, parentIds []string, msg string, JsonOutput bool, opts Options) Commit {
	var myWorkDirConfig workDirConfig

	t := time.Now()

	var snap Commit
	// var myHead Head
	snap.ID = RandHexString(SNAPSHOT_ID_LEN)
	snap.Time = t.Format("2006/01/02 15:04:05")
	snap = UpdateMessage(snap, msg, filePath)
	snap.Files, myWorkDirConfig = ReadTarFileIndex(filePath)

	if len(opts.RepoName) == 0 {
		opts.RepoName = myWorkDirConfig.DefaultRepo
	}

	if len(opts.RepoPath) == 0 {
		opts.RepoPath = myWorkDirConfig.Repos[opts.RepoName]
	}

	if len(opts.WorkDirName) == 0 {
		opts.WorkDirName = myWorkDirConfig.WorkDirName
	}

	if len(opts.Branch) == 0 {
		opts.Branch = myWorkDirConfig.Branch
	}

	snap.Branch = opts.Branch

	myRepoConfig := ReadRepoConfigFile(filepath.Join(opts.RepoPath, "config.toml"))
	branchFolder := filepath.Join(opts.RepoPath, "branches", opts.WorkDirName)
	branchPath := filepath.Join(branchFolder, myWorkDirConfig.Branch+".toml")
	myBranch := ReadBranch(branchPath)

	fancyprint.Infof("Branch: %s\nParent commit: %s\n", opts.Branch, myBranch.CommitID)
	snap.ParentIDs = append([]string{myBranch.CommitID}, parentIds...)

	chunkIDs, chunkPacks := PackFile(filePath, opts.RepoPath, myRepoConfig.ChunkerPolynomial, myRepoConfig.CompressionLevel)
	snap.ChunkIDs = chunkIDs

	snapshotFolder := filepath.Join(opts.RepoPath, "snapshots", opts.WorkDirName)
	snapshotBasename := fmt.Sprintf("%s", snap.ID[0:40])
	os.Mkdir(snapshotFolder, 0777)
	snapshotPath := filepath.Join(snapshotFolder, snapshotBasename+".json")
	WriteSnapshot(snapshotPath, snap)

	// Do I really need to track commit id in head??
	myBranch.CommitID = snap.ID

	WriteBranch(branchPath, myBranch)

	treeFolder := filepath.Join(opts.RepoPath, "trees")
	treeBasename := snap.ID[0:40]
	os.Mkdir(treeFolder, 0777)
	treePath := filepath.Join(treeFolder, treeBasename+".json")
	WriteTree(treePath, chunkPacks)

	if JsonOutput {
		PrintJson(snap.ID)
	} else if fancyprint.Verbosity >= fancyprint.NoticeLevel {
		fancyprint.SetColor(fancyprint.ColorGreen)
		fmt.Printf("Created snapshot %s (%s)\n", snap.ID[0:16], snap.ID)
		fancyprint.ResetColor()
	} else {
		fmt.Println(snap.ID)
	}

	return snap
}

// Remove PowerShell artifact of leading .\ in commit messages
func UpdateMessage(mySnapshot Commit, msg string, filePath string) Commit {
	if len(msg) == 0 {
		msg = strings.Replace(filePath[0:len(filePath)-4], ".\\", "", -1)
	}

	mySnapshot.Message = msg
	return mySnapshot
}

// Write a snapshot to disk given a file path
// TODO: Change this to WriteSnapshotFile?
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

// Create a pointer-style tag given tag name and snapshot ID
// Repo path is specified in options structure
func CreateTag(tagName string, snapshotId string, opts Options) {
	tagFolder := filepath.Join(opts.RepoPath, "tags", opts.WorkDirName)
	tagPath := filepath.Join(tagFolder, tagName+".toml")
	myTag := Branch{CommitID: snapshotId}

	fancyprint.Noticef("Tag commit: %s\n", snapshotId)
	WriteBranch(tagPath, myTag)
}

// Save the current branch head to a file
// TODO: Update this to use opts structure
// TODO: Change this to WriteBranchFile?
// TODO: Change this to take in a file stream? - Probably not, why would I need to?
func WriteBranch(branchPath string, myBranch Branch) {
	dir := filepath.Dir(branchPath)
	CreateFolder(dir)
	f, err := os.Create(branchPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create branch file %s", branchPath))
	}

	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myBranch)
	f.Close()
}

// Read the current branch head from a file
// TODO: Add some better error handling rather than the printf on line 346
func ReadBranch(branchPath string) Branch {
	var myBranch Branch
	f, err := os.Open(branchPath)

	if err != nil {
		//panic(fmt.Sprintf("Error: Could not read head file %s", headPath))
		fancyprint.Warnf("Branch file %s does not exist, returning default head struct\n", branchPath)
		return Branch{}
	}

	if _, err := toml.DecodeReader(f, &myBranch); err != nil {
		panic(fmt.Sprintf("Error:could not decode branch file %s", branchPath))
	}

	f.Close()
	return myBranch
}

// Given a partial snapshot ID, return the full snapshot ID
// by looking through the snapshots for a project
// TODO: return an error if no match
func GetFullSnapshotId(snapshotId string, opts Options) string {
	snapshotPaths := ListSnapshots(opts)

	for _, snapshotPath := range snapshotPaths {
		n := len(snapshotId) - 1
		sid := filepath.Base(snapshotPath)
		sid = sid[0 : len(sid)-5]
		// fmt.Printf("path: %s\nsid: %s\n", snapshotPath, sid)

		if len(sid) < len(snapshotId) {
			n = len(sid) - 1
		}

		if snapshotId[0:n] == sid[0:n] {
			snapshotId = sid
			break
		}
	}

	return snapshotId
}

// Given a partial snapshot ID, return the full snapshot ID
// by looking through the snapshots for a project
// TODO: return an error if no match
func (wd WorkDir) GetFullSnapshotId(snapshotId string) string {
	snapshotPaths := wd.ListSnapshots()

	for _, snapshotPath := range snapshotPaths {
		n := len(snapshotId) - 1
		sid := filepath.Base(snapshotPath)
		sid = sid[0 : len(sid)-5]
		// fmt.Printf("path: %s\nsid: %s\n", snapshotPath, sid)

		if len(sid) < len(snapshotId) {
			n = len(sid) - 1
		}

		if snapshotId[0:n] == sid[0:n] {
			snapshotId = sid
			break
		}
	}

	return snapshotId
}

// Read a snapshot given a full snapshot ID
func ReadSnapshot(snapshot string, opts Options) Commit {
	snapshotsFolder := filepath.Join(opts.RepoPath, "snapshots", opts.WorkDirName)
	snapshotPath := filepath.Join(snapshotsFolder, snapshot+".json")
	fancyprint.Debugf("Snapshot path: %s\n", snapshotPath)
	return ReadSnapshotFile(snapshotPath)
}

// Read a snapshot given a file path
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

// Read a snapshot given a partial snapshot ID
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

// Return a list of the snapshot files for a given repository and project
func ListSnapshots(opts Options) []string {
	snapshotsFolder := filepath.Join(opts.RepoPath, "snapshots", opts.WorkDirName)
	snapshotGlob := filepath.Join(snapshotsFolder, "*.json")
	// fmt.Println(snapshotGlob)
	snapshotPaths, err := filepath.Glob(snapshotGlob)

	if err != nil {
		panic(fmt.Sprintf("Error listing snapshots glob %s", snapshotGlob))
	}
	return snapshotPaths
}

// Return a list of the snapshot files for a given repository and project
func (wd WorkDir) ListSnapshots() []string {
	snapshotsFolder := filepath.Join(wd.Repo.Path, "snapshots", wd.ProjectName)
	snapshotGlob := filepath.Join(snapshotsFolder, "*.json")
	// fmt.Println(snapshotGlob)
	snapshotPaths, err := filepath.Glob(snapshotGlob)

	if err != nil {
		panic(fmt.Sprintf("Error listing snapshots glob %s", snapshotGlob))
	}
	return snapshotPaths
}

// Return the most recent snapshot structure for the current project
func LastSnapshot(opts Options) (Commit, error) {
	repoPath := opts.RepoPath
	projectName := opts.WorkDirName

	snapshotGlob := filepath.Join(repoPath, "snapshots", projectName, "*.json")
	snapshotPaths, _ := filepath.Glob(snapshotGlob)

	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	branch := opts.Branch

	// TODO: sort the snapshots by date
	for _, snapshotPath := range snapshotPaths {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		snap := ReadSnapshotFile(snapshotPath)

		if len(branch) == 0 || len(branch) > 0 && branch == snap.Branch {
			snapshotsByDate[snap.Time] = snap
			snapshotDates = append(snapshotDates, snap.Time)
		}

	}

	sort.Strings(snapshotDates)

	if len(snapshotDates) == 0 {
		return Commit{}, errors.New("no snapshots")
	}

	return snapshotsByDate[snapshotDates[len(snapshotDates)-1]], nil
}

// Print snapshots sorted in ascending order by date
// TODO: change the name to PrintSnapshotsByDate?
func PrintSnapshots(snapshotId string, maxSnapshots int, opts Options) {
	repoPath := opts.RepoPath
	projectName := opts.WorkDirName

	if len(snapshotId) > 0 {
		snap := ReadSnapshot(snapshotId, opts)
		PrintSnapshotFiles(snap, 0, opts)
		return
	}

	if maxSnapshots > 0 {
		fancyprint.Notice("Snapshot History")
	}

	snapshotGlob := filepath.Join(repoPath, "snapshots", projectName, "*.json")
	snapshotPaths, _ := filepath.Glob(snapshotGlob)

	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	// TODO: sort the snapshots by date
	for _, snapshotPath := range snapshotPaths {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		mySnapshot := ReadSnapshotFile(snapshotPath)
		snapshotsByDate[mySnapshot.Time] = mySnapshot
		snapshotDates = append(snapshotDates, mySnapshot.Time)
	}

	sort.Strings(snapshotDates)

	for i, sdate := range snapshotDates {
		snap := snapshotsByDate[sdate]
		b := opts.Branch

		if len(b) == 0 || len(b) > 0 && b == snap.Branch {
			PrintSnapshot(snap, 0, opts)
		}

		if maxSnapshots > 0 {
			if i >= maxSnapshots {
				break
			}
		}
	}
}

// Print snapshots sorted in ascending order by date
// TODO: change the name to PrintSnapshotsByDate?
func (workDir WorkDir) PrintSnapshots(snapshotId string) {
	repoPath := opts.RepoPath
	projectName := opts.WorkDirName

	if len(snapshotId) > 0 {
		snap := ReadSnapshot(snapshotId, opts)
		PrintSnapshotFiles(snap, 0, opts)
		return
	}


	snapshotGlob := filepath.Join(repoPath, "snapshots", projectName, "*.json")
	snapshotPaths, _ := filepath.Glob(snapshotGlob)

	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	// TODO: sort the snapshots by date
	for _, snapshotPath := range snapshotPaths {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		mySnapshot := ReadSnapshotFile(snapshotPath)
		snapshotsByDate[mySnapshot.Time] = mySnapshot
		snapshotDates = append(snapshotDates, mySnapshot.Time)
	}

	sort.Strings(snapshotDates)

	for i, sdate := range snapshotDates {
		snap := snapshotsByDate[sdate]
		b := opts.Branch

		if len(b) == 0 || len(b) > 0 && b == snap.Branch {
			PrintSnapshot(snap, 0, opts)
		}

		if maxSnapshots > 0 {
			if i >= maxSnapshots {
				break
			}
		}
	}
}


// Print snapshots as JSON in sorted in ascending order by date
// TODO: change the name to PrintSnapshotsByDate?
// TODO: add a sort option to ListSnapshots?
func PrintSnapshotsAsJson(snapshotId string, maxSnapshots int, snapshotFiles bool, opts Options) {
	type CommitPrint struct {
		ID      string
		Branch  string
		Message string
		Time    string
		Files   []fileInfo
	}

	repoPath := opts.RepoPath
	projectName := opts.WorkDirName

	if len(snapshotId) > 0 {
		snap := ReadSnapshot(snapshotId, opts)
		PrintJson(snap.Files)
		return
	}

	snapshotGlob := filepath.Join(repoPath, "snapshots", projectName, "*.json")
	snapshotPaths, _ := filepath.Glob(snapshotGlob)

	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	// TODO: sort the snapshots by date
	for _, snapshotPath := range snapshotPaths {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		mySnapshot := ReadSnapshotFile(snapshotPath)
		snapshotsByDate[mySnapshot.Time] = mySnapshot
		snapshotDates = append(snapshotDates, mySnapshot.Time)
	}

	sort.Strings(snapshotDates)
	printSnaps := []CommitPrint{}

	for _, sdate := range snapshotDates {
		snap := snapshotsByDate[sdate]

		ps := CommitPrint{}
		ps.ID = snap.ID
		ps.Branch = snap.Branch
		ps.Message = snap.Message
		ps.Time = snap.Time

		if snapshotFiles {
			ps.Files = snap.Files
		}

		printSnaps = append(printSnaps, ps)
	}

	PrintJson(printSnaps)
}

// Print snapshots without sorting
// TODO: Check if this is redundant
func PrintAllSnapshots(snapshotId string, opts Options) {
	// fmt.Printf("Verbosity = %d\n", opts.Verbosity)
	// print a specific revision

	fancyprint.Noticef("Branch: %s\n", opts.Branch)

	if len(snapshotId) == 0 {
		fancyprint.Notice("Snapshot History")

		for _, snapshotPath := range ListSnapshots(opts) {
			// fmt.Printf("Path: %s\n", snapshotPath)
			PrintSnapshot(ReadSnapshotFile(snapshotPath), 10, opts)
		}
	} else {
		if fancyprint.Verbosity >= fancyprint.NoticeLevel {
			fmt.Println("Snapshot")
		}

		snap := ReadSnapshot(snapshotId, opts)
		// PrintSnapshot(snap, 0, opts)
		PrintSnapshotFiles(snap, 0, opts)
	}
}

// Print a snapshot structure
func PrintSnapshot(mySnapshot Commit, maxFiles int, opts Options) {
	if fancyprint.Verbosity <= fancyprint.WarningLevel {
		fmt.Printf("%s %s %s\n", mySnapshot.ID, mySnapshot.Time, mySnapshot.Message)
		return
	}

	fancyprint.SetColor(fancyprint.ColorGreen)
	fmt.Printf("ID: %s (%s)", mySnapshot.ID[0:8], mySnapshot.ID)
	fancyprint.ResetColor()

	fmt.Printf("\n")
	fmt.Printf("Time: %s\n", mySnapshot.Time)

	if len(mySnapshot.Message) > 0 {
		fmt.Printf("Message: %s\n", mySnapshot.Message)
	}

	fmt.Printf("\n")
}

// "Files": [
//     {
//       "Path": "Arduino/",
//       "ModTime": "2019/06/27 05:25:16",
//       "Size": 0,
//       "Hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
//     },

// Print the list of files stored in a snapshot
func PrintSnapshotFiles(mySnapshot Commit, maxFiles int, opts Options) {
	for i, file := range mySnapshot.Files {
		if fancyprint.Verbosity <= fancyprint.WarningLevel {
			fmt.Printf("%s\n%d\n%s\n\n", file.ModTime, file.Size, file.Path)
		} else {
			fmt.Printf("%s ", file.ModTime)

			if file.Size >= 1e9 {
				fmt.Printf("%5.1f GB ", float64(file.Size)/1e9)
			} else if file.Size >= 1e6 {
				fmt.Printf("%5.1f MB ", float64(file.Size)/1e6)
			} else if file.Size >= 1e3 {
				fmt.Printf("%5.1f kB ", float64(file.Size)/1e3)
			} else {
				fmt.Printf("%5d B  ", file.Size)
			}

			fmt.Printf("%s\n", file.Path)
		}

		if maxFiles > 0 && i >= maxFiles {
			break
		}
	}
}
