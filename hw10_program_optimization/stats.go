package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

var (
	ErrNilReader       = errors.New("nil reader received")
	ErrEmptySourceData = errors.New("empty data source received")
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, ErrNilReader
	}

	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	i := 0

	for scanner.Scan() {
		line := scanner.Bytes()

		var user User
		if err = json.Unmarshal(line, &user); err != nil {
			return
		}
		result[i] = user
		i++
	}
	if i == 0 {
		err = ErrEmptySourceData
		return
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := re.FindString(user.Email)

		if matched != "" {
			s := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[s]++
		}
	}

	return result, nil
}
