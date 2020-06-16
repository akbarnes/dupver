package main

import (
    "os"
    "io"
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

func WriteRandom(f *os.File, numLines int, numCols int) {
    for r := 0; r < numLines; r += 1 {
        s := RandString(numCols, hexChars)
        fmt.Fprintf(f, "%s\n", s)
    }
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
    f.Close()
}

func WriteRandomTar(buf *os.File) {
    nFiles := 10

    tw := tar.NewWriter(&buf)
    var files = []string

    for i := 0; i < nFiles; i += 1 {
        files = append(files, RandString(24, hexChars) + ".txt")
    }
        
    for _, file := range files {
        bytes := RandString(50*1000*1000, hexChars)

        hdr := &tar.Header{
            Name: file.Name,
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