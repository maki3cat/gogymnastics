package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ## Problem Statement

// From the network layer, it is a common behavior to find the geographical location of a request from its source ip. Assuming you’re designing a service to resolve the source location of a request, and you can build a map between ip subnets and geo locations like this:

// | Subnet      | Location |
// | ----------- | -------- |
// | 15.0.0.0/24 | fr       |
// | 10.0.0.0/16 | cn       |
// | 10.0.1.0/25 | sh       |

// Note that we should always take the one that has more specific matching (sh over cn if possible).

// parse subnet
type Subnet struct {
	ip   []byte
	mask int
}

func parseSubnet(subnet string) *Subnet {
	parts := strings.Split(subnet, `/`)
	s := new(Subnet)
	s.ip = parseIP(parts[0])
	s.mask, _ = strconv.Atoi(parts[1])
	return s
}

// todo: we skip the error handling for now
func parseIP(ip string) []byte {
	parts := strings.Split(ip, `.`)
	res := make([]byte, 4)
	for i, p := range parts {
		val, _ := strconv.Atoi(p)
		res[i] = byte(val)
	}
	return res
}

// register location into trie
type Node struct {
	children map[bool]*Node
	bit      bool
	geo      string
}

func NewNode(bit bool) *Node {
	n := new(Node)
	n.bit = bit
	n.children = make(map[bool]*Node, 2)
	return n
}

func (n *Node) GetChild(bit bool) *Node {
	return n.children[bit]
}

func (n *Node) GetOrAddChild(bit bool) *Node {
	child, ok := n.children[bit]
	if !ok {
		child = NewNode(bit)
		n.children[bit] = child
	}
	return child
}

func getBitFromIP(ip []byte, idx int) bool {
	group := idx / 8
	b := ip[group]
	shift := 7 - (idx - group*8)
	return ((b >> shift) & 1) == 1
}

func RegisterGeo(root *Node, subnet string, loc string) {
	s := parseSubnet(subnet)
	node := root
	for i := range s.mask {
		bit := getBitFromIP(s.ip, i)
		child := node.GetOrAddChild(bit)
		node = child
	}
	node.geo = loc
	fmt.Println("node", node.bit, node.geo)
}

func GetGeo(root *Node, ip string) string {
	ipBytes := parseIP(ip)
	node := root
	res := ""
	for i := range 32 {
		// maki: here is easily mistaken, is we use node.geo; we should put it before child==nil & break
		// if we use child.geo we can put it in the last way
		// I think it is good we handle the current node first and find the child

		// current node
		if node.geo != "" {
			res = node.geo
		}
		bit := getBitFromIP(ipBytes, i)

		// child
		child := node.GetChild(bit)
		if child == nil {
			// should break or return res
			break
		}
		node = child
	}
	return res
}

// | Subnet      | Location |
// | ----------- | -------- |
// | 15.0.0.0/24 | fr       |
// | 10.0.0.0/16 | cn       |
// | 10.0.1.0/25 | sh       |
func main() {
	fmt.Println("main starts")
	dummyRoot := NewNode(true)
	RegisterGeo(dummyRoot, "15.0.0.0/24", "fr")
	RegisterGeo(dummyRoot, "10.0.0.0/16", "cn")
	RegisterGeo(dummyRoot, "10.0.1.0/25", "sh")

	// test-case 1: the original ip
	res := GetGeo(dummyRoot, "15.0.0.0")
	fmt.Println("ip:", "15.0.0.0", "geo:", res)
	res = GetGeo(dummyRoot, "10.0.0.0")
	fmt.Println("ip:", "10.0.0.0", "geo:", res)
	res = GetGeo(dummyRoot, "10.0.1.0")
	fmt.Println("ip:", "10.0.1.0", "geo:", res)

	// test-case 2: normal case exist and not exist
	res = GetGeo(dummyRoot, "15.0.0.168")
	fmt.Println("ip:", "15.0.0.168", "geo:", res)
	res = GetGeo(dummyRoot, "168.0.0.168")
	fmt.Println("ip:", "168.0.0.168", "geo:", res)

	// test-case 3: priority case
	res = GetGeo(dummyRoot, "10.0.168.0")
	fmt.Println("ip:", "10.0.168.0", "geo:", res)
	res = GetGeo(dummyRoot, "10.0.1.168")
	fmt.Println("ip:", "10.0.1.168", "geo:", res)
}
