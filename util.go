package main

import (
    "os"
    // "io"
    "path"
    "bufio"
    "fmt"
    "log"
    "strings"
    "archive/tar"
	"github.com/BurntSushi/toml"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
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


func CreateRandomTarFile(workDirFolder string, repoPath string) string {
    var fileName string

    if len(workDirFolder) == 0 {
        fileName = "test_" + RandString(24, hexChars) + ".tar"
    } else {
        fileName = path.Join(workDirFolder, "test_" + RandString(24, hexChars) + ".tar")
    }

    WriteRandomTarFile(fileName, workDirFolder, repoPath)
    return fileName
}


func WriteRandomTarFile(fileName string, workDirFolder string, repoPath string) {
    f, err := os.Create(fileName)
    check(err)
    WriteRandomTar(f, workDirFolder, repoPath)
    // WriteRandom(f, 50000, 100, hexChars)
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

    // ----------- Write config file ----------- //
	var configPath string

	if len(workDirFolder) == 0 {
		 configPath = path.Join(".dupver", "config.toml")
	} else {
		configPath = path.Join(workDirFolder, ".dupver", "config.toml")
	}

	workDirName := strings.ToLower(path.Base(workDirFolder))

	if len(repoPath) == 0 {
		repoPath = path.Join(GetHome(), ".dupver_repo")
		fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
	}		

	var myConfig workDirConfig
	myConfig.RepoPath = repoPath
    myConfig.WorkDirName = workDirName
	myEncoder := toml.NewEncoder(tw)
    myEncoder.Encode(myConfig)    
    
    hdr := &tar.Header{
        Name: configPath,
        Mode: 0600,
        Size: int64(len(bytes)),
    }
    if err := tw.WriteHeader(hdr); err != nil {
        log.Fatal(err)
    }    

    
    // ----------- Write random files ----------- //    
    var files = []string{}

    for i := 0; i < nFiles; i += 1 {
        files = append(files, RandString(24, hexChars) + ".bin")
    }
        
    for _, file := range files {
        n := 50*1000*1000
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