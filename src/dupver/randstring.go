package dupver

import (
  "math/rand"
  "time"
)


var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

const HexChars = "0123456789abcdef"

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


