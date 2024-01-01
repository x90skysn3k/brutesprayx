package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
)

func getResultString(result bool) string {
	if result {
		return "succeeded"
	} else {
		return "failed"
	}
}

func WriteToFile(filename string, content string) error {
	timestamp := time.Now().Format("2006010215")
	dir := "output"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	filename = filepath.Join(dir, filename+"_"+timestamp)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func PrintResult(service string, host string, port int, user string, pass string, result bool) {

	if result {
		pterm.Success.Println("Attempt", service, "SUCCESS on host", host, "port", port, "with username", user, "and password", pass, getResultString(result))
		content := fmt.Sprintf("Attempt %s SUCCESS on host %s port %d with username %s and password %s %s\n", service, host, port, user, pass, getResultString(result))
		filename := filepath.Base(host)
		WriteToFile(filename, content)
	}

	pterm.Color(pterm.FgLightRed).Println("Attempt", service, "on host", host, "port", port, "with username", user, "and password", pass, getResultString(result))

}
