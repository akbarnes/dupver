package main

import (
    "os"
    "fmt"
    "strings"
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
