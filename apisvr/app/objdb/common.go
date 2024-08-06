package objdb

import "sync"

// ---------------------------------------------------------------------------
// mutex
// ---------------------------------------------------------------------------
var shmMutex = &sync.Mutex{}

func SetShmSvrState(info SystemInfo) {

	shmMutex.Lock()
	SharedMem.System.SvrUtc = info.SvrUtc
	shmMutex.Unlock()
}
