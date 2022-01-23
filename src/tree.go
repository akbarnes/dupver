package dupver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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
			tree[packId] = []string{chunkId}
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
