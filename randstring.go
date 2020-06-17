package main

import (
  "math/rand"
  "time"
)


var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

const hexChars = "0123456789abcdef"

func RandHexString(length int) string {
	return RandString(length, hexChars)
}

func RandString(length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
	  b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func RandText(lines int, cols int, charset string) string {
	colChars := cols + 1

	b := make([]byte, lines*colChars)

    for r := 0; r < lines; r += 1 {
		for c := 0; c < cols; c += 1 {
			b[r*colChars + c] = charset[seededRand.Intn(len(charset))]
		}

		b[colChars - 1] = '\n'
    }

	return string(b)
}

