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

	workDirFolder := "Test_" + RandString(40, hexchars)
	workDirPath := path.Join(homeDir, "temp", workDirFolder)
	os.MkdirAll(workDirPath, 0777)
	repoPath := path.Join(homeDir, ".dupver_repo")

	workDirName := ""
	InitWorkDir(workDirFolder, workDirName, repoPath)

	cfg := ReadWorkdirConfig(workdirPath)

	if cfg.RepoPath != repoPath {
		t.Error("Incorrect repo path retrieved")
	}

	expectedWorkdirName := "test_" + RandString(40, hexchars)
	if cfg.workDirName != expectedWorkdirName {
		t.Error("Incorrect workdir name retrieved")
	}	
}