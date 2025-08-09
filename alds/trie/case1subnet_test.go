package trie

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
	for _, test := range tests {
		root := NewBitNode(false)
		BuildTrie(root, test.subnets, test.maskLen, test.value)
	}
	fmt.Println("the root has children", len(root.children), root.children)
	assert.Equal(t, "US", SearchTrie(root, []byte{10, 0, 0, 0}))
	// assert.Equal(t, "UK", SearchTrie(root, []byte{15, 0, 0, 0}))
}
