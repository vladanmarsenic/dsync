package test


import (
	"sync"
	"net/rpc"
	"net"
	"log"
	"net/http"
)

func StartLockServer(serverAddress, serviceEndpoint string) {

	lockServer := &LockServer{
		resourceLockMap:      make(map[string]*sync.Mutex),
		address: serverAddress,
	}

	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("LockServer", lockServer)
	rpcServer.HandleHTTP(serviceEndpoint, serviceEndpoint + "/debug")

	listener, err := net.Listen("tcp", serverAddress)
	if err == nil {
		log.Println("LockServer listening @" + serverAddress + serviceEndpoint)
		http.Serve(listener, nil)
	}

	log.Fatal("Unable to start LockServer @ %s%s", serverAddress, serviceEndpoint)
}
