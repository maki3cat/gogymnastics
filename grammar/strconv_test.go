package grammar

import (
	"fmt"
	"strconv"
	"testing"
)

func TestFormat(t *testing.T) {
	s := strconv.FormatBool(true)
	s1 := strconv.FormatFloat(3.1415, 'E', -1, 64)
	s2 := strconv.FormatInt(-42, 16)
	s3 := strconv.FormatInt(-42, 10)
	fmt.Println(s, s1, s2, s3)
	s4 := strconv.FormatFloat(3.1415, 'f', 2, 64)
	fmt.Println(s4)
	s5 := strconv.FormatFloat(3.1415, 'f', 2, 32)
	fmt.Println(s5)
}
