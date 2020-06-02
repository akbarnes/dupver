package main

import (
    "fmt"
    "log"
    "io"
	"os"
	"path"
    "archive/tar"
	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	WorkDirName string
	RepositoryPath string
}

type repoConfig struct {
	Version int
	ChunkerPolynomial int
}

func SaveWorkDirConfig(myConfig workDirConfig) {
	f, _ := os.Create(".dupver/config.toml")
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}



func SaveRepoConfig(repoPath string, myConfig repoConfig) {
	// TODO: add a check to make sure I don't overwrite existing
	configPath := path.Join(repoPath, "config.toml")
	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}


func ReadWorkDirConfig(filePath string) (workDirConfig) {
	var myConfig workDirConfig
	tarFile, _ := os.Open(filePath)


	// Open and iterate through the files in the archive.
	tr := tar.NewReader(tarFile)
	i := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}

		i++
		fmt.Printf("File %d: %s\n", i, hdr.Name)
	}

	tarFile.Close()

	myConfig.RepositoryPath = "C:\\Users\\305232\\DupVerRepo"
	return myConfig
}


// func ReadConfig(tarFile *os.File) {
// }
