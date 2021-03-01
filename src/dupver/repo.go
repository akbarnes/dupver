package dupver

import (
	"fmt"
	// "log"
	"os"
	"path"
	"archive/zip"

	// "path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/restic/chunker"

	"github.com/akbarnes/dupver/src/fancyprint"
)

type repoConfig struct {
	Version           int
	ChunkerPolynomial chunker.Pol
	CompressionLevel uint16
}

// Initialize a repository
func InitRepo(repoPath string, repoName string, chunkerPolynomial string, opts Options) {
	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")
		fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
	}

	CreateFolder(repoPath)
	CreateSubFolder(repoPath, "tags")
	CreateSubFolder(repoPath, "branches")
	CreateSubFolder(repoPath, "snapshots")
	CreateSubFolder(repoPath, "trees")
	CreateSubFolder(repoPath, "packs")

	snapshotsPath := path.Join(repoPath, "snapshots")
	fancyprint.Noticef("Creating folder %s\n", snapshotsPath)
	os.MkdirAll(snapshotsPath, 0777)

	treesPath := path.Join(repoPath, "trees")
	fancyprint.Noticef("Creating folder %s\n", treesPath)
	os.Mkdir(treesPath, 0777)

	fancyprint.Debugf("Chunker Polynomial: %s\n", chunkerPolynomial)
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

	// TODO: Should this print to stderr?

	if fancyprint.Verbosity >= fancyprint.NoticeLevel {
		fmt.Println("Chunker polynomial: %d\n", poly)
	} else {
		fmt.Println(poly)
	}

	var myConfig repoConfig
	myConfig.Version = 2
	myConfig.ChunkerPolynomial = poly
	// TODO: allow compression level to be specified when creating the repo
	myConfig.CompressionLevel = zip.Deflate
	SaveRepoConfig(repoPath, myConfig)
}

// Save a repository configuration to file
// TODO: Should I add SaveRepoCondfigFile?
func SaveRepoConfig(repoPath string, myConfig repoConfig) {
	// TODO: add a check to make sure I don't over`write` existing
	configPath := path.Join(repoPath, "config.toml")

	// TODO: Return an error here instead of exiting
	if _, err := os.Stat(configPath); err == nil {
		fancyprint.Warnf("Refusing to write existing repo config " + configPath)
		os.Exit(0)
	}

	fancyprint.Noticef("Creating config %s\n", configPath)
	f, err := os.Create(configPath)

	if err != nil {
		panic(fmt.Sprintf("Error creating repo folder %s", repoPath))
	}

	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}

// Read a repository configuration given a file path
// TODO: Should I add ReadRepoConfig?
func ReadRepoConfigFile(filePath string) repoConfig {
	var myConfig repoConfig
	myConfig.CompressionLevel = zip.Deflate

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
