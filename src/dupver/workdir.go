package dupver

import (
	"fmt"
	"log"

	// "bufio"
	// "io"
	"os"
	"path"
	"strings"

	// "archive/tar"
	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	WorkDirName string
	RepoPath    string
}

func FolderToWorkDirName(folder string) string {
	return strings.ReplaceAll(strings.ToLower(folder), " ", "-")
}

func InitWorkDir(workDirFolder string, workDirName string, repoPath string) {
	var configPath string
	fmt.Printf("Workdir %s, name %s, repo %s\n", workDirFolder, workDirName, repoPath)

	if len(workDirFolder) == 0 {
		fmt.Printf("Creating folder %s\n", ".dupver")
		os.Mkdir(".dupver", 0777)
		configPath = path.Join(".dupver", "config.toml")
	} else {
		fmt.Printf("Creating folder %s\n", path.Join(workDirFolder, ".dupver"))
		os.MkdirAll(path.Join(workDirFolder, ".dupver"), 0777)
		configPath = path.Join(workDirFolder, ".dupver", "config.toml")
	}

	if len(workDirName) == 0 || workDirName == "." {
		if len(workDirFolder) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			// _, folder := path.Split(dir)
			folder := path.Base(dir)
			workDirName = FolderToWorkDirName(folder)
		} else {
			workDirName = FolderToWorkDirName(workDirFolder)
		}

		fmt.Printf("Workdir name not specified, setting to %s\n", workDirName)
	}

	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")
		fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
	}

	var myConfig workDirConfig
	myConfig.RepoPath = repoPath
	myConfig.WorkDirName = workDirName
	SaveWorkDirConfig(configPath, myConfig)
}

func UpdateWorkDirName(myWorkDirConfig workDirConfig, workDirName string) workDirConfig {
	if len(workDirName) > 0 {
		myWorkDirConfig.WorkDirName = workDirName
	}

	return myWorkDirConfig
}

func SaveWorkDirConfig(configPath string, myConfig workDirConfig) {
	if _, err := os.Stat(configPath); err == nil {
		log.Fatal("Refusing to write existing project workdir config " + configPath)
	}

	f, _ := os.Create(configPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myConfig)
	f.Close()
}

func ReadWorkDirConfig(workDir string) workDirConfig {
	var configPath string

	if len(workDir) == 0 {
		configPath = path.Join(".dupver", "config.toml")
	} else {
		configPath = path.Join(workDir, ".dupver", "config.toml")
	}

	return ReadWorkDirConfigFile(configPath)
}

func ReadWorkDirConfigFile(filePath string) workDirConfig {
	var myConfig workDirConfig

	f, err := os.Open(filePath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not open project working directory config file %s", filePath))
	}

	_, err = toml.DecodeReader(f, &myConfig)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not decode TOML in project working directory config file %s", filePath))
	}

	f.Close()

	return myConfig
}