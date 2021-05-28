package dupver

import (
	"archive/tar"
	"encoding/json"
	"os/exec"

	// "bufio"
	// "errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	// "github.com/BurntSushi/toml"
	"github.com/akbarnes/dupver/src/fancyprint"
)

// DiffTool = "bcompare.exe"
// DefaultRepo = "main"

// [Repos]
//   main = "C:\\Users\\Art\\.dupver_repo"

// Halt if error parameter is not nil
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Calculate the verbosity level given parameters
func CalculateVerbosity(debug bool, verbose bool, quiet bool) int {
	if quiet {
		return 0
	}

	if debug {
		return 3
	}

	if verbose {
		return 2
	}

	return 1
}

// Create a folder path with appropriate permissions
func CreateFolder(folderName string) {
	fancyprint.Infof("Creating folder %s\n", folderName)
	os.MkdirAll(folderName, 0777)
}

// Create a subfolder given a parent folder
func CreateSubFolder(parentFolder string, subFolder string) {
	folderPath := path.Join(parentFolder, subFolder)
	CreateFolder(folderPath)
}

// Get the user's home folder path
func GetHome() string {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		// fmt.Println(pair[0])

		if pair[0] == "HOME" || pair[0] == "USERPROFILE" {
			return pair[1]
		}
	}

	fancyprint.Warn("Warning! No home variable defined")
	return ""
}

// Convert a Date/Time string into a format that is a valid file path
func TimeToPath(timeStr string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(timeStr, ":", "-"), "/", "-"), " ", "T")
}

// Copy the source file to a destination file. Any existing file
// will be overwritten and will not copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// Given a working directory and repository, create a randomly named
// tar file with a dupver workdir configuration and a set of 10 50 MB
// random binary files
func CreateRandomTarFile(workDirFolder string, repoPath string) string {
	var fileName string
	fileName = "test_" + RandString(24, HexChars) + ".tar"
	WriteRandomTarFile(fileName, workDirFolder, repoPath)
	return fileName
}

// Given a filename, create a tar file with a dupver  workdir configuration
// and and a set of 10 50 MB random binary files
func WriteRandomTarFile(fileName string, workDirFolder string, repoPath string) {
	f, err := os.Create(fileName)

	if err != nil {
		panic(fmt.Sprintf("Error creating random tar file %s", fileName))
	}

	WriteRandomTar(f, workDirFolder, repoPath)
	// WriteRandom(f, 50000, 100, HexChars)
	f.Close()
}

// Given a file handle, create a tar file with a dupver workdir configuration
// and and a set of 10 50 MB random binary files
func WriteRandomTar(buf *os.File, workDirFolder string, repoPath string) {
	nFiles := 10

	tw := tar.NewWriter(buf)

	// ----------- Write containing folder ----------- //
	hdr := &tar.Header{
		Name: workDirFolder + "/",
		Mode: 0777,
		Size: int64(0),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatal(err)
	}

	// ----------- Write config folder ----------- //
	hdr = &tar.Header{
		Name: path.Join(workDirFolder, ".dupver") + "/",
		Mode: 0777,
		Size: int64(0),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatal(err)
	}

	// ----------- Write config file ----------- //
	var configPath string

	if len(workDirFolder) == 0 {
		configPath = path.Join(".dupver", "config.toml")
	} else {
		configPath = path.Join(workDirFolder, ".dupver", "config.toml")
	}

	workDirName := strings.ToLower(filepath.Base(workDirFolder))

	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")
		fancyprint.Warnf("Repo path not specified, setting to %s\n", repoPath)
	}

	// var myConfig WorkDirConfig
	// myConfig.RepoPath = repoPath
	// myConfig.WorkDirName = workDirName
	// myEncoder := toml.NewEncoder(tw)
	// myEncoder.Encode(myConfig)

	// WorkDirName = "ep-risk-characterization"
	// Branch = "main"
	// DefaultRepo = "main"

	// [Repos]
	//   main = "C:\\Users\\305232\\.dupver_repo"

	cfgStr := fmt.Sprintf("WorkDirName = \"%s\"\n", workDirName)
	cfgStr = cfgStr + "Branch = \"main\"\n"
	cfgStr = cfgStr + "DefaultRepo = \"main\"\n\n"
	cfgStr = cfgStr + fmt.Sprintf("[Repos]\n  main = \"%s\"", strings.Replace(repoPath, "\\", "\\\\", -1))
	bytes := []byte(cfgStr)

	hdr = &tar.Header{
		Name: configPath,
		Mode: 0600,
		Size: int64(len(bytes)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatal(err)
	}

	if _, err := tw.Write(bytes); err != nil {
		log.Fatal(err)
	}

	// ----------- Write random files ----------- //
	var files = []string{}

	for i := 0; i < nFiles; i += 1 {
		files = append(files, path.Join(workDirFolder, RandString(24, HexChars)+".bin"))
	}

	for _, file := range files {
		n := 80 * 1000 * 1000
		bytes := make([]byte, n)
		seededRand.Read(bytes)

		hdr := &tar.Header{
			Name: file,
			Mode: 0600,
			Size: int64(len(bytes)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write(bytes); err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}
}

// Given a filename, create a tar file with a dupver  workdir configuration
// and and a set of 10 50 MB random binary files
func CreateRandomFiles(workDirFolder string) ([]string, error) {
	nFiles := 10
	nBytes := 50 * 1000 * 1000
	fileNames := []string{}

	for i := 0; i < nFiles; i += 1 {
		fileName := path.Join(workDirFolder, RandString(24, HexChars)+".bin")
		fileNames = append(fileNames, fileName)
		err := CreateRandomFile(fileName, nBytes)

		if err != nil {
			return []string{}, err
		}
	}

	return fileNames, nil
}

func CreateRandomFile(fileName string, nBytes int) error {
	f, err := os.Create(fileName)

	if err != nil {
		return err
	}

	WriteRandomBytes(f, nBytes)
	f.Close()

	return nil
}

// Given a file handle, create a random binary file
func WriteRandomBytes(buf *os.File, nBytes int) error {
	bytes := make([]byte, nBytes)
	seededRand.Read(bytes)

	if _, err := buf.Write(bytes); err != nil {
		return err
	}


	if err := buf.Close(); err != nil {
		return err
	}

	return nil
}

func DiffFile(f1 string, f2 string) (bool, error) {
	dat1, err1 := ioutil.ReadFile(f1)

	if err1 != nil {
		return false, err1
	}

	dat2, err2 := ioutil.ReadFile(f2)

	if err2 != nil {
		return false, err2
	}	

	if len(dat1) != len(dat2) {
		return false, nil
	}

	for i := 0; i < len(dat1); i++ {
		x1 := dat1[i]
		x2 := dat2[i]

		if x1 != x2 {
			return false, nil
		}
	}

	return true, nil
}

// Print an object as JSON to stdout
func PrintJson(a interface{}) {
	myEncoder := json.NewEncoder(os.Stdout)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(a)
}

// Print an object as JSON to stdout
func MultiPrint(a interface{}, jsonOutput bool) {
	if jsonOutput {
		myEncoder := json.NewEncoder(os.Stdout)
		myEncoder.SetIndent("", "  ")
		myEncoder.Encode(a)
	} else {
		fmt.Println(a)
	}
}

func EditFile(filePath string, prefs Preferences) {
	editCmd := exec.Command(prefs.Editor, filePath)
	editCmd.Start()
}
