package main


import (
	"flag"
    "fmt"
	"os"
    "path"
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
	// constants
	ALL_REVISIONS := 0

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

	var revision int
	flag.IntVar(&revision, "revision", 0, "Specify revision number (default is last)")
	flag.IntVar(&revision, "n", 0, "Specify revision number(shorthand)")

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
		SaveWorkDirConfig(myConfig)
	} else if checkinFlag {
        set_default(&repoPath, "$HOME/.dupver_repo")
		snapshotsPath := path.Join(repoPath, "snapshots")
		os.Mkdir(snapshotsPath, 0777)
        snapshotBasename := RandHexString(65)
        var snapshotPath string

        if len(workDirName) == 0 {
			snapshotPath = path.Join(snapshotsPath, snapshotBasename + ".toml")
        } else {
			snapshotPath = path.Join(snapshotsPath, workDirName, snapshotBasename + ".toml")
        }

		mypoly := 0x3DA3358B4DC173
		fmt.Println("Backing up ", filePath)
		snapshotFile, _ := os.Create(snapshotPath)
		PrintCommitHeader(snapshotFile, msg, filePath)
		// also save hashes for tar file to check which files are modified
		PrintTarFileIndex(filePath, snapshotFile)
		PackFile(filePath, repoPath, snapshotFile, mypoly)
		snapshotFile.Close()
	} else if checkoutFlag {
		if len(repoPath) == 0 { 
            repoPath = "$HOME/.dupver_repo"
        }

		commitHistoryPath := path.Join(repoPath, "commits.toml")
		history := ReadHistory(commitHistoryPath)
		fmt.Printf("Number of commits %d\n", len(history.Commits))
		revIndex := GetRevIndex(revision, len(history.Commits))
		fmt.Printf("Restoring commit %d\n", revIndex)
		
		if (true || len(filePath) == 0) {
			filePath = fmt.Sprintf("snapshot%d.tar", revIndex + 1)
		}

		fmt.Printf("Writing to %s\n", filePath)
		UnpackTar(filePath, history.Commits[revIndex].Chunks) 
	} else if listFlag {
		if len(repoPath) == 0 { 
            repoPath = "$HOME/.dupver_repo"
        }

		commitHistoryPath := path.Join(repoPath, "commits.toml")
		history := ReadHistory(commitHistoryPath)

		// print a specific revision
		if revision == ALL_REVISIONS {
			fmt.Printf("Commit History\n")

			for i:=0; i < len(history.Commits); i++ {
				PrintRevision(history, i, 10)
			}			
		} else {
			revIndex := GetRevIndex(revision, len(history.Commits))
			PrintRevision(history, revIndex, 0)
		}
	} else {
		fmt.Println("No command specified, exiting")
		fmt.Println("For available commands run: dupver -help")
	}
}
