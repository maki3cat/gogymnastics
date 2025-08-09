package subnet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSubnet(t *testing.T) {
	tests := []struct {
		subnet string
		ip     []byte
		mask   int
	}{
		{"10.0.0.0/25", []byte{10, 0, 0, 0}, 25},
		{"15.0.0.0/24", []byte{15, 0, 0, 0}, 24},
		{"10.1.0.0/16", []byte{10, 1, 0, 0}, 16},
	}
	for _, test := range tests {
		ip, mask, err := parseSubnet(test.subnet)
		if err != nil {
			t.Errorf("parseSubnet(%s) = %v", test.subnet, err)
		}
		if !reflect.DeepEqual(ip, test.ip) {
			t.Errorf("parseSubnet(%s) = %v, want %v", test.subnet, ip, test.ip)
		}
		if mask != test.mask {
			t.Errorf("parseSubnet(%s) = %d, want %d", test.subnet, mask, test.mask)
		}
	}
}

func TestHelperGetBit(t *testing.T) {
	tests := []struct {
		ip   []byte
		pos  int
		want bool
	}{
		{[]byte{10, 0, 0, 0}, 0, false},  // 10 = 00001010, bit 0 is 0
		{[]byte{10, 0, 0, 0}, 1, false},  // 10 = 00001010, bit 1 is 1
		{[]byte{10, 0, 0, 0}, 2, false},  // 10 = 00001010, bit 2 is 0
		{[]byte{10, 0, 0, 0}, 3, false},  // 10 = 00001010, bit 3 is 0
		{[]byte{10, 0, 0, 0}, 4, true},   // 10 = 00001010, bit 4 is 1
		{[]byte{10, 0, 0, 0}, 5, false},  // 10 = 00001010, bit 5 is 1
		{[]byte{10, 0, 0, 0}, 6, true},   // 10 = 00001010, bit 5 is 1
		{[]byte{10, 0, 0, 0}, 7, false},  // 10 = 00001010, bit 5 is 1
		{[]byte{10, 0, 0, 0}, 17, false}, // 10 = 00001010, bit 5 is 1
	}
	for _, test := range tests {
		got := helperGetBit(test.ip, test.pos)
		if got != test.want {
			t.Errorf("helperGetBit(%v, %d) = %v, want %v", test.ip, test.pos, got, test.want)
		}
	}

}

func TestBuildTrie(t *testing.T) {
	tests := []struct {
		subnets []byte
		maskLen int
		value   string
	}{
		{[]byte{10, 0, 0, 0}, 25, "US"},
	}
	root := NewBitNode(false)
	for _, test := range tests {
		BuildTrie(root, test.subnets, test.maskLen, test.value)
	}
	fmt.Println("the root has children", len(root.children), root.children)
	assert.Equal(t, "US", SearchTrie(root, []byte{10, 0, 0, 0}))
	assert.Equal(t, "", SearchTrie(root, []byte{15, 0, 0, 0}))
}

func TestSetGeoAndGetGeo(t *testing.T) {
	// Reset the global root for clean testing
	root = NewBitNode(false)

	// Test cases for setting and getting geo locations
	tests := []struct {
		subnet string
		geo    string
	}{
		{"10.0.0.0/16", "cn"},
		{"10.0.1.0/24", "sh"},
		{"15.0.0.0/24", "fr"},
		{"192.168.1.0/28", "us"},
		{"172.16.0.0/12", "de"},
	}

	// Set geo locations
	for _, test := range tests {
		err := SetGeo(test.subnet, test.geo)
		if err != nil {
			t.Errorf("SetGeo(%s, %s) failed: %v", test.subnet, test.geo, err)
		}
	}

	// Test IP lookups
	ipTests := []struct {
		ip       string
		expected string
	}{
		{"10.0.0.1", "cn"},     // matches 10.0.0.0/16
		{"10.0.0.127", "cn"},   // matches 10.0.0.0/16
		{"10.0.1.128", "sh"},   // matches 10.1.0.0/16 (more specific than /25)
		{"15.0.0.255", "fr"},   // matches 15.0.0.0/24
		{"192.168.1.5", "us"},  // matches 192.168.1.0/28
		{"172.16.255.1", "de"}, // matches 172.16.0.0/12
		{"8.8.8.8", ""},        // no match
		{"1.1.1.1", ""},        // no match
	}

	for _, test := range ipTests {
		got := GetGeo(test.ip)
		if got != test.expected {
			t.Errorf("GetGeo(%s) = %s, want %s", test.ip, got, test.expected)
		}
	}
}

func TestSetGeoInvalidSubnet(t *testing.T) {
	tests := []struct {
		subnet string
		geo    string
	}{
		{"invalid", "us"},
		{"10.0.0.0/", "us"},
		{"10.0.0.0/abc", "us"},
		{"10.0.0", "us"},
		{"10.0.0.0.1/24", "us"},
		{"256.0.0.0/24", "us"},
	}

	for _, test := range tests {
		err := SetGeo(test.subnet, test.geo)
		if err == nil {
			t.Errorf("SetGeo(%s, %s) should have failed but didn't", test.subnet, test.geo)
		}
	}
}

func TestGetGeoInvalidIP(t *testing.T) {
	tests := []string{
		"invalid",
		"10.0.0",
		"10.0.0.0.1",
		"256.0.0.0",
		"10.0.0.abc",
	}

	for _, test := range tests {
		got := GetGeo(test)
		if got != "" {
			t.Errorf("GetGeo(%s) = %s, want empty string for invalid IP", test, got)
		}
	}
}

func TestOverlappingSubnets(t *testing.T) {
	// Reset the global root for clean testing
	root = NewBitNode(false)

	// Set up overlapping subnets - more specific should take precedence
	err := SetGeo("10.0.0.0/16", "broad")
	if err != nil {
		t.Fatalf("SetGeo failed: %v", err)
	}

	err = SetGeo("10.0.1.0/24", "specific")
	if err != nil {
		t.Fatalf("SetGeo failed: %v", err)
	}

	// Test that more specific subnet takes precedence
	tests := []struct {
		ip       string
		expected string
	}{
		{"10.0.0.1", "broad"},    // matches broader /16
		{"10.0.1.1", "specific"}, // matches more specific /24
		{"10.0.2.1", "broad"},    // matches broader /16
	}

	for _, test := range tests {
		got := GetGeo(test.ip)
		if got != test.expected {
			t.Errorf("GetGeo(%s) = %s, want %s", test.ip, got, test.expected)
		}
	}
}
