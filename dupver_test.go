package main

import (
	"testing"
	"os"
	"path"
	"fmt"
	"os/exec"
	"log"
	"./dupver"
)

func TestWorkRepoInit(t *testing.T) {
	homeDir := dupver.GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	repoId := dupver.RandString(16, dupver.HexChars)
	repoFolder := ".dupver_repo_" + repoId

	repoPath := path.Join(homeDir, "temp", repoFolder)
	dupver.InitRepo(repoPath)

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

	cfg := dupver.ReadRepoConfigFile(path.Join(repoPath, "config.toml"))

	if cfg.Version != 2 {
		t.Error("Invalid repository version", cfg.Version)
	}

	if cfg.ChunkerPolynomial <= 0 {
		t.Error("Invalid chunker polynomial", cfg.ChunkerPolynomial)
	}

	os.RemoveAll(repoPath)
}

func TestWorkDirInit(t *testing.T) {
	homeDir := dupver.GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	workDirId := dupver.RandString(16, dupver.HexChars)
	workDirFolder := "Test_" + workDirId
	// workDirPath := path.Join("temp", workDirFolder)
	err := os.MkdirAll(workDirFolder, 0777)	

	if err != nil {
		t.Error("Could not create workdir")
	}

	repoPath := path.Join(homeDir, ".dupver_repo")

	workDirName := ""
	dupver.InitWorkDir(workDirFolder, workDirName, repoPath)

	cfg :=dupver.ReadWorkDirConfig(workDirFolder)

	if cfg.RepoPath != repoPath {
		t.Error("Incorrect repo path retrieved")
	}

	expectedWorkDirName := "test_" + workDirId
	if cfg.WorkDirName != expectedWorkDirName {
		t.Error("Incorrect workdir name retrieved")
	}	

	os.RemoveAll(workDirFolder)
}


func TestCommit(t *testing.T) {
	// homeDir := GetHome()
	verbosity := 1
	msg := "Commit random data"

    // ----------- Create a repo ----------- //    
	homeDir := dupver.GetHome()
	repoId := dupver.RandString(16, dupver.HexChars)
	repoFolder := ".dupver_repo_" + repoId
	repoPath := path.Join(homeDir, "temp", repoFolder)
	dupver.InitRepo(repoPath)	

    // ----------- Create a workdir ----------- //    
	workDirId := dupver.RandString(16, dupver.HexChars)
	workDirFolder := "Test_" + workDirId
	err := os.MkdirAll(workDirFolder, 0777)	

	if err != nil {
		t.Error("Could not cerate workdir folder " + workDirFolder)
	}

	workDirName := ""
	dupver.InitWorkDir(workDirFolder, workDirName, repoPath)

	// ----------- Create tar file with random data ----------- //    
	// TODO: add random permutes to data
	fileName := dupver.CreateRandomTarFile(workDirFolder, repoPath)
	fmt.Printf("Created tar file %s\n", fileName)

	// ----------- Commit the tar file  ----------- //    
	snapshot := dupver.CommitFile(fileName, msg, verbosity)

	// ----------- Commit the tar file  ----------- //   
	myWorkDirConfig := dupver.ReadWorkDirConfig(workDirFolder)
	dupver.PrintSnapshots(dupver.ListSnapshots(myWorkDirConfig), "")
	dupver.PrintSnapshots(dupver.ListSnapshots(myWorkDirConfig), snapshot)


	// ----------- Checkout the tar file  ----------- //    
	mySnapshot := dupver.ReadSnapshot(snapshot, myWorkDirConfig)
	timeStr := dupver.TimeToPath(mySnapshot.Time)
	outputFileName := fmt.Sprintf("%s-%s-%s.tar", myWorkDirConfig.WorkDirName, timeStr, snapshot[0:16])

	dupver.UnpackFile(outputFileName, myWorkDirConfig.RepoPath, mySnapshot.ChunkIDs, verbosity) 
	fmt.Printf("Wrote to %s\n", outputFileName)

	cmd := exec.Command("diff", fileName, outputFileName)
	log.Printf("Running command and waiting for it to finish...")
	output, err := cmd.Output()

	if err != nil {
		t.Error("Error comparing tar files")
	}

	if len(output) > 0 {
		t.Error("Checked out tar file dose not match input")
	}

	os.RemoveAll(workDirFolder)
	os.RemoveAll(repoPath)
}