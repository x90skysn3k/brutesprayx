package brutesprayx

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"github.com/x90skysn3k/brutesprayx/banner"
	"github.com/x90skysn3k/brutesprayx/brute"
	"github.com/x90skysn3k/brutesprayx/modules"
)

var masterServiceList = []string{"ssh", "ftp", "smtp", "mssql", "telnet", "smbnt", "postgres", "imap", "pop3", "snmp", "mysql", "vmauthd", "asterisk", "vnc"}

var version = "v2.1.0"

func Execute() {
	user := flag.String("u", "", "Username or user list to brute force")
	password := flag.String("p", "", "Password or password file to use for brute force")
	threads := flag.Int("t", 10, "Number of threads to use")
	serviceType := flag.String("s", "all", "Default all, Service type: ssh, ftp, smtp, etc")
	listServices := flag.Bool("S", false, "List all supported services")
	file := flag.String("f", "", "File to parse")
	host := flag.String("H", "", "Target in the format service://host:port")
	quiet := flag.Bool("q", false, "Supress the banner")
	timeout := flag.Int("T", 15, "Set timeout of bruteforce attempts")

	flag.Parse()

	banner.Banner(version, *quiet)

	supportedServices := getSupportedServices(*serviceType)

	if *listServices {
		pterm.DefaultSection.Println("Supported services:", strings.Join(supportedServices, ", "))
		os.Exit(1)
	}

	if *host == "" && *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	hosts, err := modules.ParseFile(*file)
	if err != nil && *file != "" {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	users := getUsers(*user)
	passwords := getPasswords(*password)
	hostsMap := make(map[modules.Host]struct{})
	for host := range hosts {
		hostsMap[host] = struct{}{}
	}

	hostsList := getHosts(hostsMap, *host)

	bar := createProgressBar(hostsList, users, passwords, supportedServices)

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

	for _, service := range supportedServices {
		wg.Add(1)
		sem <- struct{}{}
		go func(service string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			bruteforce(hostsList, users, passwords, service, *timeout, bar)
		}(service)
	}
	wg.Wait()
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	bar.Stop()
}

func getSupportedServices(serviceType string) []string {
	if serviceType != "all" {
		supportedServices := strings.Split(serviceType, ",")
		for i := range supportedServices {
			supportedServices[i] = strings.TrimSpace(supportedServices[i])
		}
		return supportedServices
	}
	return masterServiceList
}

func getUsers(user string) []string {
	var users []string
	if user != "" {
		if modules.IsFile(user) {
			var err error
			users, err = modules.ReadUsersFromFile(user)
			if err != nil {
				fmt.Println("Error reading user file:", err)
				os.Exit(1)
			}
		} else {
			users = append(users, user)
		}
	} else {
		users = modules.GetUsersFromDefaultWordlist(version)
	}
	return users
}

func getPasswords(password string) []string {
	var passwords []string
	if password != "" {
		if modules.IsFile(password) {
			var err error
			passwords, err = modules.ReadPasswordsFromFile(password)
			if err != nil {
				fmt.Println("Error reading password file:", err)
				os.Exit(1)
			}
		} else {
			passwords = append(passwords, password)
		}
	} else {
		passwords = modules.GetPasswordsFromDefaultWordlist(version)
	}
	return passwords
}

func getHosts(hosts map[modules.Host]struct{}, host string) []modules.Host {
	var hostsList []modules.Host
	for h := range hosts {
		hostsList = append(hostsList, h)
	}

	if host != "" {
		var hostObj modules.Host
		host, err := hostObj.Parse(host)
		if err != nil {
			fmt.Println("Error parsing host:", err)
			os.Exit(1)
		}
		hostsList = append(hostsList, host...)
	}
	return hostsList
}

func createProgressBar(hostsList []modules.Host, users, passwords []string, supportedServices []string) *pterm.ProgressbarPrinter {
	var nopassServices int
	for _, service := range supportedServices {
		if service == "vnc" {
			nopassServices++
		}
		if service == "snmp" {
			nopassServices++
		}
	}

	bar, _ := pterm.DefaultProgressbar.WithTotal(len(hostsList)*len(users)*len(passwords) - nopassServices*len(users)).WithTitle("Bruteforcing...").Start()
	return bar
}

func bruteforce(hostsList []modules.Host, users, passwords []string, service string, timeout int, bar *pterm.ProgressbarPrinter) {
	if service == "vnc" || service == "snmp" {
		u := ""
		for _, h := range hostsList {
			if h.Service == service {
				for _, p := range passwords {
					bruteforceHost(h, u, p, timeout, bar)
				}
			}
		}
	} else {
		for _, h := range hostsList {
			if h.Service == service {
				for _, u := range users {
					for _, p := range passwords {
						bruteforceHost(h, u, p, timeout, bar)
					}
				}
			}
		}
	}
}

func bruteforceHost(h modules.Host, u string, p string, timeout int, bar *pterm.ProgressbarPrinter) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		bruteDone := make(chan bool)
		go func() {
			brute.RunBrute(h, u, p)
			bruteDone <- true
		}()

		select {
		case <-bruteDone:
		case <-time.After(time.Duration(timeout) * time.Second):
			pterm.Color(pterm.FgRed).Println("Bruteforce timeout:", h.Service, "on host", h.Host, "port", h.Port, "with username", u, "and password", p)
		}
		bar.Increment()
	}()
	wg.Wait()
}
