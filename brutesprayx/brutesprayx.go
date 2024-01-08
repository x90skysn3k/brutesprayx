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

var version = "v2.1.2"

func Execute() {
	user := flag.String("u", "", "Username or user list to bruteforce")
	password := flag.String("p", "", "Password or password file to use for bruteforce")
	threads := flag.Int("t", 10, "Number of threads to use")
	serviceType := flag.String("s", "all", "Service type: ssh, ftp, smtp, etc; Default all")
	listServices := flag.Bool("S", false, "List all supported services")
	file := flag.String("f", "", "File to parse; Supported: Nmap, Nessus, Nexpose, Lists, etc")
	host := flag.String("H", "", "Target in the format service://host:port, CIDR ranges supported,\n default port will be used if not specified")
	quiet := flag.Bool("q", false, "Supress the banner")
	timeout := flag.Int("T", 15, "Set timeout of bruteforce attempts")

	flag.Parse()

	banner.Banner(version, *quiet)

	getSupportedServices := func(serviceType string) []string {
		if serviceType != "all" {
			supportedServices := strings.Split(serviceType, ",")
			for i := range supportedServices {
				supportedServices[i] = strings.TrimSpace(supportedServices[i])
			}
			return supportedServices
		}
		return masterServiceList
	}

	if *listServices {
		pterm.DefaultSection.Println("Supported services:", strings.Join(getSupportedServices(*serviceType), ", "))
		os.Exit(1)
	} else {
		if flag.NFlag() == 0 {
			flag.Usage()
			pterm.DefaultSection.Println("Supported services:", strings.Join(getSupportedServices(*serviceType), ", "))
			os.Exit(1)
		}
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

	var users []string
	if *user != "" {
		if modules.IsFile(*user) {
			var err error
			users, err = modules.ReadUsersFromFile(*user)
			if err != nil {
				fmt.Println("Error reading user file:", err)
				os.Exit(1)
			}
		} else {
			users = append(users, *user)
		}
	} else {
		users = modules.GetUsersFromDefaultWordlist(version)
	}

	var passwords []string
	if *password != "" {
		if modules.IsFile(*password) {
			var err error
			passwords, err = modules.ReadPasswordsFromFile(*password)
			if err != nil {
				fmt.Println("Error reading password file:", err)
				os.Exit(1)
			}
		} else {
			passwords = append(passwords, *password)
		}
	} else {
		passwords = modules.GetPasswordsFromDefaultWordlist(version)
	}

	var hostsList []modules.Host
	for h := range hosts {
		hostsList = append(hostsList, h)
	}

	if *host != "" {
		var hostObj modules.Host
		host, err := hostObj.Parse(*host)
		if err != nil {
			fmt.Println("Error parsing host:", err)
			os.Exit(1)
		}
		hostsList = append(hostsList, host...)
	}

	supportedServices := getSupportedServices(*serviceType)

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
			if service == "vnc" || service == "snmp" {
				u := ""
				for _, h := range hostsList {
					if h.Service == service {
						for _, p := range passwords {
							wg.Add(1)
							sem <- struct{}{}
							go func(h modules.Host, u string, p string) {
								defer func() {
									<-sem
									wg.Done()
									bar.Increment()
								}()
								bruteDone := make(chan bool)
								go func() {
									brute.RunBrute(h, u, p)
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
				}
			} else {
				for _, h := range hostsList {
					if h.Service == service {
						for _, u := range users {
							for _, p := range passwords {
								wg.Add(1)
								sem <- struct{}{}
								go func(h modules.Host, u string, p string) {
									defer func() {
										<-sem
										wg.Done()
										bar.Increment()
									}()
									bruteDone := make(chan bool)
									go func() {
										brute.RunBrute(h, u, p)
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
					}
				}
			}
		}(service)
	}
	wg.Wait()
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	bar.Stop()
}
