package dupver

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	// "io"
	// "bufio"
	"os"
	"strconv"
	"strings"

	// "crypto/sha256"
	// "encoding/json"

	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	WorkDirName string
	Branch string
	DefaultRepo string
	// RepoPath    string
	Repos map[string]string
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

	if opts.Verbosity >= 2 {
		fmt.Printf("Workdir %s, name %s, repo %s\n", workDirFolder, workDirName, opts.RepoPath)
	}

	if len(workDirFolder) == 0 {
		CreateFolder(".dupver", opts.Verbosity)
		configPath = filepath.Join(".dupver", "config.toml")
	} else {
		CreateSubFolder(workDirFolder, ".dupver", opts.Verbosity)
		configPath = filepath.Join(workDirFolder, ".dupver", "config.toml")
	}

	if opts.Verbosity >= 2 {
		fmt.Println("Writing workdir config file to: " + configPath)
	}

	if len(workDirName) == 0 || workDirName == "." {
		if len(workDirFolder) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			// _, folder := path.Split(dir)
			folder := filepath.Base(dir)

			if opts.Verbosity >= 3 {
				fmt.Printf("Resolving folder %s to %s\n", dir, folder)
			}
			workDirName = FolderToWorkDirName(folder)
		} else {
			workDirName = FolderToWorkDirName(workDirFolder)
		}

		if workDirName == "." || workDirName == fmt.Sprintf("%c", filepath.Separator) {
			log.Fatal("Invalid project name: " + workDirName)
		}

		if opts.Verbosity >= 1 {
			fmt.Printf("Workdir name not specified, setting to %s\n", workDirName)
		}
	}

	if len(repoPath) == 0 {
		repoPath = filepath.Join(GetHome(), ".dupver_repo")

		if opts.Verbosity >= 1 {
			fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
		}
	}

	if opts.Verbosity == 0 {
		fmt.Println(workDirName)
	} else {
		fmt.Printf("Repo name: [%s]\n", repoName)
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
func AddRepoToWorkDir(workDirPath string, repoName string, repoPath string, opts Options) {
	cfg := ReadWorkDirConfig(workDirPath)
	cfg.Repos[repoName] = repoPath
	SaveWorkDirConfig(workDirPath, cfg, true, opts)
}

// List the repositories in the working directory configuration
func ListWorkDirRepos(workDirPath string, opts Options) {
	cfg := ReadWorkDirConfig(workDirPath)
	maxLen := 0

	for name, _ := range cfg.Repos {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	fmtStr := "%" + strconv.Itoa(maxLen) + "s: %s\n"

	for name, path := range cfg.Repos {
		if opts.Verbosity == 0 {
			fmt.Printf("%s %s\n", name, path)
		} else {
			fmt.Printf(fmtStr, name, path)
		}
	}
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
func ReadWorkDirConfig(workDir string) workDirConfig {
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
func ReadWorkDirConfigFile(filePath string) workDirConfig {
	var myConfig workDirConfig

	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not open project working directory config file %s", filePath))
	}

	if _, err = toml.DecodeReader(f, &myConfig); err != nil {
		log.Fatal(fmt.Sprintf("Could not decode TOML in project working directory config file %s", filePath))
	}

	f.Close()

	return myConfig
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
		log.Fatal("Refusing to write existing project workdir config " + configPath)
	}

	if opts.Verbosity >= 2 {
		fmt.Printf("Writing config:\n%+v\n", myConfig)
		fmt.Printf("to: %s\n", configPath)
	}

	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}

// Compare the status of files in a working directory
// against a snapshot
func WorkDirStatus(workDir string, snapshot Commit, opts Options) {
	workDirPrefix := ""

	if len(workDir) == 0 {
		workDir = "."
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		workDirPrefix = filepath.Base(cwd)
	}

	if opts.Verbosity >= 2 {
		fmt.Printf("Comparing changes for wd \"%s\" (prefix: \"%s\")\n", workDir, workDirPrefix)
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
					if opts.Color {
						fmt.Printf("%s", colorCyan)
					}

					fmt.Printf("M %s\n", curPath)

					if opts.Color {
						fmt.Printf("%s", colorReset)
					}
					// fmt.Printf("M %s\n", curPath)
					changes = true
				}
			} else if opts.Verbosity >= 2 {
				if opts.Color {
					fmt.Printf("%s", colorWhite)
				}

				fmt.Printf("U %s\n", curPath)

				if opts.Color {
					fmt.Printf("%s", colorReset)
				}
			}
		} else if !strings.HasPrefix(curPath, path.Join(workDirPrefix, ".dupver")) {
			if opts.Color {
				fmt.Printf("%s", colorGreen)
			}

			fmt.Printf("+ %s\n", curPath)

			if opts.Color {
				fmt.Printf("%s", colorReset)
			}
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
			if opts.Color {
				fmt.Printf("%s", colorRed)
			}

			fmt.Printf("- %s\n", file)

			if opts.Color {
				fmt.Printf("%s", colorReset)
			}

			changes = true
		}
	}

	if !changes && opts.Verbosity >= 1 {
		fmt.Printf("No changes detected\n")
	}
}
