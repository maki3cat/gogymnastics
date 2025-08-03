package grammar

import "regexp"

var digitRegexp = regexp.MustCompile(`^\d+$`)

func isDigit(s string) bool {
	return digitRegexp.MatchString(s)
}
