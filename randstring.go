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

