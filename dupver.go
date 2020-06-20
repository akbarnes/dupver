package main


import (
	"flag"
	"fmt"
	// "os"
	// "github.com/google/subcommands"
)

const version string = "0.2.0-alpha"


func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "Archive path")
	flag.StringVar(&filePath, "f", "", "Archive path (shorthand)")

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
	posArgs := flag.Args()
	cmd := posArgs[0]
	
  	if cmd == "init-repo" {
		repoPath := posArgs[1]
		InitRepo(repoPath)
	} else if cmd == "init" {
		if len(posArgs) >= 2 {
			workDir = posArgs[1]
		}

		// Read repoPath from environment variable if empty
		InitWorkDir(workDir, workDirName, repoPath)
	} else if cmd == "commit" || cmd == "ci" {
		commitFile := posArgs[1]
		CommitFile(commitFile, msg, verbosity)
	} else if cmd == "checkout" || cmd == "co" {
		snapshotId := posArgs[1]

		cfg := ReadWorkDirConfig(workDir)
		cfg = UpdateWorkDirName(cfg, workDirName)
		cfg = UpdateRepoPath(cfg, repoPath)
		snap := ReadSnapshot(snapshotId, cfg)

		if len(filePath) == 0 {
			timeStr := TimeToPath(snap.Time)
			filePath = fmt.Sprintf("%s-%s-%s.tar", cfg.WorkDirName, timeStr, snapshotId[0:16])
		}

		UnpackFile(filePath, cfg.RepoPath, snap.ChunkIDs, verbosity) 
		fmt.Printf("Wrote to %s\n", filePath)
	} else if cmd == "log" || cmd == "list" {
		snapshotId := ""

		if len(posArgs) >= 2 {
			snapshotId = posArgs[1]
		}

		cfg := ReadWorkDirConfig(workDir)
		cfg = UpdateWorkDirName(cfg, workDirName)
		cfg = UpdateRepoPath(cfg, repoPath)
		PrintSnapshots(ListSnapshots(cfg), snapshotId)
	} else if cmd == "status" || cmd == "st" {
		snapshotId := ""

		if len(posArgs) >= 2 {
			snapshotId = posArgs[1]
		}

		cfg := ReadWorkDirConfig(workDir)
		cfg = UpdateWorkDirName(cfg, workDirName)
		cfg = UpdateRepoPath(cfg, repoPath)
		snapshotPaths := ListSnapshots(cfg)
		
		for _, snapshotPath := range snapshotPaths {
			n := len(snapshotPath)
			sid := snapshotPath[n-SNAPSHOT_ID_LEN-5:n-5]
		
			if sid[0:8] == snapshotId {
				mySnapshot := ReadSnapshotFile(snapshotPath)
				WorkDirStatus(workDir, mySnapshot)
				break
			}

		}			
	} else if cmd == "version" {
		fmt.Println("Dupver version:", version)
	} else {
		fmt.Println("No command specified, exiting")
		fmt.Println("For available commands run: dupver -help")
	}
}
