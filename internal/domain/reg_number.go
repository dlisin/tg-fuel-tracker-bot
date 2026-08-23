package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var regNumberRegexp = regexp.MustCompile(`(?i)^[АВЕКМНОРСТУХABEKMHOPCTYX]\d{3}[АВЕКМНОРСТУХABEKMHOPCTYX]{2}\d{2,3}$`)

type RegNumber string

func ParseRegNumber(value string) (RegNumber, error) {
	value = strings.TrimSpace(value)
	value = strings.ToUpper(value)
	value = strings.NewReplacer(
		" ", "",
		"А", "A",
		"В", "B",
		"Е", "E",
		"К", "K",
		"М", "M",
		"Н", "H",
		"О", "O",
		"Р", "P",
		"С", "C",
		"Т", "T",
		"У", "Y",
		"Х", "X",
	).Replace(value)

	if !regNumberRegexp.MatchString(value) {
		return "", fmt.Errorf("invalid registration number: %s", value)
	}

	return RegNumber(value), nil
}

func (r RegNumber) String() string {
	return string(r)
}
