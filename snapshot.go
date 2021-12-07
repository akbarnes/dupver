package dupver

import "time"

const SNAPSHOT_ID_LEN int = 40
const PACK_ID_LEN int = 64

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
	sid := RandHexString(SNAPSHOT_ID_LEN)
	snap := Snapshot{SnapshotTime: ts, Message: message, SnapshotId: sid}
	return snap
}

func AddFileChunkIds(files map[string]SnapshotFile, headFiles map[string]SnapshotFile, fileName string) {
	files[fileName] = headFiles[fileName]

	// for _, chunkId := range files.ChunkIds[fileName] {
	// 	snap.ChunkPackIds[chunkId] = head.ChunkPackIds[chunkId]
	// }
}
