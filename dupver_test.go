package main

import (
	"testing"
	"os"
	"path"
)

func TestInit(t *testing.T) {
	homeDir := GetHome()

	if len(homeDir) == 0 {
		t.Error("Test failed")
	}

	workDirFolder := "test_" + RandString(40, hexchars)
	workDirPath := path.Join(homeDir, "temp", workDirFolder)
	os.MkdirAll(workDirPath, 0777)
	repoPath := path.Join(homeDir, ".dupver_repo")

	workDirName := ""
	InitWorkDir(workDirFolder, workDirName, repoPath)
}