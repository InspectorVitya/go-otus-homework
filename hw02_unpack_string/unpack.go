package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packed string) (string, error) {
	var prev rune
	var builder strings.Builder
	for _, sym := range packed {
		switch {
		case unicode.IsDigit(sym) && isLetterOrSpace(prev):
			{
				if val, err := strconv.Atoi(string(sym)); err == nil && val > 0 {
					_, err = builder.WriteString(strings.Repeat(string(prev), val))
					if err != nil {
						return "", err
					}
				}
			}
		case unicode.IsDigit(sym) && !isLetterOrSpace(prev):
			return "", ErrInvalidString
		default:
			if isLetterOrSpace(prev) {
				if _, err := builder.WriteRune(prev); err != nil {
					return "", err
				}
			}
		}
		prev = sym
	}
	if isLetterOrSpace(prev) {
		if _, err := builder.WriteRune(prev); err != nil {
			return "", err
		}
	}
	return builder.String(), nil
}

func isLetterOrSpace(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsSpace(r)
}
