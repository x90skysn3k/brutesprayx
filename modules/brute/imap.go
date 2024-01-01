package brute

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/emersion/go-imap/client"
)

func BruteIMAP(host string, port int, user, password string) bool {
	var (
		conn net.Conn
		err  error
	)

	conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)

	if err != nil {
		tlsDialer := &tls.Dialer{
			NetDialer: &net.Dialer{
				Timeout: 5 * time.Second,
			},
			Config: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		conn, err = tlsDialer.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

		if err != nil {
			return false
		}
	}

	c, err := client.New(conn)
	if err != nil {
		return false
	}

	err = c.Login(user, password)
	if err != nil {
		return false
	}

	err = c.Logout()
	if err != nil {
	}

	return true
}
