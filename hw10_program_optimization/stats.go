package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//go:generate easyjson -all stats.go
//nolint:all
type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	ds, err := getDomainStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return ds, nil
}

func getDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	domainSuffix := "." + strings.ToLower(domain)

	var user User

	for scanner.Scan() {
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			continue
		}

		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, domainSuffix) {
			parts := strings.SplitN(email, "@", 2)
			result[parts[1]]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
