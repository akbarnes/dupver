package main


import (
	"flag"
    "fmt"
	"io"
	"crypto/sha256"
	"os"
	"github.com/restic/chunker"
	"compress/gzip"
	"archive/zip"
)
 	

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Unzip(src string, dest string) ([]string, error) {

    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return filenames, fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return filenames, err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return filenames, err
        }

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            return filenames, err
        }
    }
    return filenames, nil
}

func main() {
	filePtr := flag.String("file", "ACTIVSg70k.RAW", "an int")
	backupPtr := flag.Bool("backup", false, "Back up specified file")
	// restorePtr := flag.Bool("restore", false, "Restore specified file")
	msgPtr := flag.String("message", "", "commit message")
	
	flag.Parse()
	
	filePath := *filePtr
	msg := *msgPtr
	// filePath = "ACTIVSg70k.RAW"	

	if (*backupPtr == true) {
		fmt.Println("Backing up ", filePath)


		var filenames []string

		r, err := zip.OpenReader(filePath)
		if err != nil {
			return filenames, err
		}
		defer r.Close()
	
		for _, f := range r.File {
	
			// Store filename/path for returning and using later on
			fpath := filepath.Join(dest, f.Name)
	
			// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
			if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
				return filenames, fmt.Errorf("%s: illegal file path", fpath)
			}
	
			filenames = append(filenames, fpath)
	
			if f.FileInfo().IsDir() {
				// Make Folder
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}
	
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}
	
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}
	
			rc, err := f.Open()
			if err != nil {
				return filenames, err
			}
	
			_, err = io.Copy(outFile, rc)
	
			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			rc.Close()
	
			if err != nil {
				return filenames, err
			}
		}
		return filenames, nil
	

		// chunky =  Chunker(rd io.Reader, pol Pol) 
		f, _ := os.Open(filePath)
		
		// generate 32MiB of deterministic pseudo-random data
		// data := getRandom(23, 32*1024*1024)
		os.Mkdir("./data", 0777)

		// create a chunker
		mychunker := chunker.New(f, chunker.Pol(0x3DA3358B4DC173))

		// reuse this buffer
		buf := make([]byte, 8*1024*1024)

		// os.MkdirAll("data/tree")
		treePath := fmt.Sprintf("data/versions.toml")
		h, _ := os.Create(treePath)

		fmt.Fprintf(h, "[versions.2020-05-29]\n")
		fmt.Fprintf(h, "message=\"%s\"\n", msg)
		fmt.Fprintf(h, "file=\"%s\"\n", filePath)
		fmt.Fprintf(h, "chunks = [\n")

		i := 0

		for {
			chunk, err := mychunker.Next(buf)
			if err == io.EOF {
				break
			}

			if err != nil {
				panic(err)
			}
			
			i += 1
			myHash := sha256.Sum256(chunk.Data)
			fmt.Printf("Chunk %d: %d kB, %02x\n", i, chunk.Length/1024, myHash)
			fmt.Fprintf(h, "  \"%02x\",\n", myHash)

			chunkFolder := fmt.Sprintf("data/%02x", myHash[0:1])
			os.MkdirAll(chunkFolder, 0777)

			chunkPath := fmt.Sprintf("%s/%02x.gz", chunkFolder, myHash)
			g0, _ := os.Create(chunkPath)
			g := gzip.NewWriter(g0)
			g.Write(chunk.Data)
			g.Close()
			g0.Close()
		}

		f.Close()
		fmt.Fprintf(h, "]\n")
		h.Close()
	}
}
