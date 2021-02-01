package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var b strings.Builder
	if input == "" {
		return "", nil
	}
	r := []rune(input)
	if unicode.IsDigit(r[0]) {
		return "", ErrInvalidString
	}
	var prev rune
	for _, symbol := range r {
		switch {
		case unicode.IsDigit(symbol) && unicode.IsDigit(prev):
			{
				return "", ErrInvalidString
			}
		case unicode.IsDigit(symbol) && !unicode.IsDigit(prev):
			{
				if num, err := strconv.Atoi(string(symbol)); err == nil {
					_, err = b.WriteString(strings.Repeat(string(prev), num))
					if err != nil {
						return "", err
					}
				}
			}
		default:
			if unicode.IsLetter(prev) {
				b.WriteRune(prev)
			}
		}
		prev = symbol
	}
	if unicode.IsLetter(prev) {
		if _, err := b.WriteRune(prev); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}
