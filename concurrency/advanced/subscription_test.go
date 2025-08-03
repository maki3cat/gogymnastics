package advanced

import (
	"fmt"
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	// Subscribe to some feeds, and create a merged update stream.
	merged := Merge(
		Subscribe(Fetch("blog.golang.org")),
		Subscribe(Fetch("googleblog.blogspot.com")),
		Subscribe(Fetch("googledevelopers.blogspot.com")))

	// Close the subscriptions after some time.
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("closed:", merged.Close())
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}
	panic("show me the stacks")
}
