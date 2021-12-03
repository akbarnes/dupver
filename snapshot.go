package dupver

import "time"

const PACK_ID_LEN int = 64

type Snapshot struct {
	Message string
	Time string
    // SnapshotId string // Is this needed?
}

type SnapshotFile struct {
    Md5Hash string
    Size int64
    ModTime string
    ChunkIds []string
}

//type SnapshotFiles struct {
//    Properties map[string]FileProperties
//}
//
//type SnapshotTrees struct {
//    Packs map[string][]string
//}

func CreateSnapshot(message string) (Snapshot, string) {
	t := time.Now()
	ts := t.Format("2006-01-02T15-04-05")
	snap := Snapshot{Time: ts, Message: message}
	return snap, ts
}

//func (snap Snapshot) AddFileChunkIds(head Snapshot, fileName string) {
//	snap.FileChunkIds[fileName] = head.FileChunkIds[fileName]
//
//	for _, chunkId := range snap.FileChunkIds[fileName] {
//		snap.ChunkPackIds[chunkId] = head.ChunkPackIds[chunkId]
//	}
//}
