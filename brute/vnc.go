package brute

import (
	"fmt"
	"net"
	"time"

	"github.com/mitchellh/go-vnc"
)

func BruteVNC(host string, port int, user string, password string) bool {
	config := &vnc.ClientConfig{
		Auth: []vnc.ClientAuth{
			&vnc.PasswordAuth{
				Password: password,
			},
		},
	}
	const defaultTimeout = 5 * time.Second
	timer := time.NewTimer(defaultTimeout)
	defer timer.Stop()

	type result struct {
		client *vnc.ClientConn
		err    error
	}
	done := make(chan result)
	go func() {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), defaultTimeout)
		if err != nil {
			done <- result{nil, err}
			return
		}
		client, err := vnc.Client(conn, config)
		done <- result{client, err}
	}()

	select {
	case <-timer.C:
		return false
	case result := <-done:
		if result.err != nil {
			return false
		}
		result.client.Close()
		return true
	}
}
