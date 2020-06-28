package dupver

import (
	"fmt"
	"os"
	"path"
	"log"
	"github.com/BurntSushi/toml"
	"github.com/restic/chunker"
)


type repoConfig struct {
	Version int
	ChunkerPolynomial chunker.Pol
}

func InitRepo(repoPath string) {
	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")
		fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
	}	
			
	// InitRepo(workDir)
	fmt.Printf("Creating folder %s\n", repoPath)
	os.Mkdir(repoPath, 0777)

	packPath := path.Join(repoPath, "packs")
	fmt.Printf("Creating folder %s\n", packPath)
	os.Mkdir(packPath, 0777)

	snapshotsPath := path.Join(repoPath, "snapshots")
	fmt.Printf("Creating folder %s\n", snapshotsPath)
	os.MkdirAll(snapshotsPath, 0777)

	treesPath := path.Join(repoPath, "trees")
	fmt.Printf("Creating folder %s\n", treesPath)
	os.Mkdir(treesPath, 0777)	

	p, err := chunker.RandomPolynomial()
	
	if err != nil {
		panic("Error creating random polynomical while initializing repo")
	}

	var myConfig repoConfig
	myConfig.Version = 2
	myConfig.ChunkerPolynomial = p
	SaveRepoConfig(repoPath, myConfig)
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
	fmt.Printf("Creating config %s\n", configPath)

	f, err := os.Create(configPath)
	
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating repo folder %s", repoPath))
	}

	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}


func ReadRepoConfigFile(filePath string) repoConfig {
	var myConfig repoConfig

	f, err := os.Open(filePath)
	
	if err != nil {
		panic(fmt.Sprintf("Error opening repo config %s", filePath))
	}

	_, err = toml.DecodeReader(f, &myConfig)
	
	if err != nil {
		panic(fmt.Sprintf("Error decoding repo config %s", filePath))
	}

	f.Close()
	return myConfig
}

