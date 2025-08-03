package advanced

import (
	"fmt"
	"sync"
)

// WHAT WE HAVE
// basic interfaces/apis given to us
type Item struct {
	Title, Channel, GUID string
}

type Fetcher interface {
	Fetch() (items []Item, err error)
}

type fetcher struct {
	domain string
	count  int
}

func (f *fetcher) Fetch() ([]Item, error) {
	items := make([]Item, 5)
	for i := range 5 {
		f.count++
		items[i] = Item{Title: fmt.Sprintf("Hello %d", f.count), Channel: f.domain, GUID: fmt.Sprintf("%d", f.count)}
	}
	return items, nil
}

func Fetch(domain string) Fetcher {
	return &fetcher{domain: domain}
}

// WHAT WE WANT:
// Subscription: convert fetch to a stream
// Merge: merge multiple streams into one
type Subscription interface {
	Updates() <-chan Item // stream of Items
	Close() error         // shuts down the stream
}

// implementation of single subscription
func Subscribe(fetcher Fetcher) Subscription {
	// buffer the updates
	s := &sub{
		fetcher:   fetcher,
		updatesCh: make(chan Item), // for Updates
		closeCh:   make(chan struct{}),
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	fetcher   Fetcher   // fetches items
	updatesCh chan Item // delivers items to the user
	closeCh   chan struct{}
}

func (s *sub) Close() error {
	close(s.closeCh)
	return nil
}

func (s *sub) Updates() <-chan Item {
	return s.updatesCh
}

// the sender of the updates
func (s *sub) loop() {
	// loop is the sender
	defer close(s.updatesCh)
	for {
		items, err := s.fetcher.Fetch()
		if err != nil {
			continue
		}
		for _, item := range items {
			select {
			case <-s.closeCh:
				return
			default:
				select {
				case <-s.closeCh:
					return
				case s.updatesCh <- item:
				}
			}
		}
	}
}

// func Subscribe(fetcher Fetcher) Subscription {
func Merge(subs ...Subscription) Subscription {
	// we use fanout-fanin pattern
	m := &merged{
		subs:    subs,
		faninCh: make(chan Item),
		closeCh: make(chan struct{}),
	}
	go m.loop()
	return m
} // merges several streams

// implementation of merged subscription
type merged struct {
	subs    []Subscription
	faninCh chan Item
	closeCh chan struct{}
}

func (m *merged) Updates() <-chan Item {
	return m.faninCh
}

func (m *merged) Close() error {
	close(m.closeCh)
	// should we propagate the close
	for _, sub := range m.subs {
		sub.Close()
	}
	return nil
}

func (m *merged) loop() {

	// periodically call Fetch
	// send fetched items on the Updates channel
	// exit when Close is called, reporting any error
	defer close(m.faninCh)
	wg := sync.WaitGroup{}
	worker := func(sub Subscription) {
		defer wg.Done()
		for {
			select {
			case <-m.closeCh:
				return
			case item, ok := <-sub.Updates():
				if !ok {
					return
				}
				m.faninCh <- item
			}
		}
	}

	for _, sub := range m.subs {
		wg.Add(1)
		go worker(sub)
	}
	wg.Wait()
}
