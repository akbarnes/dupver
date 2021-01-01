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

// "name": "default",
// "id": "macbook-air-home",
// "repository": "",
// "storage": "/Volumes/Shared/Backups/Duplicacy/MBAir",
// "encrypted": true,
// "no_backup": false,
// "no_restore": false,
// "no_save_password": false,
// "nobackup_file": "",
// "keys": null

type workDirConfig struct {
	WorkDirName string
	BranchName  string
	DefaultRepo string
	// RepoPath    string
	Repos map[string]string
}

func FolderToWorkDirName(folder string) string {
	return strings.ReplaceAll(strings.ToLower(folder), " ", "-")
}

func InitWorkDir(workDirFolder string, workDirName string, opts Options) {
	var configPath string
	repoName := opts.RepoName
	repoPath := opts.RepoPath

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
	myConfig.BranchName = "main"
	myRepos := make(map[string]string)
	myRepos[repoName] = repoPath
	myConfig.Repos = myRepos
	myConfig.WorkDirName = workDirName
	SaveWorkDirConfigFile(configPath, myConfig, false, opts)
}

var configPath string

func AddRepoToWorkDir(workDirPath string, repoName string, repoPath string, opts Options) {
	cfg := ReadWorkDirConfig(workDirPath)
	cfg.Repos[repoName] = repoPath
	SaveWorkDirConfig(workDirPath, cfg, true, opts)
}

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

func UpdateWorkDirName(myWorkDirConfig workDirConfig, workDirName string) workDirConfig {
	if len(workDirName) > 0 {
		myWorkDirConfig.WorkDirName = workDirName
	}

	return myWorkDirConfig
}

func ReadWorkDirConfig(workDir string) workDirConfig {
	var configPath string

	if len(workDir) == 0 {
		configPath = filepath.Join(".dupver", "config.toml")
	} else {
		configPath = filepath.Join(workDir, ".dupver", "config.toml")
	}

	return ReadWorkDirConfigFile(configPath)
}

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

func SaveWorkDirConfig(workDir string, myConfig workDirConfig, forceWrite bool, opts Options) {
	var configPath string

	if len(workDir) == 0 {
		configPath = filepath.Join(".dupver", "config.toml")
	} else {
		configPath = filepath.Join(workDir, ".dupver", "config.toml")
	}

	SaveWorkDirConfigFile(configPath, myConfig, forceWrite, opts)
}

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

func WriteHead(headPath string, myHead Head, opts Options) {
	dir := filepath.Dir(headPath)
	CreateFolder(dir, opts.Verbosity)

	if opts.Verbosity >= 2 {
		fmt.Println("Writing head to " + headPath)
	}

	f, err := os.Create(headPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create head file %s", headPath))
	}

	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myHead)
	f.Close()
}

func ReadHead(headPath string, opts Options) Head {
	var myHead Head
	f, err := os.Open(headPath)

	if err != nil {
		//panic(fmt.Sprintf("Error: Could not read head file %s", headPath))
		if opts.Verbosity >= 2 {
			fmt.Printf("No head file exists, returning default head struct\n")
		}
		return Head{BranchName: "main"}
	}
	if _, err := toml.DecodeReader(f, &myHead); err != nil {
		panic(fmt.Sprintf("Error:could not decode head file %s", headPath))
	}

	f.Close()
	return myHead
}
