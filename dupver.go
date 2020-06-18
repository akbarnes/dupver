package main


import (
	"flag"
    "fmt"
	// "os"
	// "strings"
	// "time"
	// "path"
	// "path/filepath"
	// "github.com/BurntSushi/toml"
	// "github.com/restic/chunker"
)

const version string = "0.1.0-alpha"

func main() {
	var initRepoFlag bool
	var initWorkDirFlag bool
	var checkinFlag bool
	var checkoutFlag bool
	var listFlag bool
	var versionFlag bool

	flag.BoolVar(&initRepoFlag, "init-repo", false, "Initialize the repository")
	flag.BoolVar(&initWorkDirFlag, "init", false, "Initialize the working directory")
	flag.BoolVar(&checkinFlag, "commit", false, "Commit specified file")
	flag.BoolVar(&checkinFlag, "ci", false, "Commit specified file (shorthand)")
	// flag.BoolVar(&tagFlag, "tag", false, "Tag specified commit (shorthand)")


	flag.BoolVar(&checkoutFlag, "checkout", false, "Check out specified file")
	flag.BoolVar(&checkoutFlag, "co", false, "Check out specified file")

	flag.BoolVar(&listFlag, "list", false, "List revisions")
	flag.BoolVar(&listFlag, "ls", false, "List revisions (shorthand)")


	flag.BoolVar(&versionFlag, "version", false, "Print version number")
	flag.BoolVar(&versionFlag, "V", false, "Print version number (shorthand)")

	var filePath string
	flag.StringVar(&filePath, "file", "", "Archive path")
	flag.StringVar(&filePath, "f", "", "Archive path (shorthand)")

	var snapshot string
	flag.StringVar(&snapshot, "snapshot", "", "Specify snapshot id (default is last)")
	flag.StringVar(&snapshot, "s", "", "Specify snapshot id (shorthand)")
	flag.StringVar(&snapshot, "commit-id", "", "Specify commit (snaphot) id (default is last)")
	flag.StringVar(&snapshot, "c", "", "Specify commit (snapshot) id (shorthand)")

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

	var tagName string
	flag.StringVar(&tagName, "tag-name", "", "Tag name")
	flag.StringVar(&tagName, "t", "", "Tag name (shorthand)")

	var verbosity int
	flag.IntVar(&verbosity, "verbosity", 1, "Verbosity level")
	flag.IntVar(&verbosity, "v", 1, "Verbosity level (shorthand)")	

	flag.Parse()
	

	if initRepoFlag {
		InitRepo(repoPath)
	} else if initWorkDirFlag {
		InitWorkDir(workDir, workDirName, repoPath)
	} else if checkinFlag {
		CommitFile(filePath, msg, verbosity)
	} else if checkoutFlag {
		myWorkDirConfig := ReadWorkDirConfig(workDir)
		myWorkDirConfig = UpdateWorkDirName(myWorkDirConfig, workDirName)
		myWorkDirConfig = UpdateRepoPath(myWorkDirConfig, repoPath)
		mySnapshot := ReadSnapshot(snapshot, myWorkDirConfig)

		if len(filePath) == 0 {
			timeStr := TimeToPath(mySnapshot.Time)
			filePath = fmt.Sprintf("%s-%s-%s.tar", myWorkDirConfig.WorkDirName, timeStr, snapshot[0:16])
		}

		UnpackFile(filePath, myWorkDirConfig.RepoPath, mySnapshot.ChunkIDs, verbosity) 
		fmt.Printf("Wrote to %s\n", filePath)
	} else if listFlag {
		myWorkDirConfig := ReadWorkDirConfig(workDir)
		myWorkDirConfig = UpdateWorkDirName(myWorkDirConfig, workDirName)
		myWorkDirConfig = UpdateRepoPath(myWorkDirConfig, repoPath)
		PrintSnapshots(ListSnapshots(myWorkDirConfig), snapshot)
	} else if versionFlag {
		fmt.Println("Dupver version:", version)
	} else {
		fmt.Println("No command specified, exiting")
		fmt.Println("For available commands run: dupver -help")
	}
}
