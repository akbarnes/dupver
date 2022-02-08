package dupver

import (
	"fmt"
	"os"
    "os/exec"
	"path/filepath"
)

func DiffToolSnapshot(diffTool string, archiveTool string) {
    snap := ReadHead()
	fmt.Fprintf(os.Stderr, "Comparing %s\n", snap.SnapshotID[0:9])
    snap.DiffTool(diffTool, archiveTool)
}

func (snap Snapshot) DiffTool(diffTool string, archiveTool string) {
    home, err := os.UserHomeDir()
    Check(err)
    tempFolder := filepath.Join(home, ".dupver", "temp", RandHexString(24))
    snap.Checkout(tempFolder, "*", archiveTool)
    cmd := exec.Command(diffTool, tempFolder, ".")
    cmd.Run()
}

func DiffToolSnapshotFile(fileName string, diffTool string, archiveTool string) {
    snap := ReadHead()
	fmt.Fprintf(os.Stderr, "Comparing %s/%s\n", snap.SnapshotID[0:9], fileName)
    snap.DiffToolFile(fileName, diffTool, archiveTool)
}

func (snap Snapshot) DiffToolFile(fileName string, diffTool string, archiveTool string) {
    home, err := os.UserHomeDir()
    Check(err)
    tempFolder := filepath.Join(home, ".dupver", "temp", RandHexString(24))
    snap.Checkout(tempFolder, fileName, archiveTool)
    tempFile := filepath.Join(tempFolder, fileName)
    cmd := exec.Command(diffTool, tempFile, fileName) 
    cmd.Run()
}

