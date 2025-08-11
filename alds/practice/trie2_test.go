package practice

import "testing"

func TestSetGeoAndGetGeo(t *testing.T) {
	// Reset the global root for clean testing

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
		AddGeo(test.subnet, test.geo)
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
