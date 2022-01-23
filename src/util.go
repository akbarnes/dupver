package dupver

import (
	"fmt"
	"math/rand"
	"os"
    "strings"
	"path/filepath"
	"time"
    "bufio"

	"github.com/bmatcuk/doublestar"
)

const HexChars = "0123456789abcdef"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func IsWindows() bool {
    return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func ReadFilters() ([]string, error) {
	filterPath := ".dupver_ignore"
	var filters []string
	f, err := os.Open(filterPath)

	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		} else {
			err = fmt.Errorf("Ignore file %s exists but encountered error trying to open it: %w", filterPath, err)
			return []string{}, err
		}
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
        line := strings.ReplaceAll(scanner.Text(), "\n", "")

        if len(line) == 0 {
            continue
        }

		filters = append(filters, line)
	}

	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("Encountered an error while attemping to read filters from %s: %w", filterPath, err)
		return []string{}, err
	}

	return filters, nil
}

func ExcludedFile(fileName string, info os.FileInfo, filters []string) bool {
	// dupverDir := filepath.Join(WorkingDirectory, ".gover2")
	dupverDir := ".dupver"
	dupverPattern := filepath.Join(dupverDir, "**")

	if info.IsDir() {
		return true
	}

	matched, err := doublestar.PathMatch(dupverPattern, fileName)

	if err != nil && VerboseMode {
		fmt.Fprintf(os.Stderr, "Error matching %s\n", dupverDir)
	}

	if matched {
		if DebugMode {
			fmt.Fprintf(os.Stderr, "Skipping file %s in .dupver\n", fileName)
		}

		return true
	}

	for _, pattern := range filters {
		matched, err := doublestar.PathMatch(pattern, fileName)

		if err != nil && VerboseMode {
			fmt.Fprintf(os.Stderr, "Error matching %s\n", dupverDir)
		}

		if matched {
			if VerboseMode {
				fmt.Fprintf(os.Stderr, "Skipping file %s which matches with %s\n", fileName, pattern)
			}

			return true
		}
	}

	return false
}

// Return a random string of specified length with hexadecimal characters
func RandHexString(length int) string {
	return RandString(length, HexChars)
}

// Return a random string of specified length with an arbitrary character set
func RandString(length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
