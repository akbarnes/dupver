package main


import (
	"flag"
    "fmt"
	"os"
	// "log"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func main() {
	// constants
	mypoly := 0x3DA3358B4DC173
	commitLogPath := fmt.Sprintf(".dupver/versions.toml")
	ALL_REVISIONS := 0

	var initFlag bool
	var checkinFlag bool 
	var checkoutFlag bool
	var listFlag bool

	flag.BoolVar(&initFlag, "init", false, "Initialize the repository")
	flag.BoolVar(&checkinFlag, "checkin", false, "Check in specified file")
	flag.BoolVar(&checkinFlag, "ci", false, "Check in specified file")


	flag.BoolVar(&checkoutFlag, "checkout", false, "Check out specified file")
	flag.BoolVar(&checkoutFlag, "co", false, "Check out specified file")

	flag.BoolVar(&listFlag, "list", false, "List revisions")

	var filePath string
	var msg string
	var revision int
	var repo string

	flag.StringVar(&filePath, "file", "", "Archive path")
	flag.StringVar(&filePath, "f", "", "Archive path (shorthand)")

	flag.IntVar(&revision, "revision", 0, "Specify revision (default is last)")
	flag.IntVar(&revision, "r", 0, "Specify revision (shorthand)")

	flag.StringVar(&msg, "message", "", "Commit message")
	flag.StringVar(&msg, "m", "", "Commit message (shorthand)")

	flag.StringVar(&repo, "repository", "", "Commit message")
	flag.StringVar(&repo, "repo", "", "Commit message")
	
	flag.Parse()
	

	if initFlag {
		os.Mkdir("./.dupver", 0777)
		f, _ := os.Create(commitLogPath)
		f.Close()
	} else if checkinFlag {
		fmt.Println("Backing up ", filePath)
		commitFile, _ := os.OpenFile(commitLogPath, os.O_APPEND|os.O_WRONLY, 0600)
		PrintCommitHeader(commitFile, msg, filePath)
		PrintTarIndex(filePath, commitFile)
		PackTar(filePath, commitFile, mypoly)
		commitFile.Close()
	} else if checkoutFlag {
		history := ReadHistory(commitLogPath)

		fmt.Printf("Number of commits %d\n", len(history.Commits))
		revIndex := GetRevIndex(revision, len(history.Commits))
		fmt.Printf("Restoring commit %d\n", revIndex)
		
		if (true || len(filePath) == 0) {
			filePath = fmt.Sprintf("snapshot%d.tar", revIndex + 1)
		}

		fmt.Printf("Writing to %s\n", filePath)
		UnpackTar(filePath, history.Commits[revIndex].Chunks) 
	} else if listFlag {
		history := ReadHistory(commitLogPath)

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
