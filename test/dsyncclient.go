package test

import (
	"errors"
	"log"
	"net/rpc"
	"sync"

	"github.com/minio/dsync"
)

type DsyncClient struct {
	mutex           sync.Mutex
	netRPCClient    *rpc.Client
	serverAddr      string
	serviceEndpoint string
}

func NewDsyncClient(serverAddr, serviceEndpoint string) *DsyncClient {
	return &DsyncClient{
		serverAddr:      serverAddr,
		serviceEndpoint: serviceEndpoint,
	}
}

func (rpcClient *DsyncClient) RLock(args dsync.LockArgs) (status bool, err error) {
	status = false
	err = errors.New("Not yet implemented")
	return
}

func (rpcClient *DsyncClient) RUnlock(args dsync.LockArgs) (status bool, err error) {
	status = false
	err = errors.New("Not yet implemented")
	return
}

func (rpcClient *DsyncClient) Lock(args dsync.LockArgs) (status bool, err error) {
	err = rpcClient.call("LockServer.Lock", &args, &status)
	if err == rpc.ErrShutdown {
		err = rpcClient.call("LockServer.Lock", &args, &status)
	}
	return status, err
}

func (rpcClient *DsyncClient) Unlock(args dsync.LockArgs) (status bool, err error) {
	err = rpcClient.call("LockServer.Unlock", &args, &status)
	if err == rpc.ErrShutdown {
		err = rpcClient.call("LockServer.Unlock", &args, &status)
	}
	return status, err
}

func (rpcClient *DsyncClient) ForceUnlock(args dsync.LockArgs) (status bool, err error) {
	err = rpcClient.call("LockServer.ForceUnlock", &args, &status)
	if err == rpc.ErrShutdown {
		err = rpcClient.call("LockServer.ForceUnlock", &args, &status)
	}
	return status, err
}

func (rpcClient *DsyncClient) Close() (err error) {
	rpcClient.mutex.Lock()
	netRPCClient := rpcClient.netRPCClient
	rpcClient.mutex.Unlock()

	if netRPCClient != nil {
		rpcClient.mutex.Lock()
		rpcClient.netRPCClient = nil
		rpcClient.mutex.Unlock()
		netRPCClient.Close()
	}
	return nil
}

func (rpcClient *DsyncClient) ServerAddr() string {
	return rpcClient.serverAddr
}

func (rpcClient *DsyncClient) ServiceEndpoint() string {
	return rpcClient.serviceEndpoint
}

func (rpcClient *DsyncClient) call(serviceMethod string, args interface{}, reply interface{}) error {
	netRPCClient, err := rpcClient.dial()
	if err != nil {
		return err
	}
	err = netRPCClient.Call(serviceMethod, args, reply)
	if err == rpc.ErrShutdown {
		netRPCClient.Close()
	}
	return err
}

func (rpcClient *DsyncClient) dial() (client *rpc.Client,error error) {

	rpcClient.mutex.Lock()
	defer rpcClient.mutex.Unlock()

	if rpcClient.netRPCClient != nil {
		return rpcClient.netRPCClient, nil
	}

	client, error = rpc.DialHTTPPath("tcp", rpcClient.serverAddr, rpcClient.serviceEndpoint)
	if error != nil {
		log.Fatal("dialing:", error)
	}

	rpcClient.netRPCClient = client
	return
}
