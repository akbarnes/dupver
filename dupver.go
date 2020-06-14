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
)


func main() {
	var initRepoFlag bool
	var initWorkDirFlag bool
	var checkinFlag bool
	var checkoutFlag bool
	var listFlag bool

	flag.BoolVar(&initRepoFlag, "init-repo", false, "Initialize the repository")
	flag.BoolVar(&initWorkDirFlag, "init", false, "Initialize the working directory")
	flag.BoolVar(&checkinFlag, "commit", false, "Commit specified file")
	flag.BoolVar(&checkinFlag, "ci", false, "Commit specified file (shorthand)")
	// flag.BoolVar(&tagFlag, "tag", false, "Tag specified commit (shorthand)")


	flag.BoolVar(&checkoutFlag, "checkout", false, "Check out specified file")
	flag.BoolVar(&checkoutFlag, "co", false, "Check out specified file")

	flag.BoolVar(&listFlag, "list", false, "List revisions")
	flag.BoolVar(&listFlag, "ls", false, "List revisions (shorthand)")

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

	flag.Parse()
	

	if initRepoFlag {
		InitRepo(workDir)
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
		myConfig.RepoPath = repoPath
		myConfig.WorkDirName = workDirName
		SaveWorkDirConfig(configPath, myConfig)
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
		// mySnapshot.Packs, mySnapshot.Chunks = PackFile(filePath, myWorkDirConfig.RepoPath, 0x3DA3358B4DC173)
		chunkIDs, chunkPacks := PackFile(filePath, myWorkDirConfig.RepoPath, 0x3DA3358B4DC173)
		mySnapshot.ChunkIDs = chunkIDs
		// mySnapshot.ChunkPacks = chunkPacks

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

		UnpackFile(filePath, myWorkDirConfig.RepoPath, mySnapshot.ChunkIDs) 
		fmt.Printf("Wrote to %s\n", filePath)
	} else if listFlag {
		myWorkDirConfig := ReadWorkDirConfig(workDir)
		myWorkDirConfig = UpdateWorkDirName(myWorkDirConfig, workDirName)
		myWorkDirConfig = UpdateRepoPath(myWorkDirConfig, repoPath)
		PrintSnapshots(ListSnapshots(myWorkDirConfig), snapshot)
	} else {
		fmt.Println("No command specified, exiting")
		fmt.Println("For available commands run: dupver -help")
	}
}
