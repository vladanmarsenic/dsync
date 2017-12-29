package main

import (
	"log"
	"dsyncTest/test"
	"time"
	"bufio"
	"os"
	"fmt"
	"math/rand"
	"net/rpc"
	"github.com/minio/dsync"
	"strconv"
)

var (
	_ test.Locker = &test.DsyncLocker{}
)

func main(){

	testDsync()
	//testRpc()

}

func testDsync(){
	go test.StartLockServer("localhost:50000","/locker0")
	go test.StartLockServer("localhost:50001","/locker1")
	go test.StartLockServer("localhost:50002","/locker2")
	go test.StartLockServer("localhost:50003","/locker3")

	os.Setenv("LOCK_SERVERS", "localhost:50000/locker0, localhost:50001/locker1, localhost:50002/locker2, localhost:50003/locker3")
	os.Setenv("DSYNC_OWN_NODE", "0")

	defer os.Unsetenv("LOCK_SERVERS")
	defer os.Unsetenv("DSYNC_OWN_NODE")


	locker, err := test.NewDsyncLocker()
	if err != nil {
		log.Fatal("Locker could not be initialized")
	}


	resources := getResources()
	r := rand.New(rand.NewSource(99))
	for i:=0; i<10; i++ {
		dur:=time.Duration(r.Intn(5000))
		i := r.Intn(cap(resources))
	 	go lockUnlock(locker, dur * time.Millisecond, resources[i], i)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

func testRpc(){

	go test.StartLockServer("localhost:50000","/locker0")
	time.Sleep(2*time.Second)

	client, error := rpc.DialHTTPPath("tcp", "localhost:50000", "/locker0")
	if error != nil {
		log.Fatal("dialing:", error)
	}

	var result bool
	error = client.Call("LockServer.Lock", dsync.LockArgs{Resource:"Music"}, &result)
	if error != nil {
		log.Println(error)
	}

	error = client.Call("LockServer.Unlock", dsync.LockArgs{Resource:"Music"}, &result)
	if error != nil {
		log.Println(error)
	}
}

func getResources() []string{
	resources := make([]string, 1000)
	for i:=0; i< 1000; i++ {
		resources[i]="resource_"+strconv.Itoa(i)
	}
	return resources
}

func lockUnlock(locker test.Locker, d time.Duration, name string, i int){
	locker.Lock(name)
	defer locker.Release(name)
	time.Sleep(d)
	log.Println(name + "_" + strconv.Itoa(i))
}