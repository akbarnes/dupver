package dupver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
    "time"
)

const SnapshotIDLen int = 40
const PackIdLen int = 64

// TODO: change this to SerializedSnaphot
// and use Time type for SnapshotTime?
type Snapshot struct {
	Message      string
	SnapshotTime string
	SnapshotLocalTime string
	SnapshotID   string // Is this needed?
}

type Head struct {
	SnapshotTime string
	SnapshotID   string // Is this needed?
}

type SnapshotFile struct {
    Name     string
	Size     int64
	ModTime  string
	ModLocalTime  string
	ChunkIds []string
    IsArchive bool
}

// Snapshot Files
// files := map[string]SnapshotFile{}

// Pack for each Chunk
// packs := map[string]string{}

func CreateSnapshot(message string) Snapshot {
	t := time.Now().UTC()
	tl := time.Now().Local()
	ts := t.Format("2006-01-02T15-04-05")
	tsl := tl.Format("2006-01-02T15-04-05")
	sid := RandHexString(SnapshotIDLen)
	snap := Snapshot{SnapshotTime: ts, SnapshotLocalTime: tsl, Message: message, SnapshotID: sid}
	return snap
}

func AddFileChunkIds(files map[string]SnapshotFile, headFiles map[string]SnapshotFile, fileName string) {
	files[fileName] = headFiles[fileName]
}

// TODO: return err instead of panic?
func (snap Snapshot) Write() {
	snapFolder := filepath.Join(".dupver", "snapshots")

	if err := os.MkdirAll(snapFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating snapshot folder %s\n", snapFolder))
	}

	snapFile := filepath.Join(snapFolder, snap.SnapshotID+".json")
	f, err := os.Create(snapFile)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot json %s", snapFile))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(snap)
	f.Close()
}

func (snap Snapshot) WriteFiles(files []SnapshotFile) {
	filesFolder := filepath.Join(".dupver", "files")

	if err := os.MkdirAll(filesFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating files listing folder %s\n", filesFolder))
	}

	snapFile := filepath.Join(filesFolder, snap.SnapshotID+".json")
	f, err := os.Create(snapFile)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot files listing json %s", snapFile))
	}

	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(files)
	f.Close()
}

func (snap Snapshot) ReadFilesHash() map[string]SnapshotFile {
	if snap.SnapshotID == "" {
		return map[string]SnapshotFile{}
	}

    files := snap.ReadFilesList()
    fileHash := map[string]SnapshotFile{}

    for _, fileProps := range files {
        fileHash[fileProps.Name] = fileProps
    }

    return fileHash
}

func (snap Snapshot) ReadFilesList() []SnapshotFile {
	if snap.SnapshotID == "" {
		return []SnapshotFile{}
	}

	filesFolder := filepath.Join(".dupver", "files")

	if err := os.MkdirAll(filesFolder, 0777); err != nil {
		panic(fmt.Sprintf("Error creating files listing folder %s\n", filesFolder))
	}

	snapFile := filepath.Join(filesFolder, snap.SnapshotID+".json")
	f, err := os.Open(snapFile)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot files listing json %s", snapFile))
	}

	myDecoder := json.NewDecoder(f)
	files := []SnapshotFile{}

	if err := myDecoder.Decode(&files); err != nil {
		panic(fmt.Sprintf("Error: could not decode snapshot files %s\n", snapFile))
	}

	f.Close()
	return files
}


func ReadSnapshot(snapId string) Snapshot {
	snapshotPath := filepath.Join(".dupver", "snapshots", snapId+".json")

	if DebugMode {
		fmt.Fprintf(os.Stderr, "Reading %s\n", snapshotPath)
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
		fmt.Fprintf(os.Stderr, "Error:could not decode head file %s\n", snapshotPath)
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

	head := Head{SnapshotTime: snap.SnapshotTime, SnapshotID: snap.SnapshotID}

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
		fmt.Fprintf(os.Stderr, "Error:could not decode head file %s\n", headPath)
		Check(err)
	}

	f.Close()

	return ReadSnapshot(head.SnapshotID)
}
