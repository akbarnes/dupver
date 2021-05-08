package dupver

import (
	"fmt"
	// "log"
	"archive/zip"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/restic/chunker"

	"github.com/akbarnes/dupver/src/fancyprint"
)

type RepoConfig struct {
	Version           int
	ChunkerPolynomial chunker.Pol
	CompressionLevel  uint16
}

type Repo struct {
	ChunkerPolynomial 	chunker.Pol
	CompressionLevel 	uint16
	Path 				string
}

// Initialize a repository
func InitRepo(repoPath string, repoName string, chunkerPolynomial string, compressionLevel uint16, jsonOutput bool) {
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

	var myConfig RepoConfig
	myConfig.Version = 2
	myConfig.ChunkerPolynomial = poly
	// TODO: allow compression level to be specified when creating the repo
	myConfig.CompressionLevel = zip.Deflate

	// TODO: Should this print to stderr?
	if jsonOutput {
		PrintJson(myConfig)
	} else if fancyprint.Verbosity >= fancyprint.NoticeLevel {
		fmt.Println("Chunker polynomial: %+v\n", poly)
	} else {
		fmt.Println(poly)
	}

	myConfig.Save(repoPath, false)
}

// Save a project working directory configuration given
// the working directory path
func (cfg RepoConfig) Save(repoPath string, forceOverWrite bool) {
	// TODO: add a check to make sure I don't over`write` existing
	configPath := path.Join(repoPath, "config.toml")

	// TODO: Return an error here instead of exiting
	if _, err := os.Stat(configPath); err == nil && !forceOverWrite {
		fancyprint.Warnf("Refusing to write existing repo config " + configPath)
		os.Exit(0)
	}

	fancyprint.Noticef("Creating config %s\n", configPath)
	f, err := os.Create(configPath)

	if err != nil {
		panic(fmt.Sprintf("Error creating repo folder %s", repoPath))
	}

	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(cfg)
	f.Close()
}

func ReadRepoConfig(repoPath string) (RepoConfig, error) {
	configPath := filepath.Join(repoPath, "config.toml")
	return ReadRepoConfigFile(configPath)
}

func LoadRepo(repoPath string) Repo {
	cfg, err := ReadRepoConfig(repoPath)
	Check(err)
	repo := Repo{Path: repoPath}
	repo.ChunkerPolynomial = cfg.ChunkerPolynomial
	repo.CompressionLevel = cfg.CompressionLevel
	return repo
}

// Read a repository configuration given a file path
// TODO: Should I add ReadRepoConfig?
func ReadRepoConfigFile(filePath string) (RepoConfig, error) {
	var myConfig RepoConfig
	myConfig.CompressionLevel = zip.Deflate

	f, err := os.Open(filePath)

	if err != nil {
		return RepoConfig{}, err
	}

	if _, err := toml.DecodeReader(f, &myConfig); err != nil {
		return RepoConfig{}, err
	}

	f.Close()
	return myConfig, nil
}

func (cfg RepoConfig) Print() {
	fmt.Printf("Version: %d\n", cfg.Version)
	fmt.Printf("Chunker polynomial: %d\n", cfg.ChunkerPolynomial)

	compressionDescription := "Deflate"

	if cfg.CompressionLevel == 0 {
		compressionDescription = "Store"
	}

	fmt.Printf("Compression level: %d (%s)\n", cfg.CompressionLevel, compressionDescription)	
}

func (cfg RepoConfig) PrintJson() {
	PrintJson(cfg)
}