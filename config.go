package main

import (
	"os"
	"github.com/BurntSushi/toml"
)

type workDirConfig struct {
	RepositoryPath string
}


func SaveWorkDirConfig(myWorkDirConfig workDirConfig) {
	f, _ := os.Create(".dupver/config.toml")
	WriteWorkDirConfig(f, myWorkDirConfig)
}


func WriteWorkDirConfig(f *os.File, myWorkDirConfig workDirConfig) {
	myEncoder := toml.NewEncoder(f)
	myEncoder.Encode(myWorkDirConfig)
}


func ReadConfigFile(filePath string) (workDirConfig) {
	var myConfig workDirConfig
	tarFile, _ := os.Open(filePath)


	// Open and iterate through the files in the archive.
	tr := tar.NewReader(tarFile)
	i := 0
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}

		i++
		fmt.Printf("File %d: %s\n", i, hdr.Name)
	}

	fmt.Fprint(commitFile, "]\n")
	tarFile.Close()

	myConfig.RepositoryPath = "C:\\Users\\305232\\DupVerRepo"
	return myConfig
}


// func ReadConfig(tarFile *os.File) {
// }