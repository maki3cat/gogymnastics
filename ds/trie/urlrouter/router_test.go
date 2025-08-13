package urlrouter

import (
	"fmt"
	"testing"
)

func TestRegister(t *testing.T) {
	Register("/users/{userID}/location", "GET", func() {
		fmt.Println("GET /users/{userID}/location")
	})
	handler, ok := Route("/users/123/location", "GET")
	if !ok {
		t.Errorf("Route failed")
	} else {
		handler()
	}
}

// - /users/{:userID:}/location, GET/DELETE
// - /users, GET
// - /users/{:userID:}/account, GET/PATCH

func TestRegister2(t *testing.T) {
	Register("/users/{userID}/location", "GET", func() {
		fmt.Println("GET /users/{userID}/location")
	})
	Register("/users/{userID}/location", "DELETE", func() {
		fmt.Println("DELETE /users/{userID}/location")
	})
	Register("/users", "GET", func() {
		fmt.Println("GET /users")
	})

	// test location
	handler, ok := Route("/users/123/location", "GET")
	if !ok {
		t.Errorf("Route failed")
	} else {
		handler()
	}
	handler, ok = Route("/users/123/location", "DELETE")
	if !ok {
		t.Errorf("Route failed")
	} else {
		handler()
	}
	_, ok = Route("/users/882143/location", "POST")
	if ok {
		t.Errorf("Route failed")
	}

	// users
	handler, ok = Route("/users", "GET")
	if !ok {
		t.Errorf("Route failed")
	} else {
		handler()
	}
	_, ok = Route("/users", "POST")
	if ok {
		t.Errorf("Route failed")
	}
}

func TestRegister3(t *testing.T) {
	Register("/cities/*", "GET", func() {
		fmt.Println("GET /cities/*")
	})
	handler, ok := Route("/cities/shanghai", "GET")
	if !ok {
		t.Errorf("Route failed")
	} else {
		handler()
	}
}

// wildcard matching has lower priority
func TestRegister4(t *testing.T) {
	Register("/cities/*", "GET", func() {
		fmt.Println("GET /cities/*")
	})
	Register("/cities/shanghai", "GET", func() {
		fmt.Println("GET /cities/shanghai")
	})
	handler, ok := Route("/cities/shanghai", "GET")
	if !ok {
		t.Errorf("Route failed")
	} else {
		handler()
	}
}
