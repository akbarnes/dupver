package dupver

import (
	"encoding/json"
	"path"
	"path/filepath"
	"fmt"
	// "io"
	"os"
)

// Write out a JSON tree file, given a chunk to pack map
// and a filename for the tree
func WriteTree(treePath string, chunkPacks map[string]string) {
	f, err := os.Create(treePath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create tree file %s", treePath))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(chunkPacks)
	f.Close()
}

// Read a JSON tree file into a chunk-pack map 
// given a repository path
func ReadTrees(repoPath string) map[string]string {
	treesGlob := path.Join(repoPath, "trees", "*.json")
	// fmt.Println(treesGlob)
	treePaths, err := filepath.Glob(treesGlob)

	if err != nil {
		panic(fmt.Sprintf("Error reading trees %s", treesGlob))
	}

	chunkPacks := make(map[string]string)

	for _, treePath := range treePaths {
		treePacks := make(map[string]string)

		f, err := os.Open(treePath)

		if err != nil {
			panic(fmt.Sprintf("Error: could not read tree file %s", treePath))
		}

		myDecoder := json.NewDecoder(f)

		if err := myDecoder.Decode(&treePacks); err != nil {
			panic(fmt.Sprintf("Error: could not decode tree file %s", treePath))
		}

		// TODO: handle supersedes to allow repacking files
		for k, v := range treePacks {
			chunkPacks[k] = v
		}

		f.Close()
	}

	return chunkPacks
}

