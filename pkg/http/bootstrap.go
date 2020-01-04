package http

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ConvertToConfiguration(arguments map[string]*string) (*ServerConfig, error) {
	port, err := strconv.Atoi(*arguments["port"])
	if err != nil {
		return nil, fmt.Errorf("the given port \"%s\" cannot be parsed as integer", *arguments["port"])
	}
	var jwtKey []byte
	if *arguments["jwtKey"] == "" {
		jwtKey = make([]byte, 255)
		_, err = rand.Read(jwtKey)
		if err != nil {
			panic(err)
		}
	} else {
		jwtKey = []byte(*arguments["jwtKey"])
	}
	return &ServerConfig{
		BindAddress: *arguments["bindAddress"],
		Port:        port,
		AccountProvider: func() ([]Account, error) {
			return readHtPasswd(*arguments["home"] + "/.htpasswd")
		},
		JwtKey: jwtKey,
	}, nil
}

func readHtPasswd(file string) ([]Account, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("could not open passwd file: %v", err)
	}
	//noinspection GoUnhandledErrorResult
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	accounts := make([]Account, 0)
	index := 0
	for scanner.Scan() {
		line := string(scanner.Bytes())
		line = strings.TrimSpace(line)
		index = index + 1
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}
		splitted := strings.Split(line, ":")
		if len(splitted) != 2 {
			return nil, fmt.Errorf("error at line %d: Expected exactly one \":\", got %d", index, len(splitted)-1)
		}
		accounts = append(accounts, Account{
			UserName:          splitted[0],
			EncryptedPassword: splitted[1],
		})
	}
	return accounts, nil
}
