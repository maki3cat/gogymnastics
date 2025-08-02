package count

import (
	"testing"
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

func TestCount(t *testing.T) {
	ProcessLogsAndPrint(logLines)
	// fmt.Println(stats)
}

func TestCount2(t *testing.T) {
	Process(logLines)
}

// Method |             Endpoint | Code || Count
// =============================================
//  PUT   |   /users/#/locations | 204  ||  4
// POST   |               /rides | 202  ||  3
//  GET   |             /users/# | 200  ||  2
// POST   |               /users | 202  ||  2
//  PUT   |   /users/#/locations | 500  ||  1
//  GET   |      /prices/#/geo/# | 200  ||  1
//  PUT   |   /users/#/locations | 404  ||  1
//  GET   |      /rides/#/status | 200  ||  1
//  GET   |      /users/#/ride/# | 200  ||  1
