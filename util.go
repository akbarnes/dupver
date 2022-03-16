package dupver

import (
	"fmt"
	"math/rand"
	"os"
    "strings"
	"path/filepath"
	"time"
    "bufio"
    "log"
    "io"
    "os/exec"

	"github.com/bmatcuk/doublestar"
)

const HexChars = "0123456789abcdef"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func IsWindows() bool {
    return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func ToForwardSlashes(path string) string {
    return strings.ReplaceAll(path, "\\", "/")
}

func ToNativeSeparators(path string) string {
    return strings.ReplaceAll(path, "/", string(os.PathSeparator))
}


func ReadFilters() ([]string, error) {
	filterPath := ".dupver_ignore"
	f, err := os.Open(filterPath)

	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		} else {
			err = fmt.Errorf("Ignore file %s exists but encountered error trying to open it: %w", filterPath, err)
			return []string{}, err
		}
	}

    return ReadFilterFile(f)
}

func ReadArchiveTypes() ([]string, error) {
	filterPath := ".dupver_archive_types"
	f, err := os.Open(filterPath)

	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		} else {
			err = fmt.Errorf("Archive type file %s exists but encountered error trying to open it: %w", filterPath, err)
			return []string{}, err
		}
	}

    return ReadFilterFile(f)
}

func ReadFilterFile(f *os.File) ([]string, error) {
	scanner := bufio.NewScanner(f)
	var filters []string

	for scanner.Scan() {
        line := strings.TrimSuffix(scanner.Text(), "\n")

        if len(line) == 0 {
            continue
        }

		filters = append(filters, ToForwardSlashes(line))
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return filters, nil
}


func ExcludedFile(fileName string, info os.FileInfo, filters []string) bool {
	// dupverDir := filepath.Join(WorkingDirectory, ".gover2")
	dupverDir := ".dupver"
	dupverPattern := ToForwardSlashes(filepath.Join(dupverDir, "**"))

	if info.IsDir() {
		return true
	}

	matched, err := doublestar.Match(dupverPattern, fileName)

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
		matched, err := doublestar.Match(pattern, fileName)

        //fmt.Printf("file: %s\npattern: %s\n\n", fileName, pattern)
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

func ArchiveFile(fileName string, info os.FileInfo, archiveTypes []string) bool {
	if info.IsDir() {
		return false
	}

	for _, ext := range archiveTypes {
		if strings.HasSuffix(fileName, ext) {
			if VerboseMode {
				fmt.Fprintf(os.Stderr, "Preprocessing archive file %s which matches with type %s\n", fileName, ext)
			}

			return true
		}
	}

	return false
}

func GenArchiveBaseName() string {
    return RandHexString(24)
}

func GenTempArchivePath(archiveBaseName string) (string, error) {
    home, err := os.UserHomeDir()

    if err != nil {
        return "", err
    }

   return filepath.Join(home, ".dupver", "temp", archiveBaseName+".zip"), nil
}

// Currently only 7-zip is supported
func PreprocessArchive(fileName string, archiveTool string) (string, error) {
    home, err := os.UserHomeDir()

    if err != nil {
        return "", err
    }

    // Note that 7z will create folder structure as needed
    archiveBaseName := GenArchiveBaseName()
    extractFolder := filepath.Join(home, ".dupver", "temp", archiveBaseName)
    extractCmd := exec.Command(archiveTool, "x", "-o"+extractFolder, fileName)

    // TODO: add error wrapping
    if err = extractCmd.Run(); err != nil {
        return "", err
    }

    extractGlob := filepath.Join(extractFolder, "*")
    archiveFile := filepath.Join(home, ".dupver", "temp", archiveBaseName+".zip")
    compressCmd := exec.Command(archiveTool, "a", "-mm=Copy", archiveFile, extractGlob)

    if err = compressCmd.Run(); err != nil {
        return archiveFile, err
    }

    return archiveFile, nil
}

func PostprocessArchive(archiveBaseName string, outputFile string, archiveTool string) error {
    home, err := os.UserHomeDir()

    if err != nil {
        return err
    }

    // Note that 7z will create folder structure as needed
    archiveFile := filepath.Join(home, ".dupver", "temp", archiveBaseName+".zip")
    extractFolder := filepath.Join(home, ".dupver", "temp", archiveBaseName)
    extractCmd := exec.Command(archiveTool, "x", "-o"+extractFolder, archiveFile)

    // TODO: add error wrapping
//    if err = extractCmd.Run(); err != nil {
//        return fmt.Errorf("Error extracting %s to %s: %w", archiveFile, extractFolder, err)
//    }
    stderr, err := extractCmd.StderrPipe()

    if err := extractCmd.Start(); err != nil {

        log.Fatal(err)
    }

    slurp, _ := io.ReadAll(stderr)

    if err := extractCmd.Wait(); err != nil {
        log.Fatal(err)
    }

    extractGlob := filepath.Join(extractFolder, "*")
    compressCmd := exec.Command(archiveTool, "a", outputFile, extractGlob)

    if err = compressCmd.Run(); err != nil {
        return fmt.Errorf("Error compressing %s to %s: %w\n%s", extractFolder, outputFile, err, slurp)
    }

    return nil
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
