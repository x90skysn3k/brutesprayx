package modules

import (
	"time"

	"github.com/knadh/go-pop3"
)

func BrutePOP3(host string, port int, user, password string) bool {
	options := []pop3.Opt{
		{Host: host, Port: port, DialTimeout: 5 * time.Second},
		{Host: host, Port: port, TLSEnabled: true, DialTimeout: 5 * time.Second},
	}

	for _, opt := range options {
		p := pop3.New(opt)

		c, err := p.NewConn()
		if err != nil {
			continue
		}
		defer c.Quit()

		authDone := make(chan bool)
		go func() {
			err := c.Auth(user, password)
			authDone <- (err == nil)
		}()

		select {
		case authSuccess := <-authDone:
			if authSuccess {
				return true
			} else {
			}
		case <-time.After(5 * time.Second):
		}
	}

	return false
}
