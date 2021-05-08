package dupver

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/akbarnes/dupver/src/fancyprint"
)

func TestWorkRepoInit(t *testing.T) {
	debug := false
	verbose := true
	quiet := false
	monochrome := false
	fancyprint.Setup(debug, verbose, quiet, monochrome)

	homeDir := GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	repoId := RandString(16, HexChars)
	repoFolder := ".dupver_repo_" + repoId
	repoName := "test"

	repoPath := filepath.Join(homeDir, "temp", repoFolder)
	InitRepo(repoPath, repoName, "", zip.Deflate, false)

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

	cfg, err := ReadRepoConfigFile(filepath.Join(repoPath, "config.toml"))

	if err != nil {
		t.Error("could not read repo config file for repo", repoPath)
	}

	if cfg.Version != 2 {
		t.Error("Invalid repository version", cfg.Version)
	}

	if cfg.ChunkerPolynomial <= 0 {
		t.Error("Invalid chunker polynomial", cfg.ChunkerPolynomial)
	}

	os.RemoveAll(repoPath)
}

func TestWorkDirInit(t *testing.T) {
	opts := Options{}
	homeDir := GetHome()

	if len(homeDir) == 0 {
		t.Error("Could not read home directory environment variable")
	}

	workDirId := RandString(16, HexChars)
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

	InitWorkDir(workDirFolder, projectName, opts)

	cfg, _ := ReadWorkDirConfig(workDirFolder)

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
	opts := Options{}
	msg := "Commit random data"

	// ----------- Create a repo ----------- //
	homeDir := GetHome()
	repoId := RandString(16, HexChars)
	repoFolder := ".dupver_repo_" + repoId
	repoPath := filepath.Join(homeDir, "temp", repoFolder)
	repoName := "test"
	InitRepo(repoPath, repoName, "", zip.Deflate, false)

	// ----------- Create a workdir ----------- //
	workDirId := RandString(16, HexChars)
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

	InitWorkDir(workDirFolder, projectName, opts)

	// ----------- Create tar file with random data ----------- //
	// TODO: add random permutes to data
	fileName := CreateRandomTarFile(workDirFolder, repoPath)
	fmt.Printf("Created tar file %s\n", fileName)

	// ----------- Commit the tar file  ----------- //
	workDir, _ :=  LoadWorkDir(workDirFolder)
	snapshot := workDir.CommitFile(fileName, nil, msg, true)

	// ----------- Commit the tar file  ----------- //
	// TODO: Replace with PrintSnapshots
	workDir.PrintSnapshots()

	fmt.Printf("snapshot: %+v\n\n", snapshot)
	fmt.Printf("workdir: %+v\n\n", workDir)

	
	// ----------- Checkout the tar file  ----------- //
	mySnapshot := workDir.ReadSnapshot(snapshot.ID)
	timeStr := TimeToPath(mySnapshot.Time)
	outputFileName := fmt.Sprintf("%s-%s-%s.tar", workDir.ProjectName, timeStr, snapshot.ID[0:16])
	UnpackFile(outputFileName, opts.RepoPath, mySnapshot.ChunkIDs)
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


func TestCopy(t *testing.T) {
	opts := Options{}
	msg := "Commit random data"

	// ----------- Create a repo ----------- //
	homeDir := GetHome()
	repoId := RandString(16, HexChars)
	repoFolder := ".dupver_repo_" + repoId
	repoPath := filepath.Join(homeDir, "temp", repoFolder)
	repoName := "test"
	InitRepo(repoPath, repoName, "", zip.Deflate, false)

	repoId2 := RandString(16, HexChars)
	repoFolder2 := ".dupver_repo_" + repoId2
	repoPath2 := filepath.Join(homeDir, "temp", repoFolder2)
	repoName2 := "test2"
	InitRepo(repoPath2, repoName2, "", zip.Deflate, false)

	// ----------- Create a workdir ----------- //
	workDirId := RandString(16, HexChars)
	workDirFolder := "Test_" + workDirId
	err := os.MkdirAll(workDirFolder, 0777)

	if err != nil {
		t.Error("Could not create workdir folder " + workDirFolder)
	}

	projectName := "test"

	opts.WorkDirName = "test"
	opts.RepoName = "test"
	opts.RepoPath = repoPath
	opts.Branch = "main"

	InitWorkDir(workDirFolder, projectName, opts)

	// ----------- Create tar file with random data ----------- //
	// TODO: add random permutes to data
	fileName := CreateRandomTarFile(workDirFolder, repoPath)
	fmt.Printf("Created tar file %s\n", fileName)

	// ----------- Commit the tar file  ----------- //
	workDir, _ :=  LoadWorkDir(workDirFolder)
	snapshot := workDir.CommitFile(fileName, nil, msg, true)	

	// ----------- Copy to the second repo  ----------- //
	snapshotId := snapshot.ID
	CopySnapshot(snapshotId, repoPath, repoPath2, opts)

	opts2 := Options{}
	opts2.WorkDirName = opts.WorkDirName
	opts2.RepoName = opts.RepoName
	opts2.RepoPath = repoPath2
	opts2.Branch = opts.Branch

	// ----------- Commit the tar file  ----------- //
	workDir.Repo.Path = repoPath2
	workDir.PrintSnapshots()

	fmt.Printf("snapshot: %+v\n\n", snapshot)
	fmt.Printf("workdir: %+v\n\n", workDir)

	// // ----------- Checkout the tar file  ----------- //
	mySnapshot := ReadSnapshot(snapshot.ID, opts2)
	timeStr := TimeToPath(mySnapshot.Time)
	outputFileName := fmt.Sprintf("%s-%s-%s.tar", workDir.ProjectName, timeStr, snapshot.ID[0:16])

	// UnpackFile(outputFileName, myWorkDirConfig.RepoPath, mySnapshot.ChunkIDs, verbosity)
	UnpackFile(outputFileName, opts2.RepoPath, mySnapshot.ChunkIDs)
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
