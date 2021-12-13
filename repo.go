package dupver

import "time"

const SNAPSHOT_ID_LEN int = 40
const PACK_ID_LEN int = 64

// TODO: change this to SerializedSnaphot
// and use Time type for SnapshotTime?
type RepoConfig struct {
	RepoMajorVersion int64
	RepoMinorVersion int64
	DupverMajorVersion int64
	DupverMinorVersion int64
	CompressionLevel  int16
    ChunkerPoly string
}

func CreateDefaultRepoConfig() {
    cfg := RepoConfig{}
    cfg.DupverMajorVersion = DupverMajorVersion
    cfg.DupverMinorVersion = DupverMinorVersion
    cfg.RepoMajorVersion = RepoMajorVersion
    cfg.RepoMinorVersion = RepoMinorVersion
    cfg.CompressionLevel = CompressionLevel
    cfg.ChunkerPoly = "0x3abc9bff07d9e5"
}


func (cfg RepoConfig) Write() {
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

func ReadRepoConfig() (RepoConfig, error) {
	cfgPath := filepath.Join(".dupver", "repo_config.json")

	if VerboseMode {
		fmt.Printf("Reading %s\n", cfgPath)
	}

	var cfg RepoConfig
	f, err := os.Open(cfgPath)
    defer f.Close()

    if errors.Is(err, os.ErrNotExist) {
    	if VerboseMode {
		    fmt.Printf("No repo configuration present, creating")
	    }

        cfg = CreateDefaultRepoConfig()
        cfg.Write()
        return cfg, nil
	} else if err != nil {
		return RepoConfig{}, errors.New("Cannot open repo config")
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&cfg); err != nil {
		return RepoConfig{}, errors.New("Cannot decode repo config")
	}

	return cfg, nil
}


