package dupver

import "time"

const PACK_ID_LEN int = 64

type Snapshot struct {
	Message      string
	Time         string
}

type FileProperties struct {
    Hash string
    Length int64
    ModificationTime string
    CreationTime string
    Owner string
    Group string
    Permissions string
}

type SnapshotFiles struct {
	ChunkPackIds map[string]string
	FileChunkIds map[string][]string
	FileModTimes map[string]string
}

func CreateSnapshot(message string) (Snapshot, string) {
	t := time.Now()
	ts := t.Format("2006-01-02T15-04-05")
	snap := Snapshot{Time: ts, Message: message}
	snap.ChunkPackIds = make(map[string]string)
	snap.FileChunkIds = make(map[string][]string)
	snap.FileModTimes = make(map[string]string)
	return snap, ts
}

func (snap Snapshot) AddFileChunkIds(head Snapshot, fileName string) {
	snap.FileChunkIds[fileName] = head.FileChunkIds[fileName]

	for _, chunkId := range snap.FileChunkIds[fileName] {
		snap.ChunkPackIds[chunkId] = head.ChunkPackIds[chunkId]
	}
}
