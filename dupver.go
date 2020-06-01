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


func main() {
	// constants
	mypoly := 0x3DA3358B4DC173
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
	var msg string
	var revision int
	var repoPath string

	flag.StringVar(&filePath, "file", "", "Archive path")
	flag.StringVar(&filePath, "f", "", "Archive path (shorthand)")

	flag.IntVar(&revision, "revision", 0, "Specify revision number (default is last)")
	flag.IntVar(&revision, "n", 0, "Specify revision number(shorthand)")

	flag.StringVar(&msg, "message", "", "Commit message")
	flag.StringVar(&msg, "m", "", "Commit message (shorthand)")

	flag.StringVar(&repoPath, "repository", "", "Repository path")
	flag.StringVar(&repoPath, "r", "", "Repository path (shorthand)")

	flag.Parse()
	
	commitHistoryPath := path.Join(repoPath, "commits.toml")

	if initRepoFlag {
        if len(repoPath) == 0 { 
            repoPath = "$HOME/.dupver_repo"
        }

	    commitHistoryPath = path.Join(repoPath, "commits.toml")
        fmt.Printf("Creating folder %s\n", repoPath)
		os.Mkdir(repoPath, 0777)
        fmt.Printf("Creating commit history %s\n", commitHistoryPath)
		f, _ := os.Create(commitHistoryPath)
		f.Close()
	} else if initWorkDirFlag {
		os.Mkdir("./.dupver", 0777)
		var myConfig workDirConfig
		myConfig.RepositoryPath = repoPath
		SaveWorkDirConfig(myConfig)
	} else if checkinFlag {
		fmt.Println("Backing up ", filePath)
		commitFile, _ := os.OpenFile(commitHistoryPath, os.O_APPEND|os.O_WRONLY, 0600)
		PrintCommitHeader(commitFile, msg, filePath)
		PrintTarFileIndex(filePath, commitFile)
		PackFile(filePath, commitFile, mypoly)
		commitFile.Close()
	} else if checkoutFlag {
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
