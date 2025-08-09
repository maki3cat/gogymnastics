package trie

import (
	"fmt"
	"strconv"
	"strings"
)

//From the network layer, it is a common behavior to find the geographical location of a request from its source ip. Assuming you’re designing a service to resolve the source location of a request, and you can build a map between ip subnets and geo locations like this:
// 10.0.0.0/25 → cn
// 15.0.0.0/24 → fr
// 10.1.0.0/16 → sh
//here’s a sample definition of the class
//class Locator {
// fun setGeoLocation(subnet: String, geo: String)
// fun getGeoLocation(ip: String): String
//}

// assumes the subnet is valid
// @returns: the 4 bytes of the ip address, and the length of the mask
func parseSubnet(subnet string) ([]byte, int, error) {
	parts := strings.Split(subnet, "/")
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("invalid subnet")
	}
	mask, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, 0, fmt.Errorf("invalid subnet: %w", err)
	}
	ipBytes, err := parseIp(parts[0])
	if err != nil {
		return nil, 0, fmt.Errorf("invalid subnet: %w", err)
	}
	return ipBytes, mask, nil
}

func parseIp(ipStr string) ([]byte, error) {
	ip := strings.Split(ipStr, ".")
	if len(ip) != 4 {
		return nil, fmt.Errorf("invalid ip")
	}
	ipBytes := make([]byte, 4)
	for i, part := range ip {
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid ip: %w", err)
		}
		ipBytes[i] = byte(val)
	}
	return ipBytes, nil
}

// @returns: the bit at the position
func helperGetBit(ip []byte, pos int) bool {
	// from left to right
	idx := pos / 8
	bit := pos % 8
	return (ip[idx] & (1 << (7 - bit))) != 0
}

// trie related
type BitNode struct {
	value    string
	bit      bool // one bit 0-1
	children map[bool]*BitNode
}

func NewBitNode(bit bool) *BitNode {
	return &BitNode{
		value:    "",
		bit:      bit,
		children: make(map[bool]*BitNode),
	}
}

func (n *BitNode) AddChild(bit bool) *BitNode {
	if _, ok := n.children[bit]; !ok {
		fmt.Println("add child", bit)
		n.children[bit] = NewBitNode(bit)
	}
	return n.children[bit]
}

func (n *BitNode) GetChild(bit bool) *BitNode {
	return n.children[bit]
}

// root is dummy, we don't use it
func BuildTrie(root *BitNode, subnets []byte, maskLen int, value string) {
	node := root
	for i := range maskLen {
		bit := helperGetBit(subnets, i)
		child := node.GetChild(bit)
		if child == nil {
			child = node.AddChild(bit)
		}
		node = child
	}
	node.value = value
}

// root is dummy, we don't use it
func SearchTrie(root *BitNode, ip []byte) string {
	node := root
	for i := range len(ip) {
		bit := helperGetBit(ip, i)
		fmt.Println(i, bit)
		child := node.GetChild(bit)
		if child == nil {
			return ""
		}
		node = child
	}
	return node.value
}

// The functions
// fun setGeoLocation(subnet: String, geo: String)
// fun getGeoLocation(ip: String): String
var root *BitNode = NewBitNode(false)

func SetGeo(subnet string, geo string) error {
	ip, mask, err := parseSubnet(subnet)
	if err != nil {
		return fmt.Errorf("invalid subnet: %w", err)
	}
	BuildTrie(root, ip, mask, geo)
	return nil
}

func GetGeo(ip string) string {
	ipBytes, _, err := parseSubnet(ip)
	if err != nil {
		return ""
	}
	return SearchTrie(root, ipBytes)
}
