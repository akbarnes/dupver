package dupver

import (
	"archive/tar"
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Commit struct {
	TarFileName string
	ID          string
	Message     string
	Time        string
	Files       []fileInfo
	ChunkIDs    []string
	ParentIDs   []string
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

func CommitFile(filePath string, msg string, verbosity int) string {
	var myWorkDirConfig workDirConfig
	t := time.Now()

	var mySnapshot Commit
	mySnapshot.ID = RandHexString(SNAPSHOT_ID_LEN)
	mySnapshot.Time = t.Format("2006/01/02 15:04:05")
	mySnapshot.TarFileName = filePath
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
	snapshotPath := path.Join(snapshotFolder, snapshotBasename+".json")
	WriteSnapshot(snapshotPath, mySnapshot)

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

func ReadTarFileIndex(filePath string, verbosity int) ([]fileInfo, workDirConfig) {
	tarFile, err := os.Open(filePath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Could not open input tar file %s when reading index", filePath))
	}

	files, myConfig := ReadTarIndex(tarFile, verbosity)
	tarFile.Close()

	return files, myConfig
}

func ReadTarIndex(tarFile *os.File, verbosity int) ([]fileInfo, workDirConfig) {
	files := []fileInfo{}
	var myConfig workDirConfig
	// var baseFolder string
	var configPath string
	maxFiles := 10

	if verbosity >= 1 {
		fmt.Println("Files:")
	}

	// Open and iterate through the files in the archive.
	tr := tar.NewReader(tarFile)
	i := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(fmt.Sprintf("Error processing section while reading tar file index"))
		}

		// if i == 0 {
		// 	baseFolder = hdr.Name
		// 	myConfig.WorkDirName = baseFolder
		// 	configPath = path.Join(baseFolder,".dupver","config.toml")
		// 	// fmt.Printf("Base folder: %s\nConfig path: %s\n", baseFolder, configPath)
		// }

		if strings.HasSuffix(hdr.Name, ".dupver/config.toml") {
			if verbosity >= 1 {
				fmt.Printf("Reading config file %s\n", hdr.Name)
			}

			if _, err := toml.DecodeReader(tr, &myConfig); err != nil {
				panic(fmt.Sprintf("Error decoding repo configuration file %s while reading tar file index", configPath))
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

		if i <= maxFiles && verbosity >= 1 {
			fmt.Printf("%2d: %s\n", i, hdr.Name)
		}

		files = append(files, myFileInfo)
	}

	if i > maxFiles && maxFiles > 0 && verbosity >= 1 {
		fmt.Printf("...\nSkipping %d more files\n", i-maxFiles)
	}

	return files, myConfig
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

func WriteTree(treePath string, chunkPacks map[string]string) {
	f, err := os.Create(treePath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create tree file %s", treePath))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(chunkPacks)
	f.Close()
}

func ReadTrees(repoPath string) map[string]string {
	treesGlob := path.Join(repoPath, "trees", "*.json")
	// fmt.Println(treesGlob)
	treePaths, err := filepath.Glob(treesGlob)

	if err != nil {
		panic(fmt.Sprintf("Error reading trees %s", treesGlob))
	}

	chunkPacks := make(map[string]string)

	for _, treePath := range treePaths {
		treePacks := make(map[string]string)

		f, err := os.Open(treePath)

		if err != nil {
			panic(fmt.Sprintf("Error: could not read tree file %s", treePath))
		}

		myDecoder := json.NewDecoder(f)

		if err := myDecoder.Decode(&treePacks); err != nil {
			panic(fmt.Sprintf("Error: could not decode tree file %s", treePath))
		}

		// TODO: handle supersedes to allow repacking files
		for k, v := range treePacks {
			chunkPacks[k] = v
		}

		f.Close()
	}

	return chunkPacks
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

func WorkDirStatus(workDir string, snapshot Commit, verbosity int) {
	workDirPrefix := ""

	if len(workDir) == 0 {
		workDir = "."
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		workDirPrefix = filepath.Base(cwd)
	}

	if verbosity >= 2 {
		fmt.Printf("Comparing changes for wd \"%s\" (prefix: \"%s\"\n", workDir, workDirPrefix)
	}

	myFileInfo := make(map[string]fileInfo)
	deletedFiles := make(map[string]bool)
	changes := false

	for _, fi := range snapshot.Files {
		myFileInfo[fi.Path] = fi
		deletedFiles[fi.Path] = true
	}

	var CompareAgainstSnapshot = func(curPath string, info os.FileInfo, err error) error {
		// fmt.Printf("Comparing path %s\n", path)
		if len(workDirPrefix) > 0 {
			curPath = path.Join(workDirPrefix, curPath)
		}

		curPath = strings.ReplaceAll(curPath, "\\", "/")

		if info.IsDir() {
			curPath += "/"
		}

		if snapshotInfo, ok := myFileInfo[curPath]; ok {
			deletedFiles[curPath] = false

			// fmt.Printf(" mtime: %s\n", snapshotInfo.ModTime)
			// t, err := time.Parse(snapshotInfo.ModTime, "2006/01/02 15:04:05")
			// check(err)

			if snapshotInfo.ModTime != info.ModTime().Format("2006/01/02 15:04:05") {
				if !info.IsDir() {
					fmt.Printf("%sM %s%s\n", colorCyan, curPath, colorReset)
					// fmt.Printf("M %s\n", curPath)
					changes = true
				}
			} else if verbosity >= 2 {
				fmt.Printf("%sU %s%s\n", colorWhite, curPath, colorReset)
			}
		} else {
			fmt.Printf("%s+ %s%s\n", colorGreen, curPath, colorReset)
			changes = true
		}

		return nil
	}

	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)

	filepath.Walk(workDir, CompareAgainstSnapshot)

	for file, deleted := range deletedFiles {
		if strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}

		if deleted {
			fmt.Printf("%s- %s%s\n", colorRed, file, colorReset)
			changes = true
		}
	}

	if !changes && verbosity >= 1 {
		fmt.Printf("No changes detected\n")
	}
}
