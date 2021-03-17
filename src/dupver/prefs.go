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
func ReadPrefs(opts Options) (Preferences, error) {
	return ReadPrefsFile(GetPrefsPath(), opts)
}

// Load global preferences given a preferences file path
func ReadPrefsFile(filePath string, opts Options) (Preferences, error) {
	var prefs Preferences
	// TODO: set this differently for linux
	prefs.DiffTool = "bcompare"
	prefs.DefaultRepo = filepath.Join(GetHome(), ".dupver_repo")

	f, err := os.Open(filePath)

	if err != nil {
		fancyprint.Warn("Preferences file missing, creating default")
		SavePrefsFile(filePath, prefs, false, opts)
		return prefs, errors.New("Preferences file missing")
	}

	if _, err = toml.DecodeReader(f, &prefs); err != nil {
		panic(fmt.Sprintf("Invalid preferences file %s\n", filePath))
	}

	f.Close()

	return prefs, nil
}

// Save global preferences
func SavePrefs(prefs Preferences, forceWrite bool, opts Options) {
	prefsPath := filepath.Join(GetHome(), ".dupver", "prefs.toml")
	SavePrefsFile(prefsPath, prefs, forceWrite, opts)
}

// Save global preferences given a preferences file
func SavePrefsFile(prefsPath string, prefs Preferences, forceWrite bool, opts Options) {
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

// Print the current global preferences
func PrintCurrentPreferences(opts Options) {
	prefs, err := ReadPrefs(opts)

	if err != nil {
		// Todo: handle invalid configuration file
		fancyprint.Warn("Could not read preferences file.")
		os.Exit(1)
	}

	PrintPreferences(prefs)
}

// Print the global preferences structure
func PrintPreferences(prefs Preferences) {
	fmt.Printf("Editor: %s\n", prefs.Editor)
	fmt.Printf("Diff tool: %s\n", prefs.DiffTool)
	fmt.Printf("Default repository path: %s\n", prefs.DefaultRepo)
}

// Print the global preferences structure
func (prefs Preferences) Print() {
	fmt.Printf("Editor: %s\n", prefs.Editor)
	fmt.Printf("Diff tool: %s\n", prefs.DiffTool)
	fmt.Printf("Default repository path: %s\n", prefs.DefaultRepo)
}

// Print the current global preferences as JSON
func PrintCurrentPreferencesAsJson(opts Options) {
	prefs, err := ReadPrefs(opts)

	if err != nil {
		// Todo: handle invalid configuration file
		fancyprint.Warn("Could not read preferences file.")
		os.Exit(1)
	}

	PrintJson(prefs)
}

// Print the current global preferences as JSON
func (prefs Preferences) PrintJson() {
	PrintJson(prefs)
}