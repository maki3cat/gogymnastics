package practice

import (
	"strconv"
	"strings"
)

// From the network layer, it is a common behavior to find the geographical location of a request from its source ip. Assuming you’re designing a service to resolve the source location of a request, and you can build a map between ip subnets and geo locations like this:
//
// | Subnet      | Location |
// | ----------- | -------- |
// | 15.0.0.0/24 | fr       |
// | 10.0.0.0/16 | cn       |
// | 10.0.1.0/25 | sh       |
//
// Note that we should always take the one that has more specific matching (sh over cn if possible).
//
// **API**
// ```kotlin
// class Locator {
//     fun setGeoLocation(subnet: String, geo: String)
//     fun getGeoLocation(ip: String): String
// }
// ```

// todo: we skip the error handling for now

// ---parse the string to ip and mask----
func parseSubnet(subnet string) ([]byte, int) {
	parts := strings.Split(subnet, `/`)
	mask, _ := strconv.Atoi(parts[1])
	ip := parseIP(parts[0])
	return ip, mask
}

func parseIP(ip string) []byte {
	parts := strings.Split(ip, `.`)
	res := make([]byte, 4)
	for i := range 4 {
		val, _ := strconv.Atoi(parts[i])
		res[i] = byte(val)
	}
	return res
}

// maki:
// taking care if we use for range mask, we are already using the index
// get the bit value of the ip, counting from left to right
func getBit(ip []byte, bitIdx int) bool {
	// group
	// bitIdx := bitNum - 1
	idx := bitIdx / 8
	// shift
	left := bitIdx - idx*8
	part := ip[idx]
	rightShift := 7 - left
	// maki: taking the bit is use right shift & number 1 and compare with 1
	return (part>>rightShift)&1 == 1
}

// func mergeBit(ip []byte) int32 {
// 	var val int32 = 0
// 	for idx, part := range ip {
// 		leftShift := 8 * (3 - idx)
// 		val = val | (int32(part) << leftShift)
// 	}
// 	return val
// }

// func getBitV2(ip int32, bitNumber int) bool {
// 	bitIdx := bitNumber - 1
// 	rightShift := 31 - bitIdx
// 	return ip>>int32(rightShift)&1 == 1
// }

// --------build the trie tree---
type Node struct {
	bitVal   bool
	children map[bool]*Node
	Geo      string
}

func NewNode(bitVal bool) *Node {
	return &Node{
		bitVal:   bitVal,
		children: make(map[bool]*Node, 2),
		Geo:      "",
	}
}

func (n *Node) FindChild(bitVal bool) *Node {
	return n.children[bitVal]
}

func (n *Node) AddChild(bitVal bool) *Node {
	node := NewNode(bitVal)
	n.children[bitVal] = node
	return node
}

var root *Node = NewNode(true) // maki: dummy node, the value not used

func AddGeo(subnet string, geo string) {
	ip, mask := parseSubnet(subnet)
	node := root
	for i := range mask {
		bitVal := getBit(ip, i)
		child := node.FindChild(bitVal)
		if child == nil {
			child = node.AddChild(bitVal)
		}
		node = child
	}
	node.Geo = geo
}

func GetGeo(ip string) string {
	ipInBits := parseIP(ip)
	geo := ""
	node := root
	for i := range 32 {
		bitVal := getBit(ipInBits, i)
		child := node.FindChild(bitVal)
		// maki: stop when the child cannot be found
		if child == nil {
			break
		}
		// maki: replace only when the value is not empty
		if child.Geo != "" {
			geo = child.Geo
		}
		node = child
	}
	return geo
}

//     fun getGeoLocation(ip: String): String

// | Subnet      | Location |
// | ----------- | -------- |
// | 15.0.0.0/24 | fr       |
// | 10.0.0.0/16 | cn       |
// | 10.0.1.0/25 | sh       |
