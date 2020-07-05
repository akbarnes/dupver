package dupver

import (
	"fmt"
	"log"
	"path/filepath"
	"crypto/sha256"
	"bufio"
	"io"
	"os"
	"path"
	"strings"
	"archive/tar"
	"encoding/json"

	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	WorkDirName string
	RepoPath    string
}

func FolderToWorkDirName(folder string) string {
	return strings.ReplaceAll(strings.ToLower(folder), " ", "-")
}

func InitWorkDir(workDirFolder string, workDirName string, repoPath string, verbosity int) {
	var configPath string

	if verbosity >= 2 {
		fmt.Printf("Workdir %s, name %s, repo %s\n", workDirFolder, workDirName, repoPath)
	}

	if len(workDirFolder) == 0 {
		CreateFolder(".dupver", verbosity)
		configPath = path.Join(".dupver", "config.toml")
	} else {
		CreateSubFolder(workDirFolder, ".dupver", verbosity)
		configPath = path.Join(workDirFolder, ".dupver", "config.toml")
	}

	if len(workDirName) == 0 || workDirName == "." {
		if len(workDirFolder) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			// _, folder := path.Split(dir)
			folder := filepath.Base(dir)
			fmt.Printf("%s -> %s\n", dir, folder)
			workDirName = FolderToWorkDirName(folder)
		} else {
			workDirName = FolderToWorkDirName(workDirFolder)
		}

		if workDirName == "." || workDirName == fmt.Sprintf("%c", filepath.Separator) {
			log.Fatal("Invalid project name: " + workDirName)
		}

		if verbosity >= 1 {
			fmt.Printf("Workdir name not specified, setting to %s\n", workDirName)
		}
	}

	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")

		if verbosity >= 1 {
			fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
		}
	}

	if verbosity == 0 {
		fmt.Println(workDirName)
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

func WorkDirStatus(workDir string, snapshot Commit, verbosity int) {
	workDirPrefix := ""

	if len(workDir) == 0 {
		workDir = "."
		cwd, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		workDirPrefix = filepath.Base(cwd)
	}

	if verbosity >= 2 {
		fmt.Printf("Comparing changes for wd \"%s\" (prefix: \"%s\"\n", workDir, workDirPrefix)
	}

	myFileInfo := make(map[string]fileInfo)
	deletedFiles := make(map[string]bool)
	changes := false

	for _, fi := range snapshot.Files {
		myFileInfo[fi.Path] = fi
		deletedFiles[fi.Path] = true
	}

	var CompareAgainstSnapshot = func(curPath string, info os.FileInfo, err error) error {
		// fmt.Printf("Comparing path %s\n", path)
		if len(workDirPrefix) > 0 {
			curPath = path.Join(workDirPrefix, curPath)
		}

		curPath = strings.ReplaceAll(curPath, "\\", "/")

		if info.IsDir() {
			curPath += "/"
		}

		if snapshotInfo, ok := myFileInfo[curPath]; ok {
			deletedFiles[curPath] = false

			// fmt.Printf(" mtime: %s\n", snapshotInfo.ModTime)
			// t, err := time.Parse(snapshotInfo.ModTime, "2006/01/02 15:04:05")
			// check(err)

			if snapshotInfo.ModTime != info.ModTime().Format("2006/01/02 15:04:05") {
				if !info.IsDir() {
					fmt.Printf("%sM %s%s\n", colorCyan, curPath, colorReset)
					// fmt.Printf("M %s\n", curPath)
					changes = true
				}
			} else if verbosity >= 2 {
				fmt.Printf("%sU %s%s\n", colorWhite, curPath, colorReset)
			}
		} else {
			fmt.Printf("%s+ %s%s\n", colorGreen, curPath, colorReset)
			changes = true
		}

		return nil
	}

	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)

	filepath.Walk(workDir, CompareAgainstSnapshot)

	for file, deleted := range deletedFiles {
		if strings.HasPrefix(filepath.Base(file), "._") {
			continue
		}

		if deleted {
			fmt.Printf("%s- %s%s\n", colorRed, file, colorReset)
			changes = true
		}
	}

	if !changes && verbosity >= 1 {
		fmt.Printf("No changes detected\n")
	}
}


func ReadTarFileIndex(filePath string, verbosity int) ([]fileInfo, workDirConfig, Head) {
	tarFile, err := os.Open(filePath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error: Could not open input tar file %s when reading index", filePath))
	}

	files, myConfig, myHead := ReadTarIndex(tarFile, verbosity)
	tarFile.Close()

	return files, myConfig, myHead
}

func ReadTarIndex(tarFile *os.File, verbosity int) ([]fileInfo, workDirConfig, Head) {
	files := []fileInfo{}
	var myConfig workDirConfig
	var myHead Head
	// var baseFolder string
	// var configPath string
	maxFiles := 10

	if verbosity >= 1 {
		fmt.Println("Files:")
	}

	// Open and iterate through the files in the archive.
	tr := tar.NewReader(tarFile)
	i := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(fmt.Sprintf("Error processing section while reading tar file index"))
		}

		// if i == 0 {
		// 	baseFolder = hdr.Name
		// 	myConfig.WorkDirName = baseFolder
		// 	configPath = path.Join(baseFolder,".dupver","config.toml")
		// 	// fmt.Printf("Base folder: %s\nConfig path: %s\n", baseFolder, configPath)
		// }

		if strings.HasSuffix(hdr.Name, ".dupver/config.toml") {
			if verbosity >= 1 {
				fmt.Printf("Reading config file %s\n", hdr.Name)
			}

			if _, err := toml.DecodeReader(tr, &myConfig); err != nil {
				panic(fmt.Sprintf("Error decoding repo configuration file %s while reading tar file index", hdr.Name))
			}

			// fmt.Printf("Read config\nworkdir name: %s\nrepo path: %s\n", myConfig.WorkDirName, myConfig.RepoPath)
		}

		if strings.HasSuffix(hdr.Name, ".dupver/head.json") {
			if verbosity >= 1 {
				fmt.Printf("Reading head file %s\n", hdr.Name)
			}


			myDecoder := json.NewDecoder(tr)

			if err := myDecoder.Decode(&myHead); err != nil {
				panic(fmt.Sprintf("Error decoding head file %s while reading tar file index", hdr.Name))
			}

			// fmt.Printf("Read config\nworkdir name: %s\nrepo path: %s\n", myConfig.WorkDirName, myConfig.RepoPath)
		}		

		var myFileInfo fileInfo

		bytes := make([]byte, hdr.Size)

		bufr := bufio.NewReader(tr)
		_, err = bufr.Read(bytes)

		// Name              |   256B | unlimited | unlimited
		// Linkname          |   100B | unlimited | unlimited
		// Size              | uint33 | unlimited |    uint89
		// Mode              | uint21 |    uint21 |    uint57
		// Uid/Gid           | uint21 | unlimited |    uint57
		// Uname/Gname       |    32B | unlimited |       32B
		// ModTime           | uint33 | unlimited |     int89
		// AccessTime        |    n/a | unlimited |     int89
		// ChangeTime        |    n/a | unlimited |     int89
		// Devmajor/Devminor | uint21 |    uint21 |    uint57

		myFileInfo.Path = hdr.Name
		myFileInfo.Size = hdr.Size
		myFileInfo.Hash = fmt.Sprintf("%02x", sha256.Sum256(bytes))
		myFileInfo.ModTime = hdr.ModTime.Format("2006/01/02 15:04:05")

		i++

		if i <= maxFiles && verbosity >= 1 {
			fmt.Printf("%2d: %s\n", i, hdr.Name)
		}

		files = append(files, myFileInfo)
	}

	if i > maxFiles && maxFiles > 0 && verbosity >= 1 {
		fmt.Printf("...\nSkipping %d more files\n", i-maxFiles)
	}

	return files, myConfig, myHead
}
