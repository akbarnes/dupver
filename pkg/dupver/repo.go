package dupver

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/restic/chunker"
)

type RepoConfigMissingError struct {
	err string
}

func (e *RepoConfigMissingError) Error() string {
	return "Repo config is missing"
}

// TODO: change this to SerializedSnaphot
// and use Time type for SnapshotTime?
type RepoConfig struct {
	RepoMajorVersion   int64
	RepoMinorVersion   int64
	DupverMajorVersion int64
	DupverMinorVersion int64
	CompressionLevel   uint16
	PackSize           int64
	ChunkerPoly        chunker.Pol
}

const RepoReadMe = `# Repo Structure

## Directory Organization

The repository is structured is as follows:

- dupver_settings.toml
- snapshots/
- files/
- trees/
- packs/
- head.json

## Repo Settings 

Relative path: dupver_settings.json

This contains the following fields:

- RepoMajorVersion, RepoMinorVersion: Repo version
- DupverMajorVersion, DupverMinorVersion: Versio of Dupver used to initialize the repo
- CompressionLevel: Zip compression level
- PackSize: Target packfile size in bytes. It's possible and common for the packfiles to be larger than this size by about 1 MB
- ChunkerPoly: Chunker polynomial

For recovering data it's not necessary to read the preferences file.

## Snapshots 

Relative path: snapshots/<snapshot_id>.json

This folder contains a set of JSON files, one for each snapshot. 
The base name of each file is its snapshot ID, and each file contains
the following fields:

- Message: Commit message
- Username: Username of user who created the commit
- SnapshotTime: UTC time of commit
- SnapshotLocalTime: Local time of commit according to the computer that created the commit
- SnapshotID: Hex ID of commit

## File Listings 

Relative path: files/<snapshot_id>.json

This folder contains a set of JSON files, one for each snapshot. 
The base name of each file is its snapshot ID, and each file 
contains a list, where each element has the following fields

- Name: Relative file name
- Size: Size of file in bytes
- ModTime: Last modification time of the file in UTC time
- ModLocalTime: Last modification time of the file according to the computer that created the commit
- ChunkIds: List of hex chunk IDs for the file content
- IsArchive: Indicates if the file is a archive (e.g. zip, tgz, 7z, docx) and was pre-processed to store the uncompressed archive contents as a store-only zip file

## Trees 

Relative path: trees/<snapshot_id>.json

This folder contains a set of JSON files, where each file consists of a dictionary.
The keys of the dictionaries are the hex IDs for packfiles, while the values of 
the dictionaries are lists of the chunk hex IDs stored in each packfile.

## Packs

Relative path: packs/<pack_id>.json

This folder contains a set of pack files in zip format, where the base name of each file corresponds
to its pack ID. Each pack file is an archive where the stored files are chunks, whose
filenames correspond to their hex chunk IDs.

# Head

Relative path: head.json

This contains the following fields:

- SnapshotTime: UTC time of the last snapshot
- SnapshotID: Hex ID of the last snapshot
`

func CreateDefaultRepoConfig() RepoConfig {
	cfg := RepoConfig{}
	cfg.DupverMajorVersion = MajorVersion
	cfg.DupverMinorVersion = MinorVersion
	cfg.RepoMajorVersion = RepoMajorVersion
	cfg.RepoMinorVersion = RepoMinorVersion
	cfg.PackSize = PackSize
	cfg.CompressionLevel = 0
	cfg.ChunkerPoly = 0x3abc9bff07d9e5

	if RandomPoly {
		rand.Seed(time.Now().UnixNano())

		p, err := chunker.RandomPolynomial()

		if err == nil {
			cfg.ChunkerPoly = p
		} else {
			fmt.Fprintf(os.Stderr, "Error generating random polynomial, using default of %v\n", cfg.ChunkerPoly)
		}
	}

	if VerboseMode {
		fmt.Fprintf(os.Stderr, "Generated random polynomial of %v\n", cfg.ChunkerPoly)
	}

	return cfg
}

func (cfg RepoConfig) Write() {
	dupverDir := filepath.Join(WorkingDirectory, ".dupver")

	if err := os.MkdirAll(dupverDir, 0777); err != nil {
		panic(fmt.Sprintf("Error creating dupver folder %s\n", dupverDir))
	}

	cfgPath := filepath.Join(WorkingDirectory, ".dupver", "repo_config.json")
	f, err := os.Create(cfgPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create repository configuration json %s", cfgPath))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(cfg)
	f.Close()
}

func (cfg RepoConfig) WriteReadme() {
	dupverDir := filepath.Join(WorkingDirectory, ".dupver")

	if err := os.MkdirAll(dupverDir, 0777); err != nil {
		panic(fmt.Sprintf("Error creating dupver folder %s\n", dupverDir))
	}

	readmePath := filepath.Join(WorkingDirectory, ".dupver", "README.txt")
	f, err := os.Create(readmePath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create README %s", readmePath))
	}

	fmt.Fprintf(f, "%s", RepoReadMe)
	f.Close()
}

func (cfg RepoConfig) CorrectRepoVersion() bool {
	return cfg.RepoMajorVersion == RepoMajorVersion
}

func AbortIfIncorrectRepoVersion() {
	cfg, err := ReadRepoConfig(false)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read repo configuration, exiting: %v\n", err)
		os.Exit(1)
	}

	if !cfg.CorrectRepoVersion() {
		fmt.Fprintf(os.Stderr, "Incorrect repo version of %d.%d, expecting %d.x\n", cfg.RepoMajorVersion, cfg.RepoMinorVersion, RepoMajorVersion)
		os.Exit(1)
	}
}

func ReadRepoConfig(writeIfMissing bool) (RepoConfig, error) {
	cfgPath := filepath.Join(WorkingDirectory, ".dupver", "repo_config.json")

	if VerboseMode {
		fmt.Fprintf(os.Stderr, "Reading %s\n", cfgPath)
	}

	var cfg RepoConfig
	f, err := os.Open(cfgPath)
	defer f.Close()

	if errors.Is(err, os.ErrNotExist) {
		if writeIfMissing {
			if VerboseMode {
				fmt.Fprintf(os.Stderr, "Repo configuration not present, writing default\n")
			}

			cfg = CreateDefaultRepoConfig()
			cfg.Write()
			cfg.WriteReadme()
			return cfg, nil
		} else {
			return RepoConfig{}, err
		}
	} else if err != nil {
		return RepoConfig{}, errors.New("Cannot open repo config: " + cfgPath)
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&cfg); err != nil {
		panic("Cannot decode repo config")
	}

	if !cfg.CorrectRepoVersion() {
		panic(fmt.Sprintf("Invalid repo version %d.%d, expecting %d.x\n", cfg.RepoMajorVersion, cfg.RepoMinorVersion, RepoMajorVersion))
	}

	if cfg.PackSize == 0 {
		fmt.Fprintf(os.Stderr, "Warning: Repo PackSize = 0, consider setting to 524288000\n")
	}

	if VerboseMode || DebugMode {
		fmt.Fprintf(os.Stderr, "Read random polynomial of %v\n", cfg.ChunkerPoly)
	}

	return cfg, nil
}
