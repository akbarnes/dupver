package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/akbarnes/dupver/src/dupver"
)

func TestWorkRepoInit(t *testing.T) {
	opts := dupver.Options{Color: false, Verbosity: 2}
	homeDir := dupver.GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	repoId := dupver.RandString(16, dupver.HexChars)
	repoFolder := ".dupver_repo_" + repoId
	repoName := "test"

	repoPath := filepath.Join(homeDir, "temp", repoFolder)
	dupver.InitRepo(repoPath, repoName, "", opts)

	snapshotsPath := filepath.Join(repoPath, "snapshots")
	if _, err := os.Stat(snapshotsPath); err != nil {
		// path/to/whatever exists
		t.Error("Did not create snapshots folder", snapshotsPath)
	}

	treesPath := filepath.Join(repoPath, "trees")
	if _, err := os.Stat(treesPath); err != nil {
		// path/to/whatever exists
		t.Error("Did not create trees folder", treesPath)
	}

	cfg := dupver.ReadRepoConfigFile(filepath.Join(repoPath, "config.toml"))

	if cfg.Version != 2 {
		t.Error("Invalid repository version", cfg.Version)
	}

	if cfg.ChunkerPolynomial <= 0 {
		t.Error("Invalid chunker polynomial", cfg.ChunkerPolynomial)
	}

	os.RemoveAll(repoPath)
}

func TestWorkDirInit(t *testing.T) {
	opts := dupver.Options{Color: false, Verbosity: 2}
	homeDir := dupver.GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	workDirId := dupver.RandString(16, dupver.HexChars)
	workDirFolder := "Test_" + workDirId
	// workDirPath := filepath.Join("temp", workDirFolder)
	err := os.MkdirAll(workDirFolder, 0777)

	if err != nil {
		t.Error("Could not create workdir")
	}

	repoPath := filepath.Join(homeDir, ".dupver_repo")

	// workDirPath := ""
	projectName := ""

	opts.RepoName = "test"
	opts.RepoPath = repoPath
	opts.Branch = "main"

	dupver.InitWorkDir(workDirFolder, projectName, opts)

	cfg := dupver.ReadWorkDirConfig(workDirFolder)

	if cfg.Repos[cfg.DefaultRepo] != repoPath {
		t.Error("Incorrect repo path retrieved")
	}

	expectedWorkDirName := "test_" + workDirId
	if cfg.WorkDirName != expectedWorkDirName {
		t.Error("Incorrect workdir name retrieved")
	}

	os.RemoveAll(workDirFolder)
}

func TestCommit(t *testing.T) {
	opts := dupver.Options{Color: false, Verbosity: 2}
	msg := "Commit random data"

	// ----------- Create a repo ----------- //
	homeDir := dupver.GetHome()
	repoId := dupver.RandString(16, dupver.HexChars)
	repoFolder := ".dupver_repo_" + repoId
	repoPath := filepath.Join(homeDir, "temp", repoFolder)
	repoName := "test"
	dupver.InitRepo(repoPath, repoName, "", opts)

	// ----------- Create a workdir ----------- //
	workDirId := dupver.RandString(16, dupver.HexChars)
	workDirFolder := "Test_" + workDirId
	err := os.MkdirAll(workDirFolder, 0777)

	if err != nil {
		t.Error("Could not cerate workdir folder " + workDirFolder)
	}

	projectName := ""

	opts.WorkDirName = "test"
	opts.RepoName = "test"
	opts.RepoPath = repoPath
	opts.Branch = "main"

	dupver.InitWorkDir(workDirFolder, projectName, opts)

	// ----------- Create tar file with random data ----------- //
	// TODO: add random permutes to data
	fileName := dupver.CreateRandomTarFile(workDirFolder, repoPath)
	fmt.Printf("Created tar file %s\n", fileName)

	// ----------- Commit the tar file  ----------- //
	snapshot := dupver.CommitFile(fileName, []string{}, msg, opts)

	// ----------- Commit the tar file  ----------- //
	myWorkDirConfig := dupver.ReadWorkDirConfig(workDirFolder)
	dupver.PrintAllSnapshots("", opts)

	fmt.Printf("snapshot: %+v\n\n", snapshot)
	fmt.Printf("workdir config: %+v\n\n", myWorkDirConfig)

	// // ----------- Checkout the tar file  ----------- //
	mySnapshot := dupver.ReadSnapshot(snapshot.ID, opts)
	timeStr := dupver.TimeToPath(mySnapshot.Time)
	outputFileName := fmt.Sprintf("%s-%s-%s.tar", myWorkDirConfig.WorkDirName, timeStr, snapshot.ID[0:16])

	// dupver.UnpackFile(outputFileName, myWorkDirConfig.RepoPath, mySnapshot.ChunkIDs, verbosity)
	dupver.UnpackFile(outputFileName, opts.RepoPath, mySnapshot.ChunkIDs, opts)
	fmt.Printf("Wrote to %s\n", outputFileName)

	cmd := exec.Command("diff", fileName, outputFileName)
	log.Printf("Running command and waiting for it to finish...")
	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("diff %s %s\nreturned error\nDiff output:\n%s", fileName, outputFileName, output)
		t.Error("Error comparing tar files")
	}

	if len(output) > 0 {
		t.Error("Checked out tar file dose not match input")
	}

	os.RemoveAll(workDirFolder)
	os.RemoveAll(repoPath)
}
