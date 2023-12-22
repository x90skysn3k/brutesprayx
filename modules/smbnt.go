package modules

import (
	"fmt"
	"net"
	"time"

	"github.com/hirochachacha/go-smb2"
)

func BruteSMB(host string, port int, user, password string) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return false
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return false
	}
	defer s.Logoff()

	_, err = s.ListSharenames()
	if err != nil {
		return false
	}

	return true
}
