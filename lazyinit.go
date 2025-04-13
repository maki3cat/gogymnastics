package gogymnastics

import (
	"fmt"
	"sync"
)

// Pattern (1) write once at starting time, and
// read only during the lifetime of the lifetime of the server
// in this scenario, we can use a map to store the address of the nodes

// Question: if a readonly map safe for concurrent use?
// Answer: yes, a read-only map is safe for concurrent use.
var fixedAddress map[string]string
var fixedMutex map[string]*sync.Mutex

// Pattern (2) the connections have read and writes/updates for dead connection
// during the runtime of the server, so we need sync map
var connections *sync.Map

func init() {
	// set this fixedAddress only once
	fixedAddress = map[string]string{
		"node1": "localhost:8080",
		"node2": "localhost:8081",
		"node3": "localhost:8082",
		"node4": "localhost:8083",
		"node5": "localhost:8084",
	}
	fixedMutex = make(map[string]*sync.Mutex)
	for k := range fixedAddress {
		fixedMutex[k] = &sync.Mutex{}
	}
	// initialize the connections
	connections = &sync.Map{}
}

func getAddress(nodeName string) string {
	return fixedAddress[nodeName]
}

func mockingEstablishConn(address string) (string, error) {
	fmt.Printf("Establishing connection to %s\n", address)
	return "the Connection to " + address, nil
}

func GetConnection(nodeName string) (string, error) {
	conn, ok := connections.Load(nodeName)
	if ok {
		return conn.(string), nil
	}
	// only one goroutine can establish a connection
	fixedMutex[nodeName].Lock()
	defer fixedMutex[nodeName].Unlock()

	// check again if the connection is already established
	conn, ok = connections.Load(nodeName)
	if ok {
		return conn.(string), nil
	}
	// establish the connection
	address := getAddress(nodeName)
	conn, err := mockingEstablishConn(address)
	if err != nil {
		return "", err
	}
	// store the connection in the map
	connections.Store(nodeName, conn)
	return conn.(string), nil
}
