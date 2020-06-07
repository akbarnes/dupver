package main

import (
    // "fmt"
    "log"
    // "io"
	"os"
	"path"
    // "archive/tar"
	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	WorkDirName string
	RepoPath string
}

type repoConfig struct {
	Version int
	ChunkerPolynomial int
}

func UpdateWorkDirName(myWorkDirConfig workDirConfig, workDirName string)  workDirConfig{
	if len(workDirName) > 0 {
		myWorkDirConfig.WorkDirName = workDirName
	} 

	return myWorkDirConfig
}

func UpdateRepoPath(myWorkDirConfig workDirConfig, repoPath string,) workDirConfig {
	if len(repoPath) > 0 { 
		myWorkDirConfig.RepoPath = repoPath
	}

	return myWorkDirConfig
}




func SaveRepoConfig(repoPath string, myConfig repoConfig) {
	// TODO: add a check to make sure I don't over`write` existing
	configPath := path.Join(repoPath, "config.toml")
	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
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

	if _, err := toml.DecodeReader(f, &myConfig); err != nil {
		log.Fatal(err)
	}

	f.Close()

	return myConfig
}

