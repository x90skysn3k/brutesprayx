package brute

import "github.com/x90skysn3k/brutesprayx/modules"

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
	"mysql":          "mysql",
	"vnc":            "vnc",
	//"ms-wbt-server":  "rdp",
}

func MapService(service string) string {
	if mappedService, ok := NAME_MAP[service]; ok {
		return mappedService
	}
	return service
}

func RunBrute(h modules.Host, u string, p string) {
	service := MapService(h.Service)
	var result bool

	switch service {
	case "ssh":
		result = BruteSSH(h.Host, h.Port, u, p)
	case "ftp":
		result = BruteFTP(h.Host, h.Port, u, p)
	case "mssql":
		result = BruteMSSQL(h.Host, h.Port, u, p)
	case "telnet":
		result = BruteTelnet(h.Host, h.Port, u, p)
	case "smbnt":
		result = BruteSMB(h.Host, h.Port, u, p)
	case "postgres":
		result = BrutePostgres(h.Host, h.Port, u, p)
	case "smtp":
		result = BruteSMTP(h.Host, h.Port, u, p)
	case "imap":
		result = BruteIMAP(h.Host, h.Port, u, p)
	case "pop3":
		result = BrutePOP3(h.Host, h.Port, u, p)
	case "snmp":
		result = BrutePOP3(h.Host, h.Port, u, p)
	case "mysql":
		result = BruteMYSQL(h.Host, h.Port, u, p)
	case "vmauthd":
		result = BruteVMAuthd(h.Host, h.Port, u, p)
	case "asterisk":
		//warning not tested
		result = BruteAsterisk(h.Host, h.Port, u, p)
	case "vnc":
		result = BruteVNC(h.Host, h.Port, u, p)
	//case "rdp":
	//	result = brute.BruteRDP(h.Host, h.Port, u, p)
	default:
		//fmt.Printf("Unsupported service: %s\n", h.Service)
		return
	}

	modules.PrintResult(service, h.Host, h.Port, u, p, result)
}
