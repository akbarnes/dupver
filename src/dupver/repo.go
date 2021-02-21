package dupver

import (
	"fmt"
	"log"
	"os"
	"path"

	// "path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/restic/chunker"
)

type repoConfig struct {
	Version           int
	ChunkerPolynomial chunker.Pol
}

// Initialize a repository
func InitRepo(repoPath string, repoName string, chunkerPolynomial string, opts Options) {
	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")
		fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
	}

	CreateFolder(repoPath, opts.Verbosity)
	CreateSubFolder(repoPath, "tags", opts.Verbosity)
	CreateSubFolder(repoPath, "branches", opts.Verbosity)
	CreateSubFolder(repoPath, "snapshots", opts.Verbosity)
	CreateSubFolder(repoPath, "trees", opts.Verbosity)
	CreateSubFolder(repoPath, "packs", opts.Verbosity)

	snapshotsPath := path.Join(repoPath, "snapshots")

	if opts.Verbosity >= 1 {
		fmt.Printf("Creating folder %s\n", snapshotsPath)
	}

	os.MkdirAll(snapshotsPath, 0777)

	treesPath := path.Join(repoPath, "trees")

	if opts.Verbosity >= 1 {
		fmt.Printf("Creating folder %s\n", treesPath)
	}

	os.Mkdir(treesPath, 0777)

	if opts.Verbosity >= 1 {
		fmt.Printf("Chunker Polynomial: %s\n", chunkerPolynomial)
	}

	var poly chunker.Pol

	if len(chunkerPolynomial) == 0 {
		p, err := chunker.RandomPolynomial()

		if err != nil {
			panic("Error creating random polynomical while initializing repo")
		}

		poly = p
	} else {
		poly.UnmarshalJSON([]byte(chunkerPolynomial))
	}

	if opts.Verbosity >= 1 {
		fmt.Printf("Chunker polynomial: %d\n", poly)
	} else {
		fmt.Println(poly)
	}

	var myConfig repoConfig
	myConfig.Version = 2
	myConfig.ChunkerPolynomial = poly
	SaveRepoConfig(repoPath, myConfig, opts.Verbosity)
}

// Save a repository configuration to file
// TODO: Should I add SaveRepoCondfigFile?
func SaveRepoConfig(repoPath string, myConfig repoConfig, verbosity int) {
	// TODO: add a check to make sure I don't over`write` existing
	configPath := path.Join(repoPath, "config.toml")

	if _, err := os.Stat(configPath); err == nil {
		log.Fatal("Refusing to write existing repo config " + configPath)
	}

	if verbosity >= 1 {
		fmt.Printf("Creating config %s\n", configPath)
	}

	f, err := os.Create(configPath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating repo folder %s", repoPath))
	}

	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}

// Read a repository configuration given a file path
// TODO: Should I add ReadRepoConfig?
func ReadRepoConfigFile(filePath string) repoConfig {
	var myConfig repoConfig

	f, err := os.Open(filePath)

	if err != nil {
		panic(fmt.Sprintf("Error opening repo config %s", filePath))
	}

	if _, err := toml.DecodeReader(f, &myConfig); err != nil {
		panic(fmt.Sprintf("Error decoding repo config %s", filePath))
	}

	f.Close()
	return myConfig
}
