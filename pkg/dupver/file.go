package dupver

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const NumChars = 40

func HashFile(FileName string, NumChars int) (string, error) {
	var data []byte
	var err error

	data, err = ioutil.ReadFile(FileName)

	if err != nil {
		return "", fmt.Errorf("Error hashing file, unable to read content: %w", err)
	}

	sum := fmt.Sprintf("%x", sha256.Sum256(data))

	if len(sum) < NumChars || NumChars < 0 {
		NumChars = len(sum)
	}

	return sum[0:NumChars], nil
}

// Copy the source file to a destination file. Any existing file
// will be overwritten and will not copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Error copying file, cannot open source file: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("Error copying file, cannot open destination file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("Error copying file: %w", err)
	}
	return out.Close()
}
