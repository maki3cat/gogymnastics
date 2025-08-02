package count

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const url_count = 16

var userPath *regexp.Regexp
var userLocPath *regexp.Regexp
var userRidePath *regexp.Regexp
var rideStatusPath *regexp.Regexp
var pricePath *regexp.Regexp

func init() {
	userPath = regexp.MustCompile("^/users/[0-9]+$")
	userLocPath = regexp.MustCompile("^/users/[0-9]+/locations$")
	userRidePath = regexp.MustCompile("^/users/[0-9]+/ride/[0-9]+$")
	rideStatusPath = regexp.MustCompile("^/rides/[0-9]+/status$")
	pricePath = regexp.MustCompile("^/prices/[0-9]+/geo/[0-9]+$")
}

// tests
func parsePath(rawPath string) string {
	if userPath.MatchString(rawPath) {
		return "/users/#"
	}
	if userLocPath.MatchString(rawPath) {
		return "/users/#/locations"
	}
	if userRidePath.MatchString(rawPath) {
		return "/users/#/ride/#"
	}
	if rideStatusPath.MatchString(rawPath) {
		return "/rides/#/status"
	}
	if pricePath.MatchString(rawPath) {
		return "/prices/#/geo/#"
	}
	return rawPath
}

type ReqResp struct {
	Method string
	Path   string
	Status int
}

type Stat struct {
	ReqResp
	Count int // for sorting
}

func (r *Stat) String() string {
	return fmt.Sprintf("%-6s | %-20s | %-3d | %-4d", r.Method, r.Path, r.Status, r.Count)
}

func scan(logline string) *ReqResp {
	n := len("[02/Nov/2018:21:46:43 +0000] ")
	// split the line wihtout the time: n: end of it
	// separator is " "
	logInfos := strings.Split(logline[n:], " ")
	if len(logInfos) < 4 {
		fmt.Println("corrupt line", logline)
		return nil
	}
	status, err := strconv.Atoi(logInfos[3])
	if err != nil {
		fmt.Println("corrupt line", logline)
		return nil
	}
	res := ReqResp{
		Method: logInfos[0],
		Path:   parsePath(logInfos[1]),
		Status: status,
	}
	return &res
}

// :returns: sorted slice of reqRep based on count
func countAndsortReqResp(loglines []string) []*Stat {
	counts := make(map[ReqResp]int, url_count)
	for _, line := range loglines {
		reqResp := *scan(line)
		counts[reqResp] = counts[reqResp] + 1 // defaults 0
	}
	res := make([]*Stat, 0, len(counts))
	for k, v := range counts {
		val := &Stat{
			ReqResp: k,
			Count:   v,
		}
		res = append(res, val)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Count > res[j].Count
	})
	return res
}

// :returns: print
func print(stats []*Stat) {
	for _, stat := range stats {
		fmt.Println(stat.String())
	}
}

func ProcessLogsAndPrint(loglines []string) {
	stats := countAndsortReqResp(loglines)
	print(stats)
}
