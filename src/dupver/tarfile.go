package dupver

import (
	"fmt"
	// "log"
	"crypto/sha256"
	"io"
	"bufio"
	"os"
	// "path"
	// "path/filepath"
	"strings"
	"archive/tar"
	// "encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/akbarnes/dupver/src/fancyprint"
)

// Read the files, workdir configuration and head from a tar file
// given a filename
func ReadTarFileIndex(filePath string) ([]fileInfo, workDirConfig) {
	tarFile, err := os.Open(filePath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not open input tar file %s when reading index", filePath))
	}

	files, myConfig := ReadTarIndex(tarFile)
	tarFile.Close()

	return files, myConfig
}

// Read the files, workdir configuration and head from a tar file
// given a file object
func ReadTarIndex(tarFile *os.File) ([]fileInfo, workDirConfig) {
	files := []fileInfo{}
	var myConfig workDirConfig
	// var baseFolder string
	// var configPath string
	maxFiles := 10

	if fancyprint.Verbosity >= fancyprint.NoticeLevel {
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
			fancyprint.Infof("Reading config file %s\n", hdr.Name)

			if _, err := toml.DecodeReader(tr, &myConfig); err != nil {
				panic(fmt.Sprintf("Error decoding repo configuration file %s while reading tar file index", hdr.Name))
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

		if i <= maxFiles && fancyprint.Verbosity >= fancyprint.NoticeLevel {
			fmt.Printf("%2d: %s\n", i, hdr.Name)
		}

		files = append(files, myFileInfo)
	}

	if i > maxFiles && maxFiles > 0 && fancyprint.Verbosity >= fancyprint.NoticeLevel {
		fmt.Printf("...\nSkipping %d more files\n", i-maxFiles)
	}

	return files, myConfig
}
