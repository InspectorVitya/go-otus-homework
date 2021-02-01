package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	r := []rune(input)
	if unicode.IsDigit(r[0]) {
		return "", ErrInvalidString
	}
	var b strings.Builder
	var symbol rune
	var symbolNext rune
	strLen := len(r)
	for i := 0; i < strLen; i++ {
		symbol = r[i]
		if strLen <= i+1 {
			if !unicode.IsDigit(symbol) && symbol != 92 {
				b.WriteRune(symbol)
				break
			} else {
				break
			}
		} else {
			symbolNext = r[i+1]
		}
		if unicode.IsDigit(symbol) && unicode.IsDigit(symbolNext) {
			return "", ErrInvalidString
		}
		switch {
		case symbol == 92:
			{
				if !unicode.IsDigit(symbolNext) && symbolNext != 92 {
					return "", ErrInvalidString
				}
				if strLen > i+2 {
					if unicode.IsDigit(r[i+2]) {
						num, _ := strconv.Atoi(string(r[i+2]))
						b.WriteString(strings.Repeat(string(symbolNext), num))
						i++
						continue
					} else {
						b.WriteRune(symbolNext)
						i++
						continue
					}
				} else {
					b.WriteRune(symbolNext)
					break
				}
			}
		case unicode.IsDigit(symbolNext):
			{
				num, _ := strconv.Atoi(string(symbolNext))
				b.WriteString(strings.Repeat(string(symbol), num))
				continue
			}
		}
		if !unicode.IsDigit(symbol) && symbol != 92 {
			b.WriteRune(symbol)
		}
	}
	return b.String(), nil
}
