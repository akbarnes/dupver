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



func ReadWorkDirConfig(filePath string) (workDirConfig) {
	var myConfig workDirConfig


	f, _ := os.Open(filePath)

	if _, err := toml.DecodeReader(f, &myConfig); err != nil {
		log.Fatal(err)
	}

	f.Close()

	return myConfig
}

