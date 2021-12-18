package dupver

import (
    "fmt"
    "errors"
    "path/filepath"
    "encoding/json"
    "os"
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
	RepoMajorVersion int64
	RepoMinorVersion int64
	DupverMajorVersion int64
	DupverMinorVersion int64
	CompressionLevel uint16
    ChunkerPoly string
}

func CreateDefaultRepoConfig() RepoConfig {
    cfg := RepoConfig{}
    cfg.DupverMajorVersion = MajorVersion
    cfg.DupverMinorVersion = MinorVersion
    cfg.RepoMajorVersion = RepoMajorVersion
    cfg.RepoMinorVersion = RepoMinorVersion
    cfg.CompressionLevel = 0
    cfg.ChunkerPoly = "0x3abc9bff07d9e5"
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
    return (cfg.RepoMajorVersion == MajorVersion)
}

func ReadRepoConfig(writeIfMissing bool) (RepoConfig, error) {
	cfgPath := filepath.Join(".dupver", "repo_config.json")

	if VerboseMode {
		fmt.Printf("Reading %s\n", cfgPath)
	}

	var cfg RepoConfig
	f, err := os.Open(cfgPath)
    defer f.Close()

    if errors.Is(err, os.ErrNotExist) {
        if writeIfMissing {
            if VerboseMode {
                fmt.Println("Repo configuration not present, writing default")
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
		panic(fmt.Sprintf("Invalid repo version %d.%d, expecting %d.x\n", cfg.RepoMajorVersion, cfg.RepoMinorVersion, MajorRepoVersion))
    }

	return cfg, nil
}


