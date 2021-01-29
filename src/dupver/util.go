package dupver

import (
	"archive/tar"
	// "bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	// "github.com/BurntSushi/toml"
)

const colorReset string = "\033[0m"
const colorRed string = "\033[31m"
const colorGreen string = "\033[32m"
const colorYellow string = "\033[33m"
const colorBlue string = "\033[34m"
const colorPurple string = "\033[35m"
const colorCyan string = "\033[36m"
const colorWhite string = "\033[37m"

type Options struct {
	Verbosity    int
	Color        bool
	WorkDirName  string
	RepoName     string
	RepoPath     string
	Branch       string
	DestRepoName string
	DestRepoPath string
}

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

// Set the verbosity level given command-line flags
func SetVerbosity(opts Options, debug bool, verbose bool, quiet bool) Options {
	opts.Verbosity = CalculateVerbosity(debug, verbose, quiet)
	return opts
}

// Create a folder path
func CreateFolder(folderName string, verbosity int) {
	if verbosity >= 1 {
		fmt.Printf("Creating folder %s\n", folderName)
	}
	os.MkdirAll(folderName, 0777)
}

// Create a subfolder given a parent folder
func CreateSubFolder(parentFolder string, subFolder string, verbosity int) {
	folderPath := path.Join(parentFolder, subFolder)
	CreateFolder(folderPath, verbosity)
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

	fmt.Println("Warning! No home variable defined")
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
		fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
	}

	// var myConfig workDirConfig
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
		n := 50 * 1000 * 1000
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
