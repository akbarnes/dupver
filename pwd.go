package main


import (
	"os"
	"path"
	"fmt"
	"log"
)


func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// _, folder := path.Split(dir)
	folder := path.Base(dir)	

	fmt.Printf("dir: %s, folder: %s\n", dir, folder)
}