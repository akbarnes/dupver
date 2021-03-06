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

	// "io"
	// "bufio"
	// "crypto/sha256"
	// "encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/restic/chunker"

	"github.com/akbarnes/dupver/src/fancyprint"
)

type workDirConfig struct {
	WorkDirName string
	Branch      string
	DefaultRepo string
	Repos       map[string]string
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

	var myConfig workDirConfig
	// need to pass this as a parameter
	myConfig.DefaultRepo = repoName

	// TODO: specify an arbitrary branch
	myRepos := make(map[string]string)
	myRepos[repoName] = repoPath
	myConfig.Repos = myRepos
	myConfig.Branch = branch
	myConfig.WorkDirName = workDirName
	SaveWorkDirConfigFile(configPath, myConfig, false, opts)
}

// Add a new repository to the working directory configuration
func PrintCurrentWorkDirConfig(workDirPath string, opts Options) {
	cfg, err := ReadWorkDirConfig(workDirPath)

	if err != nil {
		// Todo: handle invalid configuration file
		fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
		os.Exit(0)
	}

	PrintWorkDirConfig(cfg, opts)
}

func PrintWorkDirReposConfig(cfg workDirConfig, opts Options) {
	for name, path := range cfg.Repos {
		repoCfg := ReadRepoConfig(path)
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

// Print the project working directory configuration
func PrintWorkDirConfig(cfg workDirConfig, opts Options) {
	// WorkDirName = "admin"
	// Branch = "test"
	// DefaultRepo = "store"

	// [Repos]
	//   main = "C:\\Users\\305232/.dupver_repo"

	fmt.Printf("Working directory name: %s\n", cfg.WorkDirName)
	fmt.Printf("Current branch: %s\n\n", cfg.Branch)
	PrintWorkDirReposConfig(cfg, opts)
}

// Add a new repository to the working directory configuration
// Todo: break up repos into  list of name, path key/value pairs
func PrintCurrentWorkDirConfigAsJson(workDirPath string, opts Options) {
	cfg, err := ReadWorkDirConfig(workDirPath)

	if err != nil {
		// Todo: handle invalid configuration file
		fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
		os.Exit(0)
	}

	PrintJson(cfg)
}

// Add a new repository to the working directory configuration
func AddRepoToWorkDir(workDirPath string, repoName string, repoPath string, makeDefaultRepo bool, opts Options) {
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

	if opts.JsonOutput {
		PrintJson(cfg)
	}

	SaveWorkDirConfig(workDirPath, cfg, true, opts)
}

// List the repositories in the working directory configuration
func ListWorkDirRepos(workDirPath string, opts Options) {
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
func ListWorkDirReposAsJson(workDirPath string, opts Options) {
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

		repoCfg := ReadRepoConfig(path)
		rl.ChunkerPolynomial = repoCfg.ChunkerPolynomial
		rl.CompressionLevel = repoCfg.CompressionLevel
		repoListings = append(repoListings, rl)
	}

	PrintJson(repoListings)
}

// Change the project name in the working directory configuration
func UpdateWorkDirName(myWorkDirConfig workDirConfig, workDirName string) workDirConfig {
	if len(workDirName) > 0 {
		myWorkDirConfig.WorkDirName = workDirName
	}

	return myWorkDirConfig
}

// Load a project working directory configuration given
// the working directory path
func ReadWorkDirConfig(workDir string) (workDirConfig, error) {
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
func ReadWorkDirConfigFile(filePath string) (workDirConfig, error) {
	var myConfig workDirConfig

	f, err := os.Open(filePath)

	if err != nil {
		return workDirConfig{}, errors.New("config file missing")
	}

	if _, err = toml.DecodeReader(f, &myConfig); err != nil {
		panic(fmt.Sprintf("Invalid configuration file: %s\n", filePath))
	}

	f.Close()

	return myConfig, nil
}

// Save a project working directory configuration given
// the working directory path
func SaveWorkDirConfig(workDir string, myConfig workDirConfig, forceWrite bool, opts Options) {
	var configPath string

	if len(workDir) == 0 {
		configPath = filepath.Join(".dupver", "config.toml")
	} else {
		configPath = filepath.Join(workDir, ".dupver", "config.toml")
	}

	SaveWorkDirConfigFile(configPath, myConfig, forceWrite, opts)
}

// Save a project working directory configuration given
// the project working directory configuration file path
func SaveWorkDirConfigFile(configPath string, myConfig workDirConfig, forceWrite bool, opts Options) {
	if _, err := os.Stat(configPath); err == nil && !forceWrite {
		// panic("Refusing to write existing project workdir config " + configPath)
		panic(fmt.Sprintf("Refusing to write existing project workdir config: %s\n", configPath))
	}

	fancyprint.Infof("Writing config:\n%+v\n", myConfig)
	fancyprint.Infof("to: %s\n", configPath)

	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}

// Compare the status of files in a working directory
// against a snapshot
func PrintWorkDirStatus(workDir string, snapshot Commit, opts Options) {
	workDirPrefix := ""

	if len(workDir) == 0 {
		workDir = "."
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		workDirPrefix = filepath.Base(cwd)
	}

	fancyprint.Infof("Comparing changes for wd \"%s\" (prefix: \"%s\")\n", workDir, workDirPrefix)

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

	filepath.Walk(workDir, CompareAgainstSnapshot)

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
func PrintWorkDirStatusAsJson(workDir string, snapshot Commit, opts Options) {
	type FileStatusPrint struct {
		Status string
		Path   string
	}

	fileStatus := []FileStatusPrint{}

	workDirPrefix := ""

	if len(workDir) == 0 {
		workDir = "."
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		workDirPrefix = filepath.Base(cwd)
	}

	fancyprint.Infof("Comparing changes for wd \"%s\" (prefix: \"%s\")\n", workDir, workDirPrefix)

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

	filepath.Walk(workDir, CompareAgainstSnapshot)

	for file, deleted := range deletedFiles {
		if strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}

		if deleted {
			fileStatus = append(fileStatus, FileStatusPrint{Status: "Deleted", Path: file})
			changes = true
		}
	}

	PrintJson(fileStatus)

	if !changes {
		fancyprint.Infof("No changes detected\n")
	}
}
