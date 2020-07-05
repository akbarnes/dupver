package dupver

import (
	"os"
	"path/filepath"
	// "io"
	"archive/tar"
	"bufio"
	"fmt"
	"log"
	"path"
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

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func SetVerbosity(verbose bool, quiet bool) int {
	if quiet {
		return 0
	}

	if verbose {
		return 2
	}

	return 1
}

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

func TimeToPath(timeStr string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(timeStr, ":", "-"), "/", "-"), " ", "-")
}

func CreateRandomTarFile(workDirFolder string, repoPath string) string {
	var fileName string
	fileName = "test_" + RandString(24, HexChars) + ".tar"
	WriteRandomTarFile(fileName, workDirFolder, repoPath)
	return fileName
}

func WriteRandomTarFile(fileName string, workDirFolder string, repoPath string) {
	f, err := os.Create(fileName)

	if err != nil {
		panic(fmt.Sprintf("Error creating random tar file %s", fileName))
	}

	WriteRandomTar(f, workDirFolder, repoPath)
	// WriteRandom(f, 50000, 100, HexChars)
	f.Close()
}

func WriteRandomText(f *os.File, numLines int, numCols int, charset string) {
	w := bufio.NewWriter(f)

	b := make([]byte, numCols)

	for r := 0; r < numLines; r += 1 {
		for c := range b {
			b[c] = charset[seededRand.Intn(len(charset))]
		}

		fmt.Fprintf(w, "%s\n", b)
	}

	w.Flush()
}

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
	cfgStr := fmt.Sprintf("WorkDirName = \"%s\"\nRepoPath = \"%s\"", workDirName, repoPath)
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
