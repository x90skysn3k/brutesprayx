package modules

import (
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
)

func BruteFTP(host string, port int, user, password string) bool {
	conn, err := ftp.Dial(host+":"+strconv.Itoa(port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return false
	}
	defer conn.Quit()

	err = conn.Login(user, password)
	if err != nil {
		return false
	}

	return true
}
