package test

import (
	"github.com/minio/dsync"
	"os"
	"strings"
	"time"
	"github.com/pkg/errors"
	"strconv"
	"log"
	"sync"
)

const (
	lockServersEnvKey  = "LOCK_SERVERS"
	dsyncOwnNodeEnvKey = "DSYNC_OWN_NODE"
	defaultLockTimeout = 1 * time.Second
)

var (
	_ Locker = &DsyncLocker{}
	ErrLockNotAcquired = errors.New("lock not acquired")
	ErrLockNotReleased = errors.New("lock not released")
)

type DsyncLocker struct {
	ds          *dsync.Dsync
	drwMutexMap map[string]*dsync.DRWMutex
	globalMutex *sync.Mutex
}

func NewDsyncLocker() (locker Locker, err error) {
	lockServers := strings.Split(os.Getenv(lockServersEnvKey), ",")
	dsyncOwnNode, err := strconv.Atoi(os.Getenv(dsyncOwnNodeEnvKey))
	if err != nil {
		return
	}

	lockClients := make([]dsync.NetLocker, cap(lockServers))
	for i, address := range lockServers {
		param := strings.SplitN(address, "/", 2)
		serverAddress := strings.TrimSpace(param[0])
		serviceEndpoint := "/"+ strings.TrimSpace(param[1])
		lockClients[i] = NewDsyncClient(serverAddress, serviceEndpoint)
	}

	ds, err := dsync.New(lockClients, dsyncOwnNode)
	if err != nil {
		return
	}
	locker = &DsyncLocker{ds, make(map[string]*dsync.DRWMutex), &sync.Mutex{}}
	return
}

// Maybe i do not have to take care of locks here!!!
func (d *DsyncLocker) Lock(resource string) error {
	mutex, ok := d.drwMutexMap[resource]
	if !ok {
		mutex = dsync.NewDRWMutex(resource, d.ds)
	}
	granted := mutex.GetLock(defaultLockTimeout)
	if !granted {
		return ErrLockNotAcquired
	}
	d.drwMutexMap[resource] = mutex
	return nil
}

func (d *DsyncLocker) Release(resource string) error {
	mutex, ok := d.drwMutexMap[resource]
	if !ok {
		return ErrLockNotReleased
	}
	mutex.Unlock()
	return nil
}

func (d *DsyncLocker) Status()  {
	log.Print("not yet implemented")
}
