package main

import (
	"fmt"
	"os"
	"path"
)

func InitRepo(repoPath string) {
	fmt.Printf("Creating folder %s\n", repoPath)
	os.Mkdir(repoPath, 0777)

	packPath := path.Join(repoPath, "packs")
	fmt.Printf("Creating folder %s\n", packPath)
	os.Mkdir(packPath, 0777)

	treesPath := path.Join(repoPath, "trees")
	fmt.Printf("Creating folder %s\n", treesPath)
	os.Mkdir(treesPath, 0777)		

	snapshotsPath := path.Join(repoPath, "snapshots")
	fmt.Printf("Creating folder %s\n", snapshotsPath)
	os.Mkdir(snapshotsPath, 0777)

	var myConfig repoConfig
	myConfig.Version = 1
	myConfig.ChunkerPolynomial = 0x3DA3358B4DC173
	SaveRepoConfig(repoPath, myConfig)
}