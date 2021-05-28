package dupver

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sort"
	"time"

	// "io"
	// "bufio"
	// "crypto/sha256"
	// "encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/restic/chunker"

	"github.com/akbarnes/dupver/src/fancyprint"
)

type WorkDirConfig struct {
	Version 	int
	WorkDirName string
	Branch      string
	DefaultRepo string
	Repos       map[string]string
}

type WorkDir struct {
	ProjectName string
	Path		string
	Branch      string
	Repo 		Repo
}

// Create a valid project name given a folder name
func FolderToWorkDirName(folder string) string {
	return strings.ReplaceAll(strings.ToLower(folder), " ", "-")
}

// Initialize a project working directory configuration
// given the working directory path and project name
func InitWorkDir(workDirFolder string, workDirName string, opts Options) {
	var configPath string
	repoName := opts.RepoName
	repoPath := opts.RepoPath
	branch := opts.Branch

	fancyprint.Noticef("Workdir %s, name %s, repo %s\n", workDirFolder, workDirName, opts.RepoPath)

	if len(workDirFolder) == 0 {
		CreateFolder(".dupver")
		configPath = filepath.Join(".dupver", "config.toml")
	} else {
		CreateSubFolder(workDirFolder, ".dupver")
		configPath = filepath.Join(workDirFolder, ".dupver", "config.toml")
	}

	fancyprint.Infof("Writing workdir config file to: " + configPath)

	if len(workDirName) == 0 || workDirName == "." {
		if len(workDirFolder) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			// _, folder := path.Split(dir)
			folder := filepath.Base(dir)

			fancyprint.Debugf("Resolving folder %s to %s\n", dir, folder)
			workDirName = FolderToWorkDirName(folder)
		} else {
			workDirName = FolderToWorkDirName(workDirFolder)
		}

		if workDirName == "." || workDirName == fmt.Sprintf("%c", filepath.Separator) {
			log.Fatal("Invalid project name: " + workDirName)
		}

		fancyprint.Noticef("Workdir name not specified, setting to %s\n", workDirName)
	}

	if len(repoPath) == 0 {
		repoPath = filepath.Join(GetHome(), ".dupver_repo")
		fancyprint.Noticef("Repo path not specified, setting to %s\n", repoPath)
	}

	if fancyprint.Verbosity >= fancyprint.NoticeLevel {
		fmt.Printf("Repo name: [%s]\n", repoName)
	} else {
		fmt.Println(workDirName)
	}

	var myConfig WorkDirConfig
	// need to pass this as a parameter
	myConfig.DefaultRepo = repoName

	// TODO: specify an arbitrary branch
	myRepos := make(map[string]string)
	myRepos[repoName] = repoPath
	myConfig.Repos = myRepos
	myConfig.Branch = branch
	myConfig.WorkDirName = workDirName
	myConfig.SaveFile(configPath, false)
}

// Print the project working directory configuration
func (cfg WorkDirConfig) Print() {
	// WorkDirName = "admin"
	// Branch = "test"
	// DefaultRepo = "store"

	// [Repos]
	//   main = "C:\\Users\\305232/.dupver_repo"

	fmt.Printf("Working directory name: %s\n", cfg.WorkDirName)
	fmt.Printf("Current branch: %s\n\n", cfg.Branch)
	cfg.PrintRepos()
}

func (cfg WorkDirConfig) PrintRepos() {
	for name, path := range cfg.Repos {
		repoCfg, err := ReadRepoConfig(path)

		fmt.Printf("%s: %s", name, path)
		
		if err != nil {
			fancyprint.SetColor(fancyprint.ColorRed)
			fmt.Println(" Unreadable")
			fancyprint.ResetColor()
			continue
		}

		if repoCfg.CompressionLevel == 0 {
			fmt.Print(" Store (0)")
		} else {
			fmt.Printf(" Deflate (%d)", repoCfg.CompressionLevel)
		}

		fmt.Printf(" %d", repoCfg.ChunkerPolynomial)

		if name == cfg.DefaultRepo {
			fancyprint.SetColor(fancyprint.ColorGreen)
			fmt.Print(" default")
			fancyprint.ResetColor()
		}

		fmt.Println("")
	}
}

func (cfg WorkDirConfig) PrintReposAsJson() {
	type repoConfigPrint struct {
		Name              string
		Path              string
		Default           bool
		Version           int
		ChunkerPolynomial chunker.Pol
		CompressionLevel  uint16
	}

	repoConfigs := []repoConfigPrint{}

	for name, path := range cfg.Repos {
		repoCfg, err := ReadRepoConfig(path)
		Check(err)

		rc := repoConfigPrint{Name: name, Path: path, Default: false}

		if name == cfg.DefaultRepo {
			rc.Default = true
		}

		rc.Version = repoCfg.Version
		rc.ChunkerPolynomial = repoCfg.ChunkerPolynomial
		rc.CompressionLevel = repoCfg.CompressionLevel

		repoConfigs = append(repoConfigs, rc)
	}

	PrintJson(repoConfigs)
}

func (cfg WorkDirConfig) PrintReposConfig() {
	for name, path := range cfg.Repos {
		repoCfg, err := ReadRepoConfig(path)
		Check(err)
		fmt.Printf("%s: %s", name, path)

		if repoCfg.CompressionLevel == 0 {
			fmt.Print(" Store (0)")
		} else {
			fmt.Printf(" Deflate (%d)", repoCfg.CompressionLevel)
		}

		fmt.Printf(" %d", repoCfg.ChunkerPolynomial)

		if name == cfg.DefaultRepo {
			fancyprint.SetColor(fancyprint.ColorGreen)
			fmt.Print(" default")
			fancyprint.ResetColor()
		}

		fmt.Println("")
	}
}

func (cfg WorkDirConfig) PrintJson() {
	type workDirConfigPrint struct {
		WorkDirName string
		Branch      string
		DefaultRepo string
	}

	wc := workDirConfigPrint{}
	wc.WorkDirName = cfg.WorkDirName
	wc.Branch = cfg.Branch
	wc.DefaultRepo = cfg.DefaultRepo
	PrintJson(wc)
}

// Add a new repository to the working directory configuration
func AddRepoToWorkDir(workDirPath string, repoName string, repoPath string, makeDefaultRepo bool, jsonOutput bool) {
	cfg, err := ReadWorkDirConfig(workDirPath)

	if err != nil {
		// Todo: handle invalid configuration file
		fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
		os.Exit(0)
	}

	cfg.Repos[repoName] = repoPath

	if makeDefaultRepo {
		cfg.DefaultRepo = repoName
	}

	if jsonOutput {
		PrintJson(cfg)
	}

	cfg.Save(workDirPath, true)
}

// List the repositories in the working directory configuration
func ListWorkDirRepos(workDirPath string) {
	cfg, err := ReadWorkDirConfig(workDirPath)
	maxLen := 0

	if err != nil {
		// Todo: handle invalid configuration file
		fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
		os.Exit(0)
	}

	for name, _ := range cfg.Repos {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	fmtStr := "%" + strconv.Itoa(maxLen) + "s: %s\n"

	for name, path := range cfg.Repos {
		if fancyprint.Verbosity >= fancyprint.NoticeLevel {
			fmt.Printf(fmtStr, name, path)
		} else {
			fmt.Printf("%s %s\n", name, path)
		}
	}
}

// List the repositories in the working directory configuration as JSON
func ListWorkDirReposAsJson(workDirPath string) {
	type RepoListing struct {
		Name              string
		Path              string
		Default           bool
		ChunkerPolynomial chunker.Pol
		CompressionLevel  uint16
	}

	repoListings := []RepoListing{}
	cfg, err := ReadWorkDirConfig(workDirPath)

	if err != nil {
		// Todo: handle invalid configuration file
		fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
		os.Exit(0)
	}

	for name, path := range cfg.Repos {
		rl := RepoListing{Name: name, Path: path, Default: false}

		if name == cfg.DefaultRepo {
			rl.Default = true
		}

		repoCfg, err := ReadRepoConfig(path)
		Check(err)
		rl.ChunkerPolynomial = repoCfg.ChunkerPolynomial
		rl.CompressionLevel = repoCfg.CompressionLevel
		repoListings = append(repoListings, rl)
	}

	PrintJson(repoListings)
}

// Change the project name in the working directory configuration
func UpdateWorkDirName(myWorkDirConfig WorkDirConfig, workDirName string) WorkDirConfig {
	if len(workDirName) > 0 {
		myWorkDirConfig.WorkDirName = workDirName
	}

	return myWorkDirConfig
}

// Load a project working directory configuration given
// the working directory path
func ReadWorkDirConfig(workDir string) (WorkDirConfig, error) {
	var configPath string

	if len(workDir) == 0 {
		configPath = filepath.Join(".dupver", "config.toml")
	} else {
		configPath = filepath.Join(workDir, ".dupver", "config.toml")
	}

	return ReadWorkDirConfigFile(configPath)
}

// Load a project working directory configuration given
// the project working directory configuration file path
func ReadWorkDirConfigFile(filePath string) (WorkDirConfig, error) {
	var cfg WorkDirConfig

	f, err := os.Open(filePath)

	if err != nil {
		return WorkDirConfig{}, errors.New("config file missing")
	}

	if _, err = toml.DecodeReader(f, &cfg); err != nil {
		panic(fmt.Sprintf("Invalid configuration file: %s\n", filePath))
	}

	f.Close()

	// TODO: store file version globally
	if cfg.Version != 2 {
		return cfg, errors.New("wrong version of config file")
	}

	return cfg, nil
}

// Load a project working directory configuration given
// the project working directory configuration file path
func InstantiateWorkDir(cfg WorkDirConfig) (WorkDir) {
	wdPath, err := os.Getwd()
	Check(err)
	wd := WorkDir{ProjectName: cfg.WorkDirName, Branch: cfg.Branch, Path: wdPath}
	wd.Repo = LoadRepo(cfg.Repos[cfg.DefaultRepo])
	return wd
}

func LoadWorkDir(workDirPath string) (WorkDir, error) {
	cfg, err := ReadWorkDirConfig(workDirPath)

	if err != nil {
		return WorkDir{}, err
	}

	wd := InstantiateWorkDir(cfg)
	wd.Path = workDirPath
	return wd, nil
}

// funct LoadWorkDir(workDir String)

// Save a project working directory configuration given
// the working directory path
func (cfg WorkDirConfig) Save(workDir string, forceWrite bool) {
	configPath := filepath.Join(".dupver", "config.toml")

	if len(workDir) > 0 {
		configPath = filepath.Join(workDir, ".dupver", "config.toml")
	} 

	cfg.SaveFile(configPath, forceWrite)
}

// Save a project working directory configuration given
// the project working directory configuration file path
func (cfg WorkDirConfig) SaveFile(configPath string, forceWrite bool) {
	if _, err := os.Stat(configPath); err == nil && !forceWrite {
		// panic("Refusing to write existing project workdir config " + configPath)
		panic(fmt.Sprintf("Refusing to write existing project workdir config: %s\n", configPath))
	}

	fancyprint.Infof("Writing config:\n%+v\n", cfg)
	fancyprint.Infof("to: %s\n", configPath)

	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(cfg)
	f.Close()
}

// Compare the status of files in a working directory
// against a snapshot
func (wd WorkDir) PrintStatus(snapshot Commit) {
	workDirPrefix := ""

	if len(wd.Path) == 0 {
		wd.Path = "."
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		workDirPrefix = filepath.Base(cwd)
	}

	fancyprint.Infof("Comparing changes for wd \"%s\" (prefix: \"%s\")\n", wd.Path, workDirPrefix)

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
			curPath = filepath.Join(workDirPrefix, curPath)
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
				if !info.IsDir() && !strings.HasPrefix(curPath, path.Join(workDirPrefix, ".dupver")) {
					fancyprint.SetColor(fancyprint.ColorCyan)
					fmt.Printf("M %s\n", curPath)
					fancyprint.ResetColor()
					// fmt.Printf("M %s\n", curPath)
					changes = true
				}
			} else if fancyprint.Verbosity >= fancyprint.InfoLevel {
				fancyprint.SetColor(fancyprint.ColorWhite)
				fmt.Printf("U %s\n", curPath)
				fancyprint.ResetColor()
			}
		} else if !strings.HasPrefix(curPath, path.Join(workDirPrefix, ".dupver")) {
			fancyprint.SetColor(fancyprint.ColorGreen)
			fmt.Printf("+ %s\n", curPath)
			fancyprint.ResetColor()
			changes = true
		}

		return nil
	}

	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)

	filepath.Walk(wd.Path, CompareAgainstSnapshot)

	for file, deleted := range deletedFiles {
		if strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}

		if deleted {
			fancyprint.SetColor(fancyprint.ColorRed)
			fmt.Printf("- %s\n", file)
			fancyprint.ResetColor()
			changes = true
		}
	}

	if !changes {
		fancyprint.Infof("No changes detected\n")
	}
}

// Compare the status of files in a working directory
// against a snapshot
// TODO: Create GetJSon functions/methods which are passed to PrintJson?
func (wd WorkDir) PrintStatusAsJson(snapshot Commit) {
	type FileStatusPrint struct {
		Status string
		Path   string
	}

	fileStatus := []FileStatusPrint{}

	workDirPrefix := ""

	if len(wd.Path) == 0 {
		wd.Path = "."
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		workDirPrefix = filepath.Base(cwd)
	}

	fancyprint.Infof("Comparing changes for wd \"%s\" (prefix: \"%s\")\n", wd.Path, workDirPrefix)

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
			curPath = filepath.Join(workDirPrefix, curPath)
		}

		curPath = strings.ReplaceAll(curPath, "\\", "/")

		if info.IsDir() {
			curPath += "/"
		}

		if snapshotInfo, ok := myFileInfo[curPath]; ok {
			deletedFiles[curPath] = false

			if snapshotInfo.ModTime != info.ModTime().Format("2006/01/02 15:04:05") {
				if !info.IsDir() && !strings.HasPrefix(curPath, path.Join(workDirPrefix, ".dupver")) {
					changes = true
					fileStatus = append(fileStatus, FileStatusPrint{Status: "Modified", Path: curPath})
				}
			} else if fancyprint.Verbosity >= fancyprint.InfoLevel {
				fileStatus = append(fileStatus, FileStatusPrint{Status: "Unchanged", Path: curPath})

			}
		} else if !strings.HasPrefix(curPath, path.Join(workDirPrefix, ".dupver")) {
			fileStatus = append(fileStatus, FileStatusPrint{Status: "Added", Path: curPath})
			changes = true
		}

		return nil
	}

	filepath.Walk(wd.Path, CompareAgainstSnapshot)

	for file, deleted := range deletedFiles {
		if strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}

		if deleted {
			fileStatus = append(fileStatus, FileStatusPrint{Status: "Deleted", Path: file})
			changes = true
		}
	}

	if !changes {
		fancyprint.Infof("No changes detected\n")
	}

	PrintJson(fileStatus)
}

// Given a partial snapshot ID, return the full snapshot ID
// by looking through the snapshots for a project
// TODO: return an error if no match
func (wd WorkDir) GetFullSnapshotId(snapshotId string) string {
	snapshotPaths := wd.ListSnapshotFiles()

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

// Return a list of the snapshot files for a given repository and project
func (wd WorkDir) ListSnapshotFiles() []string {
	snapshotsFolder := filepath.Join(wd.Repo.Path, "snapshots", wd.ProjectName)
	snapshotGlob := filepath.Join(snapshotsFolder, "*.json")
	// fmt.Println(snapshotGlob)
	fancyprint.Debugf("Snapshot glob: %s\n", snapshotGlob)
	snapshotPaths, err := filepath.Glob(snapshotGlob)

	if err != nil {
		panic(fmt.Sprintf("Error listing snapshots glob %s", snapshotGlob))
	}
	return snapshotPaths
}

// Read a snapshot given a full snapshot ID
func (wd WorkDir) ReadSnapshot(snapshot string) Commit {
	snapshotsFolder := filepath.Join(wd.Repo.Path, "snapshots", wd.ProjectName)
	snapshotPath := filepath.Join(snapshotsFolder, snapshot+".json")
	fancyprint.Debugf("Snapshot path: %s\n", snapshotPath)
	return ReadSnapshotFile(snapshotPath)
}

func (wd WorkDir) PrintSnapshot(snapshotId string) {
		snap := wd.ReadSnapshot(snapshotId)
		snap.Print()
		snap.PrintFiles()
}

func (wd WorkDir) PrintSnapshotFilesAsJson(snapshotId string) {
	snap := wd.ReadSnapshot(snapshotId)
	snap.PrintFilesAsJson()
}

// Print snapshots sorted in ascending order by date
// TODO: change the name to PrintSnapshotsByDate?
func (wd WorkDir) PrintSnapshots() {
	fancyprint.Notice("[Snapshot History]")
	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	// TODO: sort the snapshots by date
	for _, snapshotPath := range wd.ListSnapshotFiles() {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		snap := ReadSnapshotFile(snapshotPath)
		snapshotsByDate[snap.Time] = snap
		snapshotDates = append(snapshotDates, snap.Time)
	}

	sort.Strings(snapshotDates)

	for _, sdate := range snapshotDates {
		snap := snapshotsByDate[sdate]
		b := wd.Branch

		if len(b) == 0 || len(b) > 0 && b == snap.Branch {
			snap.Print()
		}
	}
}

// Print snapshots as JSON in sorted in ascending order by date
// TODO: change the name to PrintSnapshotsByDate?
// TODO: add a sort option to ListSnapshots?
// TODO: replace ListSnapshots + ReadSnapshot with ReadSnapshots
func (wd WorkDir) PrintSnapshotsAsJson() {
	type CommitPrint struct {
		ID      string
		Branch  string
		Message string
		Time    string
	}

	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	// TODO: sort the snapshots by date
	for _, snapshotPath := range wd.ListSnapshotFiles() {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		snap := ReadSnapshotFile(snapshotPath)
		snapshotsByDate[snap.Time] = snap
		snapshotDates = append(snapshotDates, snap.Time)
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
		printSnaps = append(printSnaps, ps)
	}

	PrintJson(printSnaps)
}


// Print snapshots as JSON in sorted in ascending order by date
// TODO: change the name to PrintSnapshotsByDate?
// TODO: add a sort option to ListSnapshots?
func (wd WorkDir) PrintSnapshotsAndFilesAsJson() {
	type CommitPrint struct {
		ID      string
		Branch  string
		Message string
		Time    string
		Files   []fileInfo
	}

	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	// TODO: sort the snapshots by date
	for _, snapshotPath := range wd.ListSnapshotFiles() {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		snap := ReadSnapshotFile(snapshotPath)
		snapshotsByDate[snap.Time] = snap
		snapshotDates = append(snapshotDates, snap.Time)
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
		ps.Files = snap.Files
		printSnaps = append(printSnaps, ps)
	}

	PrintJson(printSnaps)
}

// Return the most recent snapshot structure for the current project
func (wd WorkDir) LastSnapshot() (Commit, error) {
	snapshotGlob := filepath.Join(wd.Repo.Path, "snapshots", wd.ProjectName, "*.json")
	snapshotPaths, _ := filepath.Glob(snapshotGlob)

	snapshotsByDate := make(map[string]Commit)
	snapshotDates := []string{}

	br := wd.Branch

	// TODO: sort the snapshots by date
	for _, snapshotPath := range snapshotPaths {
		fancyprint.Debugf("Snapshot path: %s\n\n", snapshotPath)
		snap := ReadSnapshotFile(snapshotPath)

		if len(br) == 0 || len(br) > 0 && br == snap.Branch {
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

func (wd WorkDir) Commit(message string, jsonOutput bool) Commit {
	containingFolder := filepath.Dir(wd.Path)
	workdirFolder := filepath.Base(wd.Path)
	fancyprint.Debugf("%s -> %s, %s\n", wd.Path, containingFolder, workdirFolder)

	if len(message) == 0 {
		message = workdirFolder
		fancyprint.Infof("Message not specified, setting to: %s\n", message)
	}

	tarFile := CreateTar(containingFolder, workdirFolder)
	parentIds := []string{}
	commit := wd.CommitFile(tarFile, parentIds, message, jsonOutput)
	os.Remove(tarFile) // Delete the temporary file 
	return commit
}

// Commit a tar file into the repository. Project working directory name,
// branch and repository path are specified in the .dupver/config.toml
// file within the tar file
func (wd WorkDir) CommitFile(filePath string, parentIds []string, msg string, JsonOutput bool) Commit {
	t := time.Now()

	var snap Commit
	// var myHead Head
	snap.ID = RandHexString(SNAPSHOT_ID_LEN)
	snap.Time = t.Format("2006/01/02 15:04:05")
	snap = UpdateMessage(snap, msg, filePath)
	// TODO: write a version of ReadTarFileIndex that won't look for workdir config
	snap.Files = ReadTarFileIndex(filePath)
	snap.Branch = wd.Branch

	branchFolder := filepath.Join(wd.Repo.Path, "branches", wd.ProjectName)
	branchPath := filepath.Join(branchFolder, wd.Branch+".toml")
	myBranch := ReadBranch(branchPath)

	fancyprint.Infof("Branch: %s\nParent commit: %s\n", wd.Branch, myBranch.CommitID)
	snap.ParentIDs = append([]string{myBranch.CommitID}, parentIds...)

	chunkIDs, chunkPacks := PackFile(filePath, wd.Repo.Path, wd.Repo.ChunkerPolynomial, wd.Repo.CompressionLevel)
	snap.ChunkIDs = chunkIDs

	snapshotFolder := filepath.Join(wd.Repo.Path, "snapshots", wd.ProjectName)
	snapshotBasename := fmt.Sprintf("%s", snap.ID[0:40])
	os.Mkdir(snapshotFolder, 0777)
	snapshotPath := filepath.Join(snapshotFolder, snapshotBasename+".json")
	WriteSnapshot(snapshotPath, snap)

	// Do I really need to track commit id in head??
	myBranch.CommitID = snap.ID

	WriteBranch(branchPath, myBranch)

	treeFolder := filepath.Join(wd.Repo.Path, "trees")
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

func (workDir WorkDir) UnpackSnapshot(sid string, outFile string) {
	snapshotId := workDir.GetFullSnapshotId(sid)
	snap := workDir.ReadSnapshot(snapshotId)

	if len(outFile) == 0 {
		timeStr := TimeToPath(snap.Time)
		outFile = fmt.Sprintf("%s-%s-%s.tar", workDir.ProjectName, timeStr, snap.ID[0:16])
	}

	UnpackFile(outFile, workDir.Repo.Path, snap.ChunkIDs)

	if fancyprint.Verbosity <= fancyprint.WarningLevel {
		fmt.Println(outFile)
	} else {
		fmt.Printf("Wrote to %s\n", outFile)
	}
}

// Create a pointer-style tag given tag name and snapshot ID
// Repo path is specified in options structure
func (wd WorkDir) CreateTag(tagName string, snapshotId string) {
	tagFolder := filepath.Join(wd.Repo.Path, "tags", wd.ProjectName)
	tagPath := filepath.Join(tagFolder, tagName+".toml")
	myTag := Branch{CommitID: snapshotId}

	fancyprint.Noticef("Tag commit: %s\n", snapshotId)
	WriteBranch(tagPath, myTag)
}