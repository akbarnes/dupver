package dupver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type PrefsMissingError struct {
	err string
}

func (e *PrefsMissingError) Error() string {
	return "Prefs are missing"
}

// TODO: change this to SerializedSnaphot
// and use Time type for SnapshotTime?
type Prefs struct {
	PrefsMajorVersion   int64
	PrefsMinorVersion   int64
	DupverMajorVersion int64
	DupverMinorVersion int64
    Editor         string
	DiffTool       string
}

func CreateDefaultPrefs() Prefs {
    p := Prefs{Editor: "vi", DiffTool: "kdiff3"}
	p.DupverMajorVersion = MajorVersion
	p.DupverMinorVersion = MinorVersion
	p.PrefsMajorVersion = PrefsMajorVersion
	p.PrefsMinorVersion = PrefsMinorVersion

    if IsWindows() {
        p.Editor = "notepad"
    }

    return p
}

func (p Prefs) Write() {
    home, err := os.UserHomeDir()

    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to get user home directory, not writing default prefs\n")
        return
    }

	prefsPath := filepath.Join(home, ".dupver.json")

	f, err := os.Create(prefsPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create prefs json %s", prefsPath))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(p)
	f.Close()
}

func (p Prefs) CorrectPrefsVersion() bool {
	return p.PrefsMajorVersion == PrefsMajorVersion
}

func AbortIfIncorrectPrefsVersion() {
	p, err := ReadPrefs(false)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read prefs, exiting\n")
		os.Exit(1)
	}

	if !p.CorrectPrefsVersion() {
		fmt.Fprintf(os.Stderr, "Incorrect prefs version of %d.%d, expecting %d.x\n", p.PrefsMajorVersion, p.PrefsMinorVersion, PrefsMajorVersion)
		os.Exit(1)
	}
}

func ReadPrefs(writeIfMissing bool) (Prefs, error) {
    home, err := os.UserHomeDir()

    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to get user home directory, returning default prefs\n")
        return CreateDefaultPrefs(), nil
    }

	prefsPath := filepath.Join(home, ".dupver.json")

	if VerboseMode {
		fmt.Fprintf(os.Stderr, "Reading %s\n", prefsPath)
	}

	var p Prefs
	f, err := os.Open(prefsPath)
	defer f.Close()

	if errors.Is(err, os.ErrNotExist) {
		if writeIfMissing {
			if VerboseMode {
				fmt.Fprintf(os.Stderr, "Prefs not present, writing default")
			}

			p = CreateDefaultPrefs()
			p.Write()
			return p, nil
		} else {
			return Prefs{}, err
		}
	} else if err != nil {
		return Prefs{}, errors.New("Cannot open prefs")
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&p); err != nil {
		panic("Cannot decode prefs")
	}

	if !p.CorrectPrefsVersion() {
		panic(fmt.Sprintf("Invalid prefs version %d.%d, expecting %d.x\n", p.PrefsMajorVersion, p.PrefsMinorVersion, PrefsMajorVersion))
	}

	return p, nil
}
