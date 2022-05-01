package dupver

import (
    "testing"

	"github.com/restic/chunker"
)

func TestCreateDefaultRepo(test *testing.T) {
    cfg := CreateDefaultRepoConfig()

	if cfg.DupverMajorVersion != MajorVersion {
        test.Errorf("Repo config dupver major version of %d != %d", cfg.DupverMajorVersion, MajorVersion)
    }

	if cfg.DupverMinorVersion != MinorVersion {
        test.Errorf("Repo config dupver minor version of %d != %d", cfg.DupverMinorVersion, MinorVersion)
    }

	if cfg.RepoMajorVersion != RepoMajorVersion {
        test.Errorf("Repo config major version of %d != %d", cfg.RepoMajorVersion, RepoMajorVersion)
    }

	if cfg.RepoMinorVersion != RepoMinorVersion {
        test.Errorf("Repo config minor version of %d != %d", cfg.RepoMinorVersion, RepoMinorVersion)
    }

    if cfg.PackSize != PackSize {
        test.Errorf("Repo config pack size of %d != %d", cfg.PackSize, PackSize)
    }

    if cfg.CompressionLevel != 0 {
        test.Errorf("Repo config compression level of %d != %d", cfg.CompressionLevel, 0)
    }

    const p chunker.Pol = 0x3abc9bff07d9e5

    if cfg.ChunkerPoly != p {
        test.Errorf("Repo chunker poly of %d != %d", cfg.ChunkerPoly, p)
    }
}


