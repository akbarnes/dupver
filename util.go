package dupver

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/bmatcuk/doublestar"
)

const HexChars = "0123456789abcdef"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func ExcludedFile(fileName string, info os.FileInfo, filters []string) bool {
	// goverDir := filepath.Join(WorkingDirectory, ".gover2")
	goverDir := ".gover2"
	goverPattern := filepath.Join(goverDir, "**")

	if info.IsDir() {
		return true
	}

	matched, err := doublestar.PathMatch(goverPattern, fileName)

	if err != nil && VerboseMode {
		fmt.Printf("Error matching %s\n", goverDir)
	}

	if matched {
		if VerboseMode {
			fmt.Printf("Skipping file %s in .gover2\n", fileName)
		}

		return true
	}

	for _, pattern := range filters {
		matched, err := doublestar.PathMatch(pattern, fileName)

		if err != nil && VerboseMode {
			fmt.Printf("Error matching %s\n", goverDir)
		}

		if matched {
			if VerboseMode {
				fmt.Printf("Skipping file %s which matches with %s\n", fileName, pattern)
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
