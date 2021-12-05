package dupver

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/restic/chunker"
)

// TODO: return err instead of panic?
func (snap Snapshot) Write() {
	snapFolder := filepath.Join(".dupver", "snapshots")

	if err := os.MkdirAll(snapFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating snapshot folder %s\n", snapFolder))
	}

	snapFile := filepath.Join(snapFolder, snap.SnapshotId+".json")
	f, err := os.Create(snapFile)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot json %s", snapFile))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(snap)
	f.Close()
}

func (snap Snapshot) WriteFiles(files map[string]SnapshotFile) {
	filesFolder := filepath.Join(".dupver", "files")

	if err := os.MkdirAll(filesFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating files listing folder %s\n", filesFolder))
	}

	snapFile := filepath.Join(filesFolder, snap.SnapshotId+".json")
	f, err := os.Create(snapFile)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot files listing json %s", snapFile))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(files)
	f.Close()
}

func (snap Snapshot) WriteTree(packs map[string]string) {
	treesFolder := filepath.Join(".dupver", "trees")

	if err := os.MkdirAll(treesFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating trees folder %s\n", treesFolder))
	}

	treeFile := filepath.Join(treesFolder, snap.SnapshotId+".json")
	f, err := os.Create(treeFile)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot tree json %s", treeFile))
	}

	tree := map[string][]string{}

	// Remember that I'll only encounter each chunk id once
	for chunkId, packId := range packs {
		if _, ok := tree[packId]; ok {
			tree[packId] = append(tree[packId], chunkId)
		} else {
			tree[packId] = []string{}
		}
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(tree)
	f.Close()
}

func ReadTrees() map[string]string {
	treesFolder := filepath.Join(".dupver", "trees")

	if err := os.MkdirAll(treesFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating trees folder %s\n", treesFolder))
	}

	treesGlob := filepath.Join(treesFolder, "*.json")
	treePaths, err := filepath.Glob(treesGlob)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not glob trees %s", treesGlob))
	}

	packs := map[string]string{}

	for _, treePath := range treePaths {
		tree := ReadTree(treePath)

		for packId, chunkIds := range tree {
			for _, chunkId := range chunkIds {
				packs[chunkId] = packId
			}
		}
	}

	return packs
}

// Read a tree given a file path
func ReadTree(treePath string) map[string][]string {
	tree := map[string][]string{}
	f, err := os.Open(treePath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not read tree file %s", treePath))
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&tree); err != nil {
		panic(fmt.Sprintf("Error: could not decode tree file %s\n", treePath))
	}

	f.Close()
	return tree
}

// func ReadSnapshot(snapId string) Snapshot {
// 	snapshotPath := filepath.Join(".dupver", "snapshots", snapId+".json")

// 	if VerboseMode {
// 		fmt.Printf("Reading %s\n", snapshotPath)
// 	}

// 	return ReadSnapshotFile(snapId)
// }

// // Read a snapshot given a file path
// func ReadSnapshotFile(snapshotPath string) Snapshot {
// 	var snap Snapshot
// 	f, err := os.Open(snapshotPath)

// 	// ChunkPackIds map[string]string
// 	// FileChunkIds map[string][]string
// 	// FileModTimes map[string]string

// 	if err != nil {
// 		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
// 		snap := Snapshot{}
// 		snap.ChunkPackIds = make(map[string]string)
// 		snap.FileChunkIds = make(map[string][]string)
// 		snap.FileModTimes = make(map[string]string)
// 		return snap
// 	}

// 	myDecoder := json.NewDecoder(f)

// 	if err := myDecoder.Decode(&snap); err != nil {
// 		fmt.Printf("Error:could not decode head file %s\n", snapshotPath)
// 		Check(err)
// 	}

// 	f.Close()
// 	return snap
// }

// func WriteHead(snapshotPath string) {
// 	headPath := filepath.Join(".dupver", "head.json")
// 	f, err := os.Create(headPath)

// 	if err != nil {
// 		panic(fmt.Sprintf("Error: Could not create head file %s", headPath))
// 	}

// 	myEncoder := json.NewEncoder(f)
// 	myEncoder.SetIndent("", "  ")
// 	myEncoder.Encode(snapshotPath)
// 	f.Close()
// }

// // Read a snapshot given a file path
// func ReadHead() Snapshot {
// 	headPath := filepath.Join(".dupver", "head.json")
// 	f, err := os.Open(headPath)

// 	if err != nil {
// 		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
// 		snap := Snapshot{}
// 		snap.ChunkPackIds = make(map[string]string)
// 		snap.FileChunkIds = make(map[string][]string)
// 		snap.FileModTimes = make(map[string]string)
// 		return snap
// 	}

// 	snapshotId := ""
// 	myDecoder := json.NewDecoder(f)

// 	if err := myDecoder.Decode(&snapshotId); err != nil {
// 		fmt.Printf("Error:could not decode head file %s\n", headPath)
// 		Check(err)
// 	}

// 	f.Close()

// 	snapshotPath := filepath.Join(".dupver", "snapshots", snapshotId+".json")
// 	return ReadSnapshotFile(snapshotPath)
// }

func CreatePackFile(packId string) (*os.File, error) {
	packFolderPath := filepath.Join(".dupver", "packs", packId[0:2])
	os.MkdirAll(packFolderPath, 0777)
	packPath := filepath.Join(packFolderPath, packId+".zip")

	if VerboseMode {
		fmt.Printf("Creating pack: %s\n", packId[0:16])
	}

	// TODO: only create pack file if we need to save stuff - set to nil initially
	packFile, err := os.Create(packPath)

	if err != nil {
		if VerboseMode {
			fmt.Printf("Error creating pack file %s", packPath)
		}

		return nil, err
	}

	return packFile, nil
}

func WriteChunkToPack(zipWriter *zip.Writer, chunkId string, chunk chunker.Chunk) error {
	var header zip.FileHeader
	header.Name = chunkId
	header.Method = CompressionLevel

	writer, err := zipWriter.CreateHeader(&header)

	if err != nil {
		if VerboseMode {
			fmt.Printf("Error creating zip header\n")
		}

		return err
	}

	if _, err := writer.Write(chunk.Data); err != nil {
		if VerboseMode {
			fmt.Printf("Error writing chunk %s to zip file\n", chunkId)
		}

		return err
	}

	return nil
}

// func ExtractChunkFromPack(outFile *os.File, chunkId string, packId string) error {
// 	dupverDir := filepath.Join(WorkingDirectory, ".dupver")
// 	packFolderPath := path.Join(dupverDir, "packs", packId[0:2])
// 	packPath := path.Join(packFolderPath, packId+".zip")
// 	packFile, err := zip.OpenReader(packPath)

// 	if err != nil {
// 		if VerboseMode {
// 			fmt.Printf("Error extracting pack %s[%s]\n", packId, chunkId)
// 		}
// 		return err
// 	}

// 	defer packFile.Close()
// 	return ExtractChunkFromZipFile(outFile, packFile, chunkId)
// }

// func ExtractChunkFromZipFile(outFile *os.File, packFile *zip.ReadCloser, chunkId string) error {
// 	for _, f := range packFile.File {

// 		if f.Name == chunkId {
// 			// fmt.Printf("Contents of %s:\n", f.Name)
// 			chunkFile, err := f.Open()

// 			if err != nil {
// 				if VerboseMode {
// 					fmt.Printf("Error opening chunk %s\n", chunkId)
// 				}

// 				return err
// 			}

// 			_, err = io.Copy(outFile, chunkFile)

// 			if err != nil {
// 				if VerboseMode {
// 					fmt.Printf("Error reading chunk %s\n", chunkId)
// 				}

// 				return err
// 			}

// 			chunkFile.Close()
// 		}
// 	}

// 	return nil
// }
