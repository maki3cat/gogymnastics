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

func TestWildCardMathching(t *testing.T) {
	r := NewRouter()
	r.RegisterFunc("/*", func(string) string {
		fmt.Println("wildcard matching")
		return ""
	})
	r.HandleRequest("/user/123/location")
	r.HandleRequest("/location")
}

func TestWildCardMathchingComplex(t *testing.T) {
	r := NewRouter()
	r.RegisterFunc("/user/*", func(string) string {
		fmt.Println("user wildcard matching")
		return ""
	})
	r.RegisterFunc("/*", func(string) string {
		fmt.Println("general wildcard matching")
		return ""
	})
	r.HandleRequest("/user/account")
	r.HandleRequest("/location")
}

func TestWildCardMathchingComplex2(t *testing.T) {
	r := NewRouter()
	r.RegisterFunc("/user/*", func(string) string {
		fmt.Println("user wildcard matching")
		return ""
	})
	r.RegisterFunc("/*", func(string) string {
		fmt.Println("general wildcard matching")
		return ""
	})
	r.RegisterFunc("/user/account", func(string) string {
		fmt.Println("precise matching")
		return ""
	})
	r.HandleRequest("/user/account")
}
func TestRouter(t *testing.T) {
	router := NewRouter()

	router.RegisterFunc("/hello", func(_ string) string {
		return "exact hello"
	})

	router.RegisterFunc("/user/{id}", func(id string) string {
		return "user id=" + id
	})

	router.RegisterFunc("/files/*", func(_ string) string {
		return "wildcard files"
	})

	router.RegisterFunc("/files/special", func(_ string) string {
		return "special file"
	})

	tests := []struct {
		path     string
		expected string
	}{
		// Exact match
		{"/hello", "exact hello"},

		// Digit param
		{"/user/123", "user id=123"},

		// Wildcard
		{"/files/anything/here", "wildcard files"},

		// Exact match beats wildcard
		{"/files/special", "special file"},

		// Path not found
		{"/not/exist", ""},
	}

	for _, tt := range tests {
		param, handler := router.findHandler(tt.path)
		var got string
		if handler != nil {
			got = (*handler)(param)
		}

		if got != tt.expected {
			t.Errorf("path=%q: expected %q, got %q", tt.path, tt.expected, got)
		}
	}
}
