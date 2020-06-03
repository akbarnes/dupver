package main


import (
	"flag"
    "fmt"
	"os"
	"strings"
	"time"
	"path"
	"path/filepath"
	"github.com/BurntSushi/toml"
)


func check(e error) {
    if e != nil {
        panic(e)
    }
}


func set_default(s *string, d string) {
	if len(*s) == 0 {
		*s = d
	}
}


func main() {
	SNAPSHOT_ID_LEN := 40 

	var initRepoFlag bool
	var initWorkDirFlag bool
	var checkinFlag bool
	var checkoutFlag bool
	var listFlag bool

	flag.BoolVar(&initRepoFlag, "init-repo", false, "Initialize the repository")
	flag.BoolVar(&initWorkDirFlag, "init", false, "Initialize the working directory")
	flag.BoolVar(&checkinFlag, "checkin", false, "Check in specified file")
	flag.BoolVar(&checkinFlag, "ci", false, "Check in specified file")


	flag.BoolVar(&checkoutFlag, "checkout", false, "Check out specified file")
	flag.BoolVar(&checkoutFlag, "co", false, "Check out specified file")

	flag.BoolVar(&listFlag, "list", false, "List revisions")


	var filePath string
	flag.StringVar(&filePath, "file", "", "Archive path")
	flag.StringVar(&filePath, "f", "", "Archive path (shorthand)")

	var snapshot string
	flag.StringVar(&snapshot, "snapshot", "", "Specify snapshot id (default is last)")
	flag.StringVar(&snapshot, "s", "", "Specify snapshot id (shorthand)")

	var msg string
	flag.StringVar(&msg, "message", "", "Commit message")
	flag.StringVar(&msg, "m", "", "Commit message (shorthand)")

	var repoPath string
	flag.StringVar(&repoPath, "repository", "", "Repository path")
	flag.StringVar(&repoPath, "r", "", "Repository path (shorthand)")

	var workDirName string
	flag.StringVar(&workDirName, "workdir", "", "Working directory name")
	flag.StringVar(&workDirName, "w", "", "Working directory name (shorthand)")

	flag.Parse()
	

	if initRepoFlag {
        set_default(&repoPath, "$HOME/.dupver_repo")
        fmt.Printf("Creating folder %s\n", repoPath)
		os.Mkdir(repoPath, 0777)

		packPath := path.Join(repoPath, "packs")
        fmt.Printf("Creating folder %s\n", packPath)
		os.Mkdir(packPath, 0777)

	    snapshotsPath := path.Join(repoPath, "snapshots")
        fmt.Printf("Creating folder %s\n", snapshotsPath)
		os.Mkdir(snapshotsPath, 0777)

		var myConfig repoConfig
		myConfig.Version = 1
		myConfig.ChunkerPolynomial = 0x3DA3358B4DC173
		SaveRepoConfig(repoPath, myConfig)
	} else if initWorkDirFlag {
        set_default(&repoPath, "$HOME/.dupver_repo")
		os.Mkdir("./.dupver", 0777)
		var myConfig workDirConfig
		myConfig.RepositoryPath = repoPath
		myConfig.WorkDirName = workDirName
		SaveWorkDirConfig(myConfig)
	} else if checkinFlag {
		t := time.Now()

		var mySnapshot commit

        set_default(&repoPath, "$HOME/.dupver_repo")
		snapshotsPath := path.Join(repoPath, "snapshots")
		os.Mkdir(snapshotsPath, 0777)
        snapshotId := RandHexString(SNAPSHOT_ID_LEN)
        mySnapshot.ID = snapshotId
		snapshotDate := t.Format("2006-01-02-T15-04-05")
		mySnapshot.Time = snapshotDate
        snapshotBasename := fmt.Sprintf("%s-%s", snapshotDate, snapshotId[0:40])

        var snapshotPath string
		if len(workDirName) == 0 {
			panic("WorkDirName not specified")
		} 

        snapshotFolder := path.Join(repoPath, "snapshots", workDirName)
		os.Mkdir(snapshotFolder, 0777)
		snapshotPath = path.Join(snapshotFolder, snapshotBasename + ".toml")
		mypoly := 0x3DA3358B4DC173
		fmt.Println("Backing up ", filePath)
		mySnapshot.TarFileName = filePath

		if len(msg) == 0 {
			msg =  strings.Replace(filePath[0:len(filePath)-4], ".\\", "", -1)
		}

		mySnapshot.Message = msg

		// also save hashes for tar file to check which files are modified
		mySnapshot.Files, _ = ReadTarFileIndex(filePath)
		mySnapshot.Chunks = PackFile(filePath, repoPath, mypoly)

		snapshotFile, _ := os.Create(snapshotPath)
		myEncoder := toml.NewEncoder(snapshotFile)
		myEncoder.Encode(mySnapshot)
		snapshotFile.Close()
	} else if checkoutFlag {
        set_default(&repoPath, "$HOME/.dupver_repo")
		snapshotGlob := path.Join(repoPath, "snapshots", workDirName, "*.toml")
		fmt.Println(snapshotGlob)
		snapshotPaths, _ := filepath.Glob(snapshotGlob)

		var mySnapshot commit
		foundSnapshot := false

		for _, snapshotPath := range snapshotPaths {
			n := len(snapshotPath)
			snapshotId := snapshotPath[n-SNAPSHOT_ID_LEN-5:n-5]
		
			if snapshotId[0:len(snapshot)] == snapshot {
				mySnapshot = ReadSnapshot(snapshotPath)
				foundSnapshot = true
				break
			}
		}
		
		if !foundSnapshot {
			panic("Could not find snapshot")
		}

		if len(filePath) == 0 {
			filePath = fmt.Sprintf("%s-%s.tar", workDirName, snapshot)
		}

		UnpackTar(filePath, repoPath, mySnapshot.Chunks) 
		fmt.Printf("Wrote to %s\n", filePath)
	} else if listFlag {
		if len(repoPath) == 0 { 
            repoPath = "$HOME/.dupver_repo"
        }

		snapshotGlob := path.Join(repoPath, "snapshots", workDirName, "*.toml")
		fmt.Println(snapshotGlob)
		snapshotPaths, _ := filepath.Glob(snapshotGlob)

		// print a specific revision
		if len(snapshot) == 0 {
			fmt.Printf("Snapshot History\n")

			for _, snapshotPath := range snapshotPaths {
				fmt.Printf("Path: %s\n", snapshotPath)
				PrintSnapshot(ReadSnapshot(snapshotPath), 10)
			}			
		} else {
			fmt.Println("Snapshot")

			for _, snapshotPath := range snapshotPaths {
				n := len(snapshotPath)
				snapshotId := snapshotPath[n-SNAPSHOT_ID_LEN-5:n-5]
			
				if snapshotId[0:8] == snapshot {
					PrintSnapshot(ReadSnapshot(snapshotPath), 0)
				}
			}	

		}
	} else {
		fmt.Println("No command specified, exiting")
		fmt.Println("For available commands run: dupver -help")
	}
}
