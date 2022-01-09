package dupver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
    "time"
    "math/rand"
	"path/filepath"

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
	dupverDir := filepath.Join(".dupver")

	if err := os.MkdirAll(dupverDir, 0777); err != nil {
		panic(fmt.Sprintf("Error creating dupver folder %s\n", dupverDir))
	}

	cfgPath := filepath.Join(".dupver", "repo_config.json")
	f, err := os.Create(cfgPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create repository configuration json %s", cfgPath))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(cfg)
	f.Close()
}

func (cfg RepoConfig) CorrectRepoVersion() bool {
	return cfg.RepoMajorVersion == RepoMajorVersion
}

func AbortIfIncorrectRepoVersion() {
	cfg, err := ReadRepoConfig(false)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read repo configuration, exiting\n")
		os.Exit(1)
	}

	if !cfg.CorrectRepoVersion() {
		fmt.Fprintf(os.Stderr, "Incorrect repo version of %d.%d, expecting %d.x\n", cfg.RepoMajorVersion, cfg.RepoMinorVersion, RepoMajorVersion)
		os.Exit(1)
	}
}

func ReadRepoConfig(writeIfMissing bool) (RepoConfig, error) {
	cfgPath := filepath.Join(".dupver", "repo_config.json")

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
			return cfg, nil
		} else {
			return RepoConfig{}, err
		}
	} else if err != nil {
		return RepoConfig{}, errors.New("Cannot open repo config")
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
