package main

import (
	"testing"
	"os"
	"path"
	"fmt"
)

func TestWorkRepoInit(t *testing.T) {
	homeDir := GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	repoId := RandString(16, hexChars)
	repoFolder := ".dupver_repo_" + repoId

	repoPath := path.Join(homeDir, "temp", repoFolder)
	InitRepo(repoPath)

	snapshotsPath := path.Join(repoPath, "snapshots")
	if _, err := os.Stat(snapshotsPath); err != nil {
		// path/to/whatever exists
		t.Error("Did not create snapshots folder", snapshotsPath)
	} 

	treesPath := path.Join(repoPath, "trees")
	if _, err := os.Stat(treesPath); err != nil {
		// path/to/whatever exists
		t.Error("Did not create trees folder", treesPath)
	} 	

	cfg := ReadRepoConfigFile(path.Join(repoPath, "config.toml"))

	if cfg.Version != 1 {
		t.Error("Invalid repository version", cfg.Version)
	}

	if cfg.ChunkerPolynomial <= 0 {
		t.Error("Invalid chunker polynomial", cfg.ChunkerPolynomial)
	}
}

func TestWorkDirInit(t *testing.T) {
	homeDir := GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	workDirId := RandString(16, hexChars)
	workDirFolder := "Test_" + workDirId
	// workDirPath := path.Join("temp", workDirFolder)
	err := os.MkdirAll(workDirFolder, 777)

	if err != nil {
		t.Error("Could not create workdir")
	}

	repoPath := path.Join(homeDir, ".dupver_repo")

	workDirName := ""
	InitWorkDir(workDirFolder, workDirName, repoPath)

	return

	cfg := ReadWorkDirConfig(workDirFolder)

	if cfg.RepoPath != repoPath {
		t.Error("Incorrect repo path retrieved")
	}

	expectedWorkDirName := "test_" + workDirId
	if cfg.WorkDirName != expectedWorkDirName {
		t.Error("Incorrect workdir name retrieved")
	}	
}


func TestCommit(t *testing.T) {
	// homeDir := GetHome()
	// verbosity := 1
	// msg := "Commit random data"

	fileName := CreateRandomTarFile()
	fmt.Printf("Created tar file %s", fileName)

	// CommitFile(filePath, msg, verbosity)


	// expectedWorkDirName := "test_" + workDirId
	// if cfg.WorkDirName != expectedWorkDirName {
	// 	t.Error("Incorrect workdir name retrieved")
	// }	
}