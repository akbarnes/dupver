package main


import (
	"flag"
    "fmt"
	"os"
	"strings"
	"time"
	"path"
	// "path/filepath"
	// "github.com/BurntSushi/toml"
	"github.com/restic/chunker"
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
		if len(repoPath) == 0 {
			repoPath = path.Join(GetHome(), ".dupver_repo")
			fmt.Printf("Repo path not specified, setting to %s\n", repoPath)
		}	
				
		// InitRepo(workDir)
		fmt.Printf("Creating folder %s\n", repoPath)
		os.Mkdir(repoPath, 0777)
	
		packPath := path.Join(repoPath, "packs")
		fmt.Printf("Creating folder %s\n", packPath)
		os.Mkdir(packPath, 0777)
	
		snapshotsPath := path.Join(repoPath, "snapshots")
		fmt.Printf("Creating folder %s\n", snapshotsPath)
		os.MkdirAll(snapshotsPath, 0777)
	
		treesPath := path.Join(repoPath, "trees")
		fmt.Printf("Creating folder %s\n", treesPath)
		os.Mkdir(treesPath, 0777)	

		p, err := chunker.RandomPolynomial()
		check(err)
	
		var myConfig repoConfig
		myConfig.Version = 1
		myConfig.ChunkerPolynomial = p
		SaveRepoConfig(repoPath, myConfig)
	} else if initWorkDirFlag {
		InitWorkDir(workDir, workDirName, repoPath)
	} else if checkinFlag {
		var myWorkDirConfig workDirConfig
		t := time.Now()

		var mySnapshot commit
        mySnapshot.ID = RandHexString(SNAPSHOT_ID_LEN)
		mySnapshot.Time = t.Format("2006/01/02 15:04:05")
		mySnapshot.TarFileName = filePath
		// mySnapshot = UpdateTags(mySnapshot, tagName)
		mySnapshot = UpdateMessage(mySnapshot, msg, filePath)		
		mySnapshot.Files, myWorkDirConfig = ReadTarFileIndex(filePath)
		myRepoConfig := ReadRepoConfigFile(path.Join(myWorkDirConfig.RepoPath, "config.toml"))
		
		chunkIDs, chunkPacks := PackFile(filePath, myWorkDirConfig.RepoPath, myRepoConfig.ChunkerPolynomial, verbosity)
		mySnapshot.ChunkIDs = chunkIDs

		snapshotFolder := path.Join(myWorkDirConfig.RepoPath, "snapshots", myWorkDirConfig.WorkDirName)
        snapshotBasename := fmt.Sprintf("%s-%s", t.Format("2006-01-02-T15-04-05"), mySnapshot.ID[0:40])		
		os.Mkdir(snapshotFolder, 0777)
		snapshotPath := path.Join(snapshotFolder, snapshotBasename + ".json")
		WriteSnapshot(snapshotPath, mySnapshot)

		treeFolder := path.Join(myWorkDirConfig.RepoPath, "trees")
        treeBasename := mySnapshot.ID[0:40]
		os.Mkdir(treeFolder, 0777)
		treePath := path.Join(treeFolder, treeBasename + ".json")
		WriteTree(treePath, chunkPacks)
	} else if checkoutFlag {
		myWorkDirConfig := ReadWorkDirConfig(workDir)
		myWorkDirConfig = UpdateWorkDirName(myWorkDirConfig, workDirName)
		myWorkDirConfig = UpdateRepoPath(myWorkDirConfig, repoPath)
		mySnapshot := ReadSnapshot(snapshot, myWorkDirConfig)

		if len(filePath) == 0 {
			timeStr := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(mySnapshot.Time, ":", "-"), "/", "-"), " ", "-")
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
