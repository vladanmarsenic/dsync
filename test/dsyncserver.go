package test

import (
	"errors"
	"sync"
	"log"
	"github.com/minio/dsync"
)

type LockServer struct {
	resourceLockMap      map[string]*sync.Mutex
	address string
}

func (ls *LockServer) RLock(args *dsync.LockArgs, reply *bool) error {
	return errors.New("Not yet supported")
}

func (ls *LockServer) RUnlock(args *dsync.LockArgs, reply *bool) error {
	return errors.New("Not yet supported")
}

func (ls *LockServer) Lock(args *dsync.LockArgs, reply *bool) error {
	mutex, ok :=ls.resourceLockMap[args.Resource]
	if !ok {
		mutex = &sync.Mutex{}
		ls.resourceLockMap[args.Resource] = mutex
	}
	mutex.Lock();

	log.Println("Resource " + args.Resource + " locked @" +ls.address)
	*reply = true
	return nil
}

func (ls *LockServer) Unlock(args *dsync.LockArgs, reply *bool) error {
	mutex, ok :=ls.resourceLockMap[args.Resource]
	if !ok {
		return errors.New("Resource "+args.Resource+" does not exist @" + ls.address)
	}
	mutex.Unlock()

	log.Println("Resource " + args.Resource + " unlocked @" + ls.address)
	*reply = true
	return nil
}

func (ls *LockServer) ForceUnlock(args *dsync.LockArgs, reply *bool) error {
	*reply = false
	return errors.New("ForceUnlock is not supported")
}