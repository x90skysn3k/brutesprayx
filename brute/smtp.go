package brute

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

func BruteSMTP(host string, port int, user, password string) bool {
	auth := smtp.PlainAuth("", user, password, host)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	smtpClient, err := smtp.NewClient(conn, host)
	if err != nil {
		return false
	}
	defer smtpClient.Quit()

	tlsDialer := &tls.Dialer{
		NetDialer: &net.Dialer{
			Timeout: 5 * time.Second,
		},
		Config: &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: true,
		},
	}

	tlsConn, err := tlsDialer.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return false
	}
	defer tlsConn.Close()

	if err := smtpClient.StartTLS(tlsDialer.Config); err != nil {
		return false
	}

	if err := smtpClient.Auth(auth); err == nil {
		return true
	} else {
	}

	return false
}
