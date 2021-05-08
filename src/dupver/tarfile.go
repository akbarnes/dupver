package dupver

import (
	"fmt"
	"log"
	"crypto/sha256"
	"io"
	"bufio"
	"os"
	"os/exec"
	// "path"
	"path/filepath"
	// "strings"
	"archive/tar"
	// "encoding/json"

	// "github.com/BurntSushi/toml"
	"github.com/akbarnes/dupver/src/fancyprint"
)

// Write a project working directory to a tar file in temp
// given a working directory path and the path of its parent folder
func CreateTar(parentPath string, commitPath string) string {
	tarFile := RandHexString(40) + ".tar"
	tarFolder := filepath.Join(GetHome(), "temp")
	tarPath := filepath.Join(tarFolder, tarFile)

	// InitRepo(workDir)
	fancyprint.Debugf("Tar path: %s\n", tarPath)
	fancyprint.Debugf("Creating folder %s\n", tarFolder)

	os.Mkdir(tarFolder, 0777)

	CompressTar(parentPath, commitPath, tarPath)
	return tarPath
}

// Write a project working directory to a tar file
// given a working directory path, parent folder path and tar file path
func CompressTar(parentPath string, commitPath string, tarPath string) string {
	if len(tarPath) == 0 {
		tarPath = commitPath + ".tar"
	}

	cleanCommitPath := filepath.Clean(commitPath)

	tarCmd := exec.Command("tar", "cfv", tarPath, cleanCommitPath)
	tarCmd.Dir = parentPath
	fancyprint.Debugf("Running tar cfv %s %s\n", tarPath, cleanCommitPath)
	output, err := tarCmd.CombinedOutput()

	if err != nil {
		log.Fatal(fmt.Sprintf("Tar command failed\nOutput:\n%s\nError:\n%s\n", output, err))
	} 
	
	fancyprint.Debugf("Ran tar command with output:\n%s\n", output)
	return tarPath
}

// Read the files, workdir configuration and head from a tar file
// given a filename
func ReadTarFileIndex(filePath string) []fileInfo {
	tarFile, err := os.Open(filePath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not open input tar file %s when reading index", filePath))
	}

	files := ReadTarIndex(tarFile)
	tarFile.Close()

	return files
}

// Read the files, workdir configuration and head from a tar file
// given a file object
func ReadTarIndex(tarFile *os.File) []fileInfo {
	files := []fileInfo{}

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
		files = append(files, myFileInfo)
		i++
	}

	if fancyprint.Verbosity >= fancyprint.NoticeLevel {
		fmt.Printf("%d files/directories stored\n", i)
	}

	return files
}
