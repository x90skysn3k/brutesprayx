package modules

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Host struct {
	Service string
	Host    string
	Port    int
	CIDR    string
}

type NexposeNode struct {
	Address   string            `xml:"address,attr"`
	Endpoints []NexposeEndpoint `xml:"endpoints>endpoint"`
}

type NexposeEndpoint struct {
	Port     string         `xml:"port,attr"`
	Status   string         `xml:"status,attr"`
	Protocol string         `xml:"protocol,attr"`
	Service  NexposeService `xml:"services>service"`
}

type NexposeService struct {
	Name string `xml:"name,attr"`
}

type NessusReport struct {
	Hosts []NessusHost `xml:"Report>ReportHost"`
}

type NessusHost struct {
	Name  string       `xml:"name,attr"`
	Items []NessusItem `xml:"ReportItem"`
}

type NessusItem struct {
	Port    string `xml:"port,attr"`
	SvcName string `xml:"svc_name,attr"`
}

type NmapRun struct {
	Hosts []NmapHost `xml:"host"`
}

type NmapHost struct {
	Addresses []NmapAddress `xml:"address"`
	Ports     []NmapPort    `xml:"ports>port"`
}

type NmapAddress struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}

type NmapPort struct {
	PortId    string      `xml:"portid,attr"`
	Protocol  string      `xml:"protocol,attr"`
	PortState NmapState   `xml:"state"`
	Service   NmapService `xml:"service"`
}

type NmapState struct {
	State string `xml:"state,attr"`
}

type NmapService struct {
	Name string `xml:"name,attr"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ParseGNMAP(filename string) (map[Host]int, error) {
	supported := []string{"ssh", "ftp", "postgres", "telnet", "mysql", "ms-sql-s", "shell",
		"vnc", "imap", "imaps", "nntp", "pcanywheredata", "pop3", "pop3s",
		"exec", "login", "microsoft-ds", "smtp", "smtps", "submission",
		"svn", "iss-realsecure", "snmptrap", "snmp", "ms-wbt-server"}

	hosts := make(map[Host]int)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, name := range supported {
			matches := regexp.MustCompile(fmt.Sprintf(`([0-9][0-9]*)/open/[a-z][a-z]*//%s`, name))
			portMatches := matches.FindStringSubmatch(line)
			if len(portMatches) == 0 {
				continue
			}
			port, _ := strconv.Atoi(portMatches[1])
			ipMatches := regexp.MustCompile(`[0-9]+(?:\.[0-9]+){3}`).FindAllString(line, -1)

			for _, ip := range ipMatches {
				h := Host{Service: name, Host: ip, Port: port}
				hosts[h] = 1
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}
func ParseJSON(filename string) (map[Host]int, error) {
	supported := []string{"ssh", "ftp", "postgres", "telnet", "mysql", "ms-sql-s", "shell",
		"vnc", "imap", "imaps", "nntp", "pcanywheredata", "pop3", "pop3s",
		"exec", "login", "microsoft-ds", "smtp", "smtps", "submission",
		"svn", "iss-realsecure", "snmptrap", "snmp"}

	hosts := make(map[Host]int)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err != nil {
			break
		}
		host, _ := data["host"].(string)
		port, _ := data["port"].(string)
		name, _ := data["service"].(string)
		if contains(supported, name) {
			p, _ := strconv.Atoi(port)
			h := Host{Service: name, Host: host, Port: p}
			hosts[h] = 1
		}
	}

	return hosts, nil
}
func ParseXML(filename string) (map[Host]int, error) {
	supported := []string{"ssh", "ftp", "postgresql", "telnet", "mysql", "ms-sql-s", "rsh",
		"vnc", "imap", "imaps", "nntp", "pcanywheredata", "pop3", "pop3s",
		"exec", "login", "microsoft-ds", "smtp", "smtps", "submission",
		"svn", "iss-realsecure", "snmptrap", "snmp"}

	hosts := make(map[Host]int)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var report NmapRun
	err = decoder.Decode(&report)
	if err != nil {
		return nil, err
	}

	for _, host := range report.Hosts {
		ip := host.Addresses[0].Addr
		for _, port := range host.Ports {
			if port.PortState.State == "open" {
				name := port.Service.Name
				if contains(supported, name) {
					p, _ := strconv.Atoi(port.PortId)
					h := Host{Service: name, Host: ip, Port: p}
					hosts[h] = 1
				}
			}
		}
	}

	return hosts, nil
}
func ParseNexpose(filename string) (map[Host]int, error) {
	supported := []string{"ssh", "ftp", "postgresql", "telnet", "mysql", "ms-sql-s", "rsh",
		"vnc", "imap", "imaps", "nntp", "pcanywheredata", "pop3", "pop3s",
		"exec", "login", "microsoft-ds", "smtp", "smtps", "submission",
		"svn", "iss-realsecure", "snmptrap", "snmp", "cifs"}

	hosts := make(map[Host]int)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var nodes []NexposeNode
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "node" {
				var node NexposeNode
				decoder.DecodeElement(&node, &se)
				nodes = append(nodes, node)
			}
		}
	}

	for _, node := range nodes {
		ip := node.Address
		for _, port := range node.Endpoints {
			if port.Status == "open" {
				name := port.Service.Name
				name = strings.ToLower(name)
				if contains(supported, name) {
					p, _ := strconv.Atoi(port.Port)
					h := Host{Service: name, Host: ip, Port: p}
					hosts[h] = 1
				}
			}
		}
	}
	return hosts, nil
}
func ParseNessus(filename string) (map[Host]int, error) {
	supported := []string{"ssh", "ftp", "postgresql", "telnet", "mysql", "ms-sql-s", "rsh",
		"vnc", "imap", "imaps", "nntp", "pcanywheredata", "pop3", "pop3s",
		"exec", "login", "microsoft-ds", "smtp", "smtps", "submission",
		"svn", "iss-realsecure", "snmptrap", "snmp", "cifs"}

	hosts := make(map[Host]int)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var report NessusReport
	err = decoder.Decode(&report)
	if err != nil {
		return nil, err
	}
	for _, host := range report.Hosts {
		ip := host.Name
		for _, port := range host.Items {
			if port.Port != "0" {
				name := port.SvcName
				if contains(supported, name) {
					p, _ := strconv.Atoi(port.Port)
					h := Host{Service: name, Host: ip, Port: p}
					hosts[h] = 1
				}
			}
		}
	}

	return hosts, nil
}
func ParseList(filename string) (map[Host]int, error) {
	supportedServices := []string{"ssh", "ftp", "smtp", "mssql", "telnet", "smbnt", "postgres", "imap", "pop3", "snmp", "mysql", "iss-realsecure"}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hosts := make(map[Host]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		service := parts[0]
		ip := parts[1]
		port, _ := strconv.Atoi(parts[2])
		h := Host{Service: service, Host: ip, Port: port}
		hosts[h] = 1

		var found bool
		for _, services := range supportedServices {
			if service == services {
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("unsupported service: %s", h.Service)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}
func (h *Host) Parse(host string) error {
	supportedServices := []string{"ssh", "ftp", "smtp", "mssql", "telnet", "smbnt", "postgres", "imap", "pop3", "snmp", "mysql", "vmauthd"}

	parts := strings.Split(host, "://")
	if len(parts) != 2 {
		return fmt.Errorf("invalid host format: %s", host)
	}

	h.Service = parts[0]
	remaining := parts[1]

	portIndex := strings.LastIndex(remaining, ":")
	if portIndex == -1 {
		return fmt.Errorf("invalid host format: %s", host)
	}

	portStr := remaining[portIndex+1:]
	remaining = remaining[:portIndex]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port in host: %s", host)
	}
	var found bool
	for _, service := range supportedServices {
		if h.Service == service {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("unsupported service: %s", h.Service)
	}

	h.Port = port
	h.Host = remaining

	return nil
}
