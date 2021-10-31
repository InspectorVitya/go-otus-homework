package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, tld string) (DomainStat, error) {
	result := DomainStat{}
	if tld == "" {
		return result, nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Bytes()

		email := jsoniter.Get(line, "Email").ToString()

		if strings.HasSuffix(email, "."+tld) {
			num := result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])] = num
		}
	}

	return result, nil
}
