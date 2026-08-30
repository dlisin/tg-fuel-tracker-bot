package util

import (
	"strconv"
	"strings"
)

func ParseInt64(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseInt(s, 10, 64)
}

func ParseFloat64(s string) (float64, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}
