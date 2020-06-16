package main

import (
	"testing"
	"os"
	"path"
)

func TestInit(t *testing.T) {
	homeDir := GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	workDirID := RandString(16, hexChars)
	workDirFolder := "Test_" + workDirID
	// workDirPath := path.Join("temp", workDirFolder)
	err := os.MkdirAll(workDirFolder, 777)

	if err != nil {
		t.Error("Could not create workdir")
	}

	repoPath := path.Join(homeDir, ".dupver_repo")

	workDirName := ""
	InitWorkDir(workDirFolder, workDirName, repoPath)

	cfg := ReadWorkDirConfig(workDirFolder)

	if cfg.RepoPath != repoPath {
		t.Error("Incorrect repo path retrieved")
	}

	expectedWorkDirName := "test_" + workDirID
	if cfg.WorkDirName != expectedWorkDirName {
		t.Error("Incorrect workdir name retrieved")
	}	
}