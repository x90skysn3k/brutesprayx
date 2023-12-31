package modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func ReadUsersFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	users := []string{}
	for scanner.Scan() {
		users = append(users, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func ReadPasswordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	passwords := []string{}
	for scanner.Scan() {
		passwords = append(passwords, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return passwords, nil
}

func GetUsersFromDefaultWordlist() []string {
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

func GetPasswordsFromDefaultWordlist() []string {
	wordlistPath := filepath.Join("wordlist", service, "pass")

	f, err := os.Open(wordlistPath)
	if err != nil {
		fmt.Printf("Error opening %s wordlist: %s\n", "pass", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	users := []string{}
	for scanner.Scan() {
		users = append(users, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s wordlist: %s\n", "pass", err)
		os.Exit(1)
	}

	return users
}
