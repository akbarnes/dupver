package main


import (
	"flag"
    "fmt"
	"os"
	"strings"
	"time"
	"path"
	"path/filepath"
	// "github.com/BurntSushi/toml"
)


func check(e error) {
    if e != nil {
        panic(e)
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

	var workDir string
	flag.StringVar(&workDir, "workdir", "", "Working directory")
	flag.StringVar(&workDir, "d", "", "Working directory (shorthand)")

	var workDirName string
	flag.StringVar(&workDirName, "workdir-name", "", "Working directory name")
	flag.StringVar(&workDirName, "w", "", "Working directory name (shorthand)")

	flag.Parse()
	

	if initRepoFlag {
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
		var configPath string

		if len(workDir) == 0 {
 			os.Mkdir(".dupver", 0777)
 			configPath = path.Join(".dupver", "config.toml")
		} else {
 			os.Mkdir(path.Join(workDir, ".dupver"), 0777)
 			configPath = path.Join(workDir, ".dupver", "config.toml")
		}

		if len(workDirName) == 0 {
			if len(workDir) == 0 {
				panic("Both workDir and workDirName are empty")
			} else {
				workDirName = strings.ToLower(path.Base(workDir))
			}
		}

		var myConfig workDirConfig
		myConfig.RepositoryPath = repoPath
		myConfig.WorkDirName = workDirName
		SaveWorkDirConfig(configPath, myConfig)
	} else if checkinFlag {
		t := time.Now()

		var mySnapshot commit
		var myWorkDirConfig workDirConfig

		snapshotsPath := path.Join(repoPath, "snapshots")
		os.Mkdir(snapshotsPath, 0777)
        snapshotId := RandHexString(SNAPSHOT_ID_LEN)
        mySnapshot.ID = snapshotId
		snapshotDate := t.Format("2006-01-02-T15-04-05")
		mySnapshot.Time = t.Format("2006/01/02 15:04:05")
        snapshotBasename := fmt.Sprintf("%s-%s", snapshotDate, snapshotId[0:40])
		mySnapshot.Files, myWorkDirConfig = ReadTarFileIndex(filePath)

        var snapshotPath string
		if len(workDirName) == 0 {
			workDirName = myWorkDirConfig.WorkDirName
		} 

		if len(repoPath) == 0 {
			repoPath = myWorkDirConfig.RepositoryPath
		}

		fmt.Printf("Workdir name: %s\nRepo path: %s\n", workDirName, repoPath)

        snapshotFolder := path.Join(repoPath, "snapshots", workDirName)
		os.Mkdir(snapshotFolder, 0777)
		snapshotPath = path.Join(snapshotFolder, snapshotBasename + ".toml")
		mypoly := 0x3DA3358B4DC173
		fmt.Printf("Checking in %s as snapshot %s\n", filePath, snapshotId[0:8])
		mySnapshot.TarFileName = filePath

		if len(msg) == 0 {
			msg =  strings.Replace(filePath[0:len(filePath)-4], ".\\", "", -1)
		}

		mySnapshot.Message = msg

		// also save hashes for tar file to check which files are modified
		mySnapshot.Chunks = PackFile(filePath, repoPath, mypoly)

		WriteSnapshot(snapshotPath, mySnapshot)
	} else if checkoutFlag {
		var configPath string

		if len(workDir) == 0 {
			configPath = path.Join(".dupver", "config.toml")
		} else {
			configPath = path.Join(workDir, ".dupver", "config.toml")
		}

		myWorkDirConfig := ReadWorkDirConfig(configPath)

		if len(workDirName) == 0 {
			workDirName = myWorkDirConfig.WorkDirName
		}

		if len(repoPath) == 0 { 
            repoPath = myWorkDirConfig.RepositoryPath
        }

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
			filePath = fmt.Sprintf("%s-%s-%s.tar", workDirName, mySnapshot.Time, snapshot)
		}

		UnpackTar(filePath, repoPath, mySnapshot.Chunks) 
		fmt.Printf("Wrote to %s\n", filePath)
	} else if listFlag {
		var configPath string

		if len(workDir) == 0 {
			configPath = path.Join(".dupver", "config.toml")
		} else {
			configPath = path.Join(workDir, ".dupver", "config.toml")
		}

		myWorkDirConfig := ReadWorkDirConfig(configPath)

		if len(workDirName) == 0 {
			workDirName = myWorkDirConfig.WorkDirName
		}

		if len(repoPath) == 0 { 
            repoPath = myWorkDirConfig.RepositoryPath
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
