package main

import (
    "os"
    // "io"
    "bufio"
    "fmt"
    "log"
    "strings"
    "archive/tar"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func GetHome() string {
    for _, e := range os.Environ() {
        pair := strings.SplitN(e, "=", 2)
        // fmt.Println(pair[0])

        if pair[0] == "HOME" || pair[0] == "USERPROFILE" {
            return pair[1]
        } 
    }

    fmt.Println("Warning! No home variable defined")
    return ""
}


func CreateRandomTarFile() string {
    fileName := "test_" + RandString(24, hexChars) + ".tar"
    WriteRandomTarFile(fileName)
    return fileName
}


func WriteRandomTarFile(fileName string) {
    f, err := os.Create(fileName)
    check(err)
    WriteRandomTar(f)
    // WriteRandom(f, 50000, 100, hexChars)
    f.Close()
}


func WriteRandomText(f *os.File, numLines int, numCols int, charset string) {
	w := bufio.NewWriter(f)

	b := make([]byte, numCols)

    for r := 0; r < numLines; r += 1 {
		for c := range b {
			b[c] = charset[seededRand.Intn(len(charset))]
		}
			  
		fmt.Fprintf(w, "%s\n", b)
    }

    w.Flush()
}


func WriteRandomTar(buf *os.File) {
    nFiles := 10

    tw := tar.NewWriter(buf)
    var files = []string{}

    for i := 0; i < nFiles; i += 1 {
        files = append(files, RandString(24, hexChars) + ".txt")
    }
        
    for _, file := range files {
        n := 50*1000*1000
        bytes := make([]byte, n)
        seededRand.Read(bytes)

        hdr := &tar.Header{
            Name: file,
            Mode: 0600,
            Size: int64(len(bytes)),
        }
        if err := tw.WriteHeader(hdr); err != nil {
            log.Fatal(err)
        }
        if _, err := tw.Write(bytes); err != nil {
            log.Fatal(err)
        }
    }
    if err := tw.Close(); err != nil {
        log.Fatal(err)
    }
}