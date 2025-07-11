package main

import (
	"fmt"
	"sync"
)

// COMMENTS AND EVALUATION FROM GPT-4

// CORRECTNESS
// 1. Read-only Map for Fixed Addresses
// You're right that a map which is only written once and then read concurrently is safe
// as long as no writes happen after initialization.
// Go's memory model guarantees safe reads from a fully-initialized map
// if no further writes occur. ✅

// 2. Per-node Mutexes for Synchronizing Connection Creation
// Using a per-node *sync.Mutex is a solid choice
// for reducing contention compared to a global mutex.
// This ensures that only one goroutine per node establishes the connection,
// preventing redundant work. ✅

// 3. Double-checked Locking
// You're doing the "double-checked locking" pattern:
// First check without lock (fast path).
// Lock and re-check to avoid duplicate work. ✅ Correctly implemented.

// 4. Sync.Map for Concurrent Read/Write
// snyc.Map is ideal here because you're reading/writing connections concurrently,
// and its performance is good for cases with many reads and occasional writes. ✅ Good choice.

// EFFICIENT ASPECTS:
// 1. Fine-grained Synchronizaiton/Locking using fixedMutex ✅
// 2. Avoidance of Unnecessary Synchronization: Avoiding locking on read path ✅
// 3. Caching: Avoiding redundant work by caching connections ✅

type ConnectionManager struct {
	// Pattern (1) write once at starting time, and
	// read only during the lifetime of the lifetime of the server
	// in this scenario, we can use a map to store the address of the nodes
	addresses map[string]string
	mutexes   map[string]*sync.Mutex

	// Pattern (2) the connections have read and writes/updates for dead connection
	// during the runtime of the server, so we need sync map
	connections *sync.Map
}

var connManager *ConnectionManager

func init() {
	connManager = &ConnectionManager{}
	// set this fixedAddress only once
	connManager.addresses = map[string]string{
		"node1": "localhost:8080",
		"node2": "localhost:8081",
		"node3": "localhost:8082",
		"node4": "localhost:8083",
		"node5": "localhost:8084",
	}
	connManager.mutexes = make(map[string]*sync.Mutex)
	for k := range connManager.addresses {
		connManager.mutexes[k] = &sync.Mutex{}
	}
	// initialize the connections
	connManager.connections = &sync.Map{}
}

func (c *ConnectionManager) getAddress(nodeName string) string {
	return c.addresses[nodeName]
}

func (c *ConnectionManager) mockingEstablishConn(address string) (string, error) {
	fmt.Printf("Establishing connection to %s\n", address)
	return "the Connection to " + address, nil
}

func (c *ConnectionManager) GetConnection(nodeName string) (string, error) {
	// to avoid future panics
	if c.addresses[nodeName] == "" {
		return "", fmt.Errorf("node %s does not exist", nodeName)
	}

	conn, ok := c.connections.Load(nodeName)
	if ok {
		return conn.(string), nil
	}
	// only one goroutine can establish a connection
	c.mutexes[nodeName].Lock()
	defer c.mutexes[nodeName].Unlock()

	// check again if the connection is already established
	conn, ok = c.connections.Load(nodeName)
	if ok {
		return conn.(string), nil
	}
	// establish the connection
	address := c.getAddress(nodeName)
	conn, err := c.mockingEstablishConn(address)
	if err != nil {
		return "", err
	}
	// store the connection in the map
	c.connections.Store(nodeName, conn)
	return conn.(string), nil
}
