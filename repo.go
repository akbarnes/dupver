package main

import (
	// "fmt"
	"os"
	"path"
	"github.com/BurntSushi/toml"
)


type repoConfig struct {
	Version int
	ChunkerPolynomial int
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


// func InitRepo(repoPath string) {
// 	fmt.Printf("Creating folder %s\n", repoPath)
// 	os.Mkdir(repoPath, 0777)

// 	packPath := path.Join(repoPath, "packs")
// 	fmt.Printf("Creating folder %s\n", packPath)
// 	os.Mkdir(packPath, 0777)

// 	snapshotsPath := path.Join(repoPath, "snapshots")
// 	fmt.Printf("Creating folder %s\n", snapshotsPath)
// 	os.MkdirAll(snapshotsPath, 0777)

// 	treesPath := path.Join(repoPath, "trees")
// 	fmt.Printf("Creating folder %s\n", treesPath)
// 	os.Mkdir(treesPath, 0777)	

// 	var myConfig repoConfig
// 	myConfig.Version = 1
// 	myConfig.ChunkerPolynomial = 0x3DA3358B4DC173
// 	SaveRepoConfig(repoPath, myConfig)
// }