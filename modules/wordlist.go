package modules

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v29/github"
	"github.com/pterm/pterm"
)

func downloadFileFromGithub(repoOwner, repoName, filePath, localPath string) error {
	client := github.NewClient(nil)
	content, _, _, err := client.Repositories.GetContents(context.Background(), repoOwner, repoName, filePath, nil)
	if err != nil {
		return err
	}

	data, err := content.GetContent()
	if err != nil {
		return err
	}

	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	bar := pterm.DefaultProgressbar.WithTotal(len(data))
	bar.Start()
	for _, c := range data {
		_, err := file.WriteString(string(c))
		if err != nil {
			return err
		}
		bar.Increment()
	}
	bar.Stop()

	return nil
}

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
	wordlistPath := filepath.Join("wordlist", "user")

	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		err := downloadFileFromGithub("x90skysn3k", "brutesprayx", "wordlist/user", wordlistPath)
		if err != nil {
			fmt.Printf("Error downloading user wordlist: %s\n", err)
			os.Exit(1)
		}
	}

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
	wordlistPath := filepath.Join("wordlist", "password")

	if _, err := os.Stat(wordlistPath); os.IsNotExist(err) {
		err := downloadFileFromGithub("x90skysn3k", "brutesprayx", "wordlist/password", wordlistPath)
		if err != nil {
			fmt.Printf("Error downloading password wordlist: %s\n", err)
			os.Exit(1)
		}
	}

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
