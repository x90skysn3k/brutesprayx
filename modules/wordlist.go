package modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func GetUsersFromWordlist(service string) []string {
	wordlistPath := filepath.Join("wordlist", service, "user")

	f, err := os.Open(wordlistPath)
	if err != nil {
		fmt.Printf("Error opening %s wordlist: %s\n", "user", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	users := []string{}
	for scanner.Scan() {
		users = append(users, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s wordlist: %s\n", "user", err)
		os.Exit(1)
	}

	return users
}

func ReadFromWordlist(file string, listType string) string {
	wordlistPath := filepath.Join("wordlist", file, listType)

	f, err := os.Open(wordlistPath)
	if err != nil {
		fmt.Printf("Error opening %s wordlist: %s\n", listType, err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		fmt.Printf("Error reading %s wordlist: %s\n", listType, scanner.Err())
		os.Exit(1)
	}

	return scanner.Text()
}
