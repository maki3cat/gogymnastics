package main

import (
	"fmt"
	"testing"
)

func TestNormalPath(t *testing.T) {
	r := NewRouter()
	r.HandleRequest("/user")
	r.RegisterFunc("/user", func(string) string {
		fmt.Println("the path /user is called")
		return ""
	})
	r.HandleRequest("/user")
}

func TestPathWithParameter(t *testing.T) {
	r := NewRouter()
	r.HandleRequest("/user/123/location")
	r.RegisterFunc("/user/{userID}/location", func(userID string) string {
		fmt.Println("the path /user/{userID}/localtion is called")
		fmt.Println("found the userID is", userID)
		return ""
	})
	r.HandleRequest("/user/123/location")
}
