package main

import (
    "fmt"
    // "log"
	// "io"
	"strings"
	"os"
	"path"
    // "archive/tar"
	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	WorkDirName string
	RepoPath string
}

func InitWorkDir(workDirFolder string, workDirName string, repoPath string) {
	var configPath string
	// fmt.Printf("Workdir %s, name %s, repo %s\n", workDirFolder, workDirName, repoPath)

	if len(workDirFolder) == 0 {
		 os.Mkdir(".dupver", 0777)
		 configPath = path.Join(".dupver", "config.toml")
	} else {
		 os.MkdirAll(path.Join(workDirFolder, ".dupver"), 0777)
		 configPath = path.Join(workDirFolder, ".dupver", "config.toml")
	}

	if len(workDirName) == 0 {
		if len(workDirFolder) == 0 {
			panic("Both workDir and workDirName are empty")
		} else {
			workDirName = strings.ToLower(path.Base(workDirFolder))
		}
	}

	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")
		fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
	}		

	var myConfig workDirConfig
	myConfig.RepoPath = repoPath
	myConfig.WorkDirName = workDirName
	SaveWorkDirConfig(configPath, myConfig)
}

func UpdateWorkDirName(myWorkDirConfig workDirConfig, workDirName string)  workDirConfig{
	if len(workDirName) > 0 {
		myWorkDirConfig.WorkDirName = workDirName
	} 

	return myWorkDirConfig
}


func SaveWorkDirConfig(configPath string, myConfig workDirConfig) {
	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}

func ReadWorkDirConfig(workDir string) workDirConfig {
	var configPath string

	if len(workDir) == 0 {
		configPath = path.Join(".dupver", "config.toml")
	} else {
		configPath = path.Join(workDir, ".dupver", "config.toml")
	}

	return ReadWorkDirConfigFile(configPath)
}


func ReadWorkDirConfigFile(filePath string) workDirConfig {
	var myConfig workDirConfig

	f, _ := os.Open(filePath)

	_, err := toml.DecodeReader(f, &myConfig)
	check(err)

	f.Close()

	return myConfig
}

