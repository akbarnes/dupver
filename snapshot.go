package dupver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
    "time"
)

const SnapshotIdLen int = 40
const PackIdLen int = 64

// TODO: change this to SerializedSnaphot
// and use Time type for SnapshotTime?
type Snapshot struct {
	Message      string
	SnapshotTime string
	SnapshotId   string // Is this needed?
}

type Head struct {
	SnapshotTime string
	SnapshotId   string // Is this needed?
}

type SnapshotFile struct {
	Size     int64
	ModTime  string
	ChunkIds []string
}

// Snapshot Files
// files := map[string]SnapshotFile{}

// Pack for each Chunk
// packs := map[string]string{}

func CreateSnapshot(message string) Snapshot {
	t := time.Now()
	ts := t.Format("2006-01-02T15-04-05")
	sid := RandHexString(SnapshotIdLen)
	snap := Snapshot{SnapshotTime: ts, Message: message, SnapshotId: sid}
	return snap
}

func AddFileChunkIds(files map[string]SnapshotFile, headFiles map[string]SnapshotFile, fileName string) {
	files[fileName] = headFiles[fileName]

	// for _, chunkId := range files.ChunkIds[fileName] {
	// 	snap.ChunkPackIds[chunkId] = head.ChunkPackIds[chunkId]
	// }
}

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

func (snap Snapshot) ReadFilesList() map[string]SnapshotFile {
	if snap.SnapshotId == "" {
		return map[string]SnapshotFile{}
	}

	filesFolder := filepath.Join(".dupver", "files")

	if err := os.MkdirAll(filesFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating files listing folder %s\n", filesFolder))
	}

	snapFile := filepath.Join(filesFolder, snap.SnapshotId+".json")
	f, err := os.Open(snapFile)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot files listing json %s", snapFile))
	}

	myDecoder := json.NewDecoder(f)
	files := map[string]SnapshotFile{}

	if err := myDecoder.Decode(&files); err != nil {
		panic(fmt.Sprintf("Error: could not decode snapshot files %s\n", snapFile))
	}

	f.Close()
	return files
}

func ReadSnapshot(snapId string) Snapshot {
	snapshotPath := filepath.Join(".dupver", "snapshots", snapId+".json")

	if VerboseMode {
		fmt.Printf("Reading %s\n", snapshotPath)
	}

	return ReadSnapshotJson(snapshotPath)
}

// Read a snapshot given a file path
func ReadSnapshotJson(snapshotPath string) Snapshot {
	var snap Snapshot
	f, err := os.Open(snapshotPath)

	if err != nil {
		return Snapshot{}
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&snap); err != nil {
		fmt.Printf("Error:could not decode head file %s\n", snapshotPath)
		Check(err)
	}

	f.Close()
	return snap
}

func (snap Snapshot) WriteHead() {
	headPath := filepath.Join(".dupver", "head.json")
	f, err := os.Create(headPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create head file %s", headPath))
	}

	head := Head{SnapshotTime: snap.SnapshotTime, SnapshotId: snap.SnapshotId}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(head)
	f.Close()
}

// Read all snapshots and sort by date
func ReadAllSnapshots() []Snapshot {
	snapshotGlob := filepath.Join(".dupver", "snapshots", "*.json")
	snapshotPaths, err := filepath.Glob(snapshotGlob)
	Check(err)
	snaps := []Snapshot{}

	for _, snapshotPath := range snapshotPaths {
		snaps = append(snaps, ReadSnapshotJson(snapshotPath))
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].SnapshotTime < snaps[j].SnapshotTime
	})

	return snaps
}

// Read the head snapshot and files list
func ReadHead() Snapshot {
	headPath := filepath.Join(".dupver", "head.json")
	f, err := os.Open(headPath)

	if err != nil {
		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
		return Snapshot{}
	}

	head := Head{}
	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&head); err != nil {
		fmt.Printf("Error:could not decode head file %s\n", headPath)
		Check(err)
	}

	f.Close()

	return ReadSnapshot(head.SnapshotId)
}
