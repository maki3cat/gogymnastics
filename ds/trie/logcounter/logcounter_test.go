package logcounter

import (
	"fmt"
	"testing"
)

func TestLogParser(t *testing.T) {
	lines := []string{
		"[02/Nov/2018:21:46:31 +0000] PUT /users/12345/locations HTTP/1.1 204 iphone-3",
		"[02/Nov/2018:21:46:31 +0000] PUT /users/6098/locations HTTP/1.1 204 iphone-3",
		"[02/Nov/2018:21:46:32 +0000] PUT /users/3911/locations HTTP/1.1 204 moto-x",
		"[02/Nov/2018:21:46:33 +0000] PUT /users/9933/locations HTTP/1.1 404 moto-x",
		"[02/Nov/2018:21:46:33 +0000] PUT /users/3911/locations HTTP/1.1 500 moto-x",
	}
	root := NewNode("")
	for _, line := range lines {
		log := parseLog(line)
		CountLog(root, log)
	}
	fmt.Println("--------------------------------")
	PrintCounts(root)
}

func TestFindRealPath(t *testing.T) {
	rawPath := "/users/12345/locations"
	realPath := findRealPath(rawPath)
	fmt.Println(realPath)
}
