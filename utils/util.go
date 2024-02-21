package utils

import "strings"

func TrimSpacesLR(s string) string {
	// Trim spaces from left
	s = strings.TrimLeft(s, " ")

	// Trim spaces from right
	s = strings.TrimRight(s, " ")

	return s
}
