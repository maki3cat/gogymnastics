package grammar

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	s := "hello world get host"
	parts := strings.Split(s, " ")
	fmt.Println(parts)
}

func TestIsDigit(t *testing.T) {
	assert.True(t, isDigit("123"))
	assert.False(t, isDigit("123a"))
}

func TestStringConversions(t *testing.T) {
	s := strconv.FormatBool(true)
	fmt.Println(s)
	s = strconv.FormatFloat(3.1415, 'E', -1, 64)
	fmt.Println(s)

	// format to hex
	s = strconv.FormatInt(-42, 16)
	fmt.Println(s)
	s = strconv.FormatUint(42, 16)
	fmt.Println(s)

	// test
	s = strconv.Itoa(42)
	fmt.Println(s)

	// itoa, atoi
	i, err := strconv.Atoi("32464542")
	assert.NoError(t, err)
	fmt.Println(i)
	i, err = strconv.Atoi("user")
	assert.Error(t, err)
	fmt.Println(i)
}

func TestTypeSafeAssertion(t *testing.T) {
	// maki: casting an interface to a specific type,
	// use the ok expression to check
	// maki: also for errors
	var x any
	x = "hello"
	if s, ok := x.(string); ok {
		fmt.Println("string:", s)
	} else {
		fmt.Println("not a string")
	}
}

func TestTypeConversion(t *testing.T) {
	x := []byte{'a', 'b', 'c', 'A', 'B', 'C'}
	s := string(x)
	fmt.Println(s, x)

	// convert to string
	s = string(x) // this is a no-op
	fmt.Println(s, x)
}
