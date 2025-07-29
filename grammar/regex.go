package grammar

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var logLines = []string{
	"[02/Nov/2018:21:46:31 +0000] PUT /users/12345/locations HTTP/1.1 204 iphone-3",
	"[02/Nov/2018:21:46:31 +0000] PUT /users/6098/locations HTTP/1.1 204 iphone-3",
	"[02/Nov/2018:21:46:32 +0000] PUT /users/3911/locations HTTP/1.1 204 moto-x",
	"[02/Nov/2018:21:46:33 +0000] PUT /users/9933/locations HTTP/1.1 404 moto-x",
	"[02/Nov/2018:21:46:33 +0000] PUT /users/3911/locations HTTP/1.1 500 moto-x",
	"[02/Nov/2018:21:46:34 +0000] GET /rides/9943222/status HTTP/1.1 200 moto-x",
	"[02/Nov/2018:21:46:34 +0000] POST /rides HTTP/1.1 202 iphone-2",
	"[02/Nov/2018:21:46:35 +0000] POST /users HTTP/1.1 202 iphone-5",
	"[02/Nov/2018:21:46:35 +0000] POST /rides HTTP/1.1 202 iphone-5",
	"[02/Nov/2018:21:46:37 +0000] POST /rides HTTP/1.1 202 iphone-4",
	"[02/Nov/2018:21:46:38 +0000] GET /users/994/ride/16 HTTP/1.1 200 iphone-5",
	"[02/Nov/2018:21:46:39 +0000] POST /users HTTP/1.1 202 iphone-3",
	"[02/Nov/2018:21:46:40 +0000] PUT /users/8384721/locations HTTP/1.1 204 iphone-3",
	"[02/Nov/2018:21:46:41 +0000] GET /users/342111 HTTP/1.1 200 iphone-5",
	"[02/Nov/2018:21:46:42 +0000] GET /users/9933 HTTP/1.1 200 iphone-5",
	"[02/Nov/2018:21:46:43 +0000] GET /prices/20180103/geo/12 HTTP/1.1 200 iphone-5",
}

// regex matched urls:
// /users/[0-9]+/locations -> /users/#/locations
// /rides/[0-9]+/status -> /rides/#/status
// /users/[0-9]+/ride/[0-9]+ -> /users/#/ride/#
// /prices/[0-9]+/geo/[0-9]+ -> /prices/#/geo/#
// /users/[0-9]+ -> /users/#
// others, keep the original
var (
	usersLocationsRE *regexp.Regexp
	ridesStatusRE    *regexp.Regexp
	pricesGeoRE      *regexp.Regexp
	usersRideRE      *regexp.Regexp
	usersRE          *regexp.Regexp
)

func init() {
	usersLocationsRE = regexp.MustCompile(`/users/[0-9]+/locations`)
	ridesStatusRE = regexp.MustCompile(`/rides/[0-9]+/status`)
	pricesGeoRE = regexp.MustCompile(`/prices/[0-9]+/geo/[0-9]+`)
	usersRideRE = regexp.MustCompile(`/users/[0-9]+/ride/[0-9]+`)
	usersRE = regexp.MustCompile(`/users/[0-9]+`)
}

func extractPath(rawPath string) string {
	if usersLocationsRE.MatchString(rawPath) {
		return "/users/#/locations"
	}
	if ridesStatusRE.MatchString(rawPath) {
		return "/rides/#/status"
	}
	if pricesGeoRE.MatchString(rawPath) {
		return "/prices/#/geo/#"
	}
	if usersRideRE.MatchString(rawPath) {
		return "/users/#/ride/#"
	}
	if usersRE.MatchString(rawPath) {
		return "/users/#"
	}
	return rawPath
}

type Event struct {
	Method   string
	Endpoint string
	Code     int
	Count    int
}

func (e Event) line(isHeader bool) string {
	// Define column widths for alignment
	const (
		methodWidth   = 6
		endpointWidth = 25
		codeWidth     = 6
		countWidth    = 6
	)
	if isHeader {
		return fmt.Sprintf("|| %-*s | %-*s | %*s | %*s ||",
			methodWidth, "Method",
			endpointWidth, "Endpoint",
			codeWidth, "Code",
			countWidth, "Count")
	}
	return fmt.Sprintf("|| %-*s | %-*s | %*d | %*d ||",
		methodWidth, e.Method,
		endpointWidth, e.Endpoint,
		codeWidth, e.Code,
		countWidth, e.Count)
}

func scanLines() []string {
	return logLines
}

func count(logLines []string) map[Event]int {
	counts := make(map[Event]int)
	startIdx := 28
	for _, line := range logLines {
		// split the time by length
		// _ := line[:15]
		// split the rest by space
		parts := strings.Split(strings.TrimSpace(line[startIdx:]), " ")
		code, err := strconv.Atoi(strings.TrimSpace(parts[3]))
		if err != nil {
			fmt.Println("error converting code to int", err, parts[3])
			continue
		}
		counts[Event{
			Method:   strings.TrimSpace(parts[0]),
			Endpoint: extractPath(strings.TrimSpace(parts[1])),
			Code:     code,
		}] += 1
	}
	return counts
}

func sortEvents(counts map[Event]int) []Event {
	events := make([]Event, 0, len(counts))
	for event, count := range counts {
		event.Count = count
		events = append(events, event)
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].Count > events[j].Count
	})
	return events
}

func PrintLogs() {
	counts := count(scanLines())
	events := sortEvents(counts)
	fmt.Println(Event{}.line(true))
	for _, event := range events {
		fmt.Println(event.line(false))
	}
}
