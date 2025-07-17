package medium

import (
	"fmt"
	"sync"
	"testing"
)

func TestGetConnectionTwice(t *testing.T) {

	fmt.Println("---1st TESTING---")
	nodeID := "node2"

	// Test the connection establishment
	conn, err := connManager.GetConnection(nodeID)
	if err != nil {
		t.Fatalf("Failed to establish connection: %v", err)
	}
	if conn != "the Connection to localhost:8081" {
		t.Fatalf("Expected 'the Connection to localhost:8081', got '%s'", conn)
	}

	// Test the GetConnection function
	conn, err = connManager.GetConnection(nodeID)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}
	if conn != "the Connection to localhost:8081" {
		t.Fatalf("Expected 'the Connection to localhost:8081', got '%s'", conn)
	}
	// How to check the EstablishConn is called for only once?
}

func TestParallelGetConnection(t *testing.T) {
	fmt.Println("---2nd TESTING---")
	nodeID := "node3"
	// Test concurrent access to GetConnection
	var wg sync.WaitGroup
	getConn := func(threadID int) {
		defer wg.Done()
		connManager.GetConnection(nodeID)
		fmt.Printf("thread %d: get the connection to %s\n", threadID, nodeID)
	}
	for i := range 10 {
		wg.Add(1)
		go getConn(i)
	}
	wg.Wait()
}
