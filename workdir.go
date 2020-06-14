package main

import (
    // "fmt"
    // "log"
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

