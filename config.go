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
	RepositoryPath string
}

type repoConfig struct {
	Version int
	ChunkerPolynomial int
}

func UpdateWorkDirName(workDirName *string, myWorkDirConfig workDirConfig) {
	if len(*workDirName) == 0 {
		*workDirName = myWorkDirConfig.WorkDirName
	} 
}



func SaveRepoConfig(repoPath string, myConfig repoConfig) {
	// TODO: add a check to make sure I don't overwrite existing
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

