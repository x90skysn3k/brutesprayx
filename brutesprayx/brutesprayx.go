package brutesprayx

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"github.com/x90skysn3k/brutesprayx/modules"
	"github.com/x90skysn3k/brutesprayx/parse"
)

var NAME_MAP = map[string]string{
	"ms-sql-s":       "mssql",
	"microsoft-ds":   "smbnt",
	"cifs":           "smbnt",
	"postgresql":     "postgres",
	"smtps":          "smtp",
	"submission":     "smtp",
	"imaps":          "imap",
	"pop3s":          "pop3",
	"iss-realsecure": "vmauthd",
	"snmptrap":       "snmp",
	//"ms-wbt-server":  "rdp",
}

func mapService(service string) string {
	if mappedService, ok := NAME_MAP[service]; ok {
		return mappedService
	}
	return service
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func parseFile(filename string) (map[parse.Host]int, error) {
	in_format := ""
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	line := scanner.Text()

	if line[0] == '{' {
		in_format = "json"
	} else if strings.HasPrefix(line, "# Nmap") {
		if !scanner.Scan() {
			return nil, scanner.Err()
		}
		line = scanner.Text()
		if !strings.HasPrefix(line[1:], "Nmap") {
			in_format = "gnmap"
		}
	} else if strings.HasPrefix(line, "<NexposeReport ") {
		in_format = "xml_nexpose"
	} else if strings.Contains(line, "<?xml ") {
		if !scanner.Scan() {
			return nil, scanner.Err()
		}
		line = scanner.Text()
		if strings.Contains(line, "nmaprun") {
			in_format = "xml"

		} else if strings.HasPrefix(line, "<NessusClientData") {
			in_format = "xml_nessus"
		}
	} else {
		in_format = "list"
	}

	if in_format == "" {
		fmt.Println("File is not correct format!")
		os.Exit(0)
	}

	switch in_format {
	case "gnmap":
		hosts, err := parse.ParseGNMAP(filename)
		return hosts, err
	case "json":
		hosts, err := parse.ParseJSON(filename)
		return hosts, err
	case "xml":
		hosts, err := parse.ParseXML(filename)
		return hosts, err
	case "xml_nexpose":
		hosts, err := parse.ParseNexpose(filename)
		return hosts, err
	case "xml_nessus":
		hosts, err := parse.ParseNessus(filename)
		return hosts, err
	case "list":
		hosts, err := parse.ParseList(filename)
		return hosts, err
	default:
		return nil, fmt.Errorf("unsupported file type: %s", in_format)
	}
}

func readUsersFromFile(filename string) ([]string, error) {
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

func readPasswordsFromFile(filename string) ([]string, error) {
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

func writeToFile(filename string, content string) error {
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

func countHosts(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func isFile(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil && filepath.Ext(fileName) == "" {
		return true
	}
	return false
}

func brute(h parse.Host, u string, p string) {
	service := mapService(h.Service)
	var result bool

	switch service {
	case "ssh":
		result = modules.BruteSSH(h.Host, h.Port, u, p)
	case "ftp":
		result = modules.BruteFTP(h.Host, h.Port, u, p)
	case "mssql":
		result = modules.BruteMSSQL(h.Host, h.Port, u, p)
	case "telnet":
		result = modules.BruteTelnet(h.Host, h.Port, u, p)
	case "smbnt":
		result = modules.BruteSMB(h.Host, h.Port, u, p)
	case "postgres":
		result = modules.BrutePostgres(h.Host, h.Port, u, p)
	case "smtp":
		result = modules.BruteSMTP(h.Host, h.Port, u, p)
	case "imap":
		result = modules.BruteIMAP(h.Host, h.Port, u, p)
	case "pop3":
		result = modules.BrutePOP3(h.Host, h.Port, u, p)
	case "snmp":
		result = modules.BrutePOP3(h.Host, h.Port, u, p)
	//case "rdp":
	//	result = modules.BruteRDP(h.Host, h.Port, u, p)
	default:
		//fmt.Printf("Unsupported service: %s\n", h.Service)
		return
	}

	printResult(service, h.Host, h.Port, u, p, result)
}

func printResult(service string, host string, port int, user string, pass string, result bool) {

	if result {
		pterm.Success.Println("Attempt", service, "SUCCESS on host", host, "port", port, "with username", user, "and password", pass, getResultString(result))
		content := fmt.Sprintf("Attempt %s SUCCESS on host %s port %d with username %s and password %s %s\n", service, host, port, user, pass, getResultString(result))
		filename := filepath.Base(host)
		writeToFile(filename, content)
	}

	pterm.Color(pterm.FgLightRed).Println("Attempt", service, "on host", host, "port", port, "with username", user, "and password", pass, getResultString(result))

}

func getResultString(result bool) string {
	if result {
		return "succeeded"
	} else {
		return "failed"
	}
}

func Execute() {
	user := flag.String("u", "", "Username or user list to brute force")
	password := flag.String("p", "", "Password or password file to use for brute force")
	threads := flag.Int("t", 10, "Number of threads to use")
	serviceType := flag.String("s", "all", "Default all, Service type: ssh, ftp, smtp, etc")
	file := flag.String("f", "", "File to parse")
	host := flag.String("H", "", "Target in the format service://host:port")
	quiet := flag.Bool("q", false, "Supress the banner")
	timeout := flag.Int("T", 15, "Set timeout of bruteforce attempts")

	flag.Parse()

	modules.Banner(*quiet)

	if *user == "" || *password == "" {
		flag.Usage()
		os.Exit(1)
	}

	getSupportedServices := func(serviceType string) []string {
		if serviceType != "all" {
			supportedServices := strings.Split(serviceType, ",")
			for i := range supportedServices {
				supportedServices[i] = strings.TrimSpace(supportedServices[i])
			}
			return supportedServices
		}
		return nil
	}

	hosts, err := parseFile(*file)
	if err != nil && *file != "" {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	var users []string
	if isFile(*user) {
		var err error
		users, err = readUsersFromFile(*user)
		if err != nil {
			fmt.Println("Error reading user file:", err)
			os.Exit(1)
		}
	} else {
		users = append(users, *user)
	}

	var passwords []string
	if isFile(*password) {
		var err error
		passwords, err = readPasswordsFromFile(*password)
		if err != nil {
			fmt.Println("Error reading password file:", err)
			os.Exit(1)
		}
	} else {
		passwords = append(passwords, *password)
	}

	var hostsList []parse.Host
	for h := range hosts {
		hostsList = append(hostsList, h)
	}

	if *host != "" {
		var hostObj parse.Host
		if err := hostObj.Parse(*host); err != nil {
			fmt.Println("Error parsing host:", err)
			os.Exit(1)
		}
		hostsList = append(hostsList, hostObj)
	}

	bar, _ := pterm.DefaultProgressbar.WithTotal(len(hostsList) * len(users) * len(passwords)).WithTitle("Bruteforcing...").Start()
	var wg sync.WaitGroup
	sem := make(chan struct{}, *threads)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		pterm.DefaultSection.Println("\nReceived an interrupt signal, shutting down...")
		time.Sleep(5 * time.Second)
		bar.Stop()
		os.Exit(0)
	}()

	for _, h := range hostsList {
		wg.Add(1)
		sem <- struct{}{}
		go func(h parse.Host) {
			defer func() {
				<-sem
				wg.Done()
			}()
			for _, u := range users {
				for _, p := range passwords {
					wg.Add(1)
					sem <- struct{}{}
					go func(h parse.Host, u string, p string) {
						defer func() {
							<-sem
							wg.Done()
							bar.Increment()
						}()
						service := mapService(h.Service)
						if *serviceType != "all" && !contains(supportedServices, service) {
							return
						}
						bruteDone := make(chan bool)
						go func() {
							brute(h, u, p)
							bruteDone <- true
						}()

						select {
						case <-bruteDone:
						case <-time.After(time.Duration(*timeout) * time.Second):
							pterm.Color(pterm.FgRed).Println("Bruteforce timeout:", h.Service, "on host", h.Host, "port", h.Port, "with username", u, "and password", p)
						}
					}(h, u, p)
				}
			}
		}(h)
	}
	wg.Wait()
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	bar.Stop()
	if len(getSupportedServices(*serviceType)) > 0 {
		pterm.DefaultSection.Println("Supported services:", strings.Join(getSupportedServices(*serviceType), ", "))
	}
}
