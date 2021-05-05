// Package implements working directory, repository and preferences handling
// for the Dupver application
package dupver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/akbarnes/dupver/src/fancyprint"
)

type Preferences struct {
	DiffTool    string
	Editor      string
	DefaultRepo string
}

type Options struct {
	WorkDirName  string
	RepoName     string
	RepoPath     string
	Branch       string
	DestRepoName string
	DestRepoPath string
	JsonOutput   bool
}

func GetPrefsPath() string {
	return filepath.Join(GetHome(), ".dupver", "prefs.toml")
}

// Load global preferences
func ReadPrefs() (Preferences, error) {
	return ReadPrefsFile(GetPrefsPath())
}

// Load global preferences given a preferences file path
func ReadPrefsFile(filePath string) (Preferences, error) {
	var prefs Preferences
	// TODO: set this differently for linux
	prefs.DiffTool = "bcompare"
	prefs.DefaultRepo = filepath.Join(GetHome(), ".dupver_repo")

	f, err := os.Open(filePath)

	if err != nil {
		fancyprint.Warn("Preferences file missing, creating default")
		prefs.SaveFile(filePath, false)
		return prefs, errors.New("Preferences file missing")
	}

	if _, err = toml.DecodeReader(f, &prefs); err != nil {
		panic(fmt.Sprintf("Invalid preferences file %s\n", filePath))
	}

	f.Close()

	return prefs, nil
}

// Save global preferences
func (prefs Preferences) Save(forceWrite bool) {
	prefsPath := filepath.Join(GetHome(), ".dupver", "prefs.toml")
	prefs.SaveFile(prefsPath, forceWrite)
}

// Save global preferences given a preferences file
func (prefs Preferences) SaveFile(prefsPath string, forceWrite bool) {
	if _, err := os.Stat(prefsPath); err == nil && !forceWrite {
		// panic("Refusing to write existing project workdir config " + configPath)
		panic(fmt.Sprintf("Refusing to write existing preferences %s\n", prefsPath))
	}

	fancyprint.Infof("Writing prefs:\n%+v\n", prefs)
	fancyprint.Infof("to: %s\n", prefsPath)

	CreateSubFolder(GetHome(), ".dupver")

	f, _ := os.Create(prefsPath)
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(prefs)
	f.Close()
}

// Print the global preferences structure
func (prefs Preferences) Print() {
	fmt.Printf("Editor: %s\n", prefs.Editor)
	fmt.Printf("Diff tool: %s\n", prefs.DiffTool)
	fmt.Printf("Default repository path: %s\n", prefs.DefaultRepo)
}

// Print the current global preferences as JSON
func (prefs Preferences) PrintJson() {
	PrintJson(prefs)
}