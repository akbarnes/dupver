package dupver

import (
	"fmt"
	"log"
	// "path"
	"path/filepath"
	// "io"
	// "bufio"
	"os"
	"strings"
	// "crypto/sha256"
	// "encoding/json"

	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	WorkDirName string
	RepoPath    string
}

func FolderToWorkDirName(folder string) string {
	return strings.ReplaceAll(strings.ToLower(folder), " ", "-")
}

func InitWorkDir(workDirFolder string, workDirName string, repoPath string, verbosity int) {
	var configPath string

	if verbosity >= 2 {
		fmt.Printf("Workdir %s, name %s, repo %s\n", workDirFolder, workDirName, repoPath)
	}

	if len(workDirFolder) == 0 {
		CreateFolder(".dupver", verbosity)
		configPath = filepath.Join(".dupver", "config.toml")
	} else {
		CreateSubFolder(workDirFolder, ".dupver", verbosity)
		configPath = filepath.Join(workDirFolder, ".dupver", "config.toml")
	}

	if len(workDirName) == 0 || workDirName == "." {
		if len(workDirFolder) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			// _, folder := path.Split(dir)
			folder := filepath.Base(dir)
			fmt.Printf("%s -> %s\n", dir, folder)
			workDirName = FolderToWorkDirName(folder)
		} else {
			workDirName = FolderToWorkDirName(workDirFolder)
		}

		if workDirName == "." || workDirName == fmt.Sprintf("%c", filepath.Separator) {
			log.Fatal("Invalid project name: " + workDirName)
		}

		if verbosity >= 1 {
			fmt.Printf("Workdir name not specified, setting to %s\n", workDirName)
		}
	}

	if len(repoPath) == 0 {
		repoPath = filepath.Join(GetHome(), ".dupver_repo")

		if verbosity >= 1 {
			fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
		}
	}

	if verbosity == 0 {
		fmt.Println(workDirName)
	}

	var myConfig workDirConfig
	myConfig.RepoPath = repoPath
	myConfig.WorkDirName = workDirName
	SaveWorkDirConfig(configPath, myConfig)
}

func UpdateWorkDirName(myWorkDirConfig workDirConfig, workDirName string) workDirConfig {
	if len(workDirName) > 0 {
		myWorkDirConfig.WorkDirName = workDirName
	}

	return myWorkDirConfig
}

func SaveWorkDirConfig(configPath string, myConfig workDirConfig) {
	if _, err := os.Stat(configPath); err == nil {
		log.Fatal("Refusing to write existing project workdir config " + configPath)
	}

	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
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

	_, err = toml.DecodeReader(f, &myConfig)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not decode TOML in project working directory config file %s", filePath))
	}

	f.Close()

	return myConfig
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

func WriteHead(headPath string, myHead Head, verbosity int) {
	dir := filepath.Dir(headPath)
	CreateFolder(dir, verbosity)

	if verbosity >= 2 {
		fmt.Println("Writing head to " +  headPath)
	}

	f, err := os.Create(headPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create head file %s", headPath))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(myHead)
	f.Close()
}

func ReadHead(headPath string) Head {
	var myHead Head
	f, err := os.Open(headPath)

	if err != nil {
		//panic(fmt.Sprintf("Error: Could not read head file %s", headPath))
		fmt.Printf("No head file exists, returning defaut head struct\n")
		return Head{BranchName: "main"}
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&myHead); err != nil {
		panic(fmt.Sprintf("Error:could not decode head file %s", headPath))
	}

	f.Close()
	return myHead
}


