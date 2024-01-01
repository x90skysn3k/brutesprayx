package brute

import (
	"fmt"
	"strings"

	"github.com/wenerme/astgo/ami"
)

func BruteAsterisk(host string, port int, user, password string) bool {
	boot := make(chan *ami.Message, 1)

	conn, err := ami.Connect(
		fmt.Sprintf("%s:%d", host, port),
		ami.WithAuth(user, password),
		ami.WithSubscribe(ami.SubscribeFullyBootedChanOnce(boot)),
	)
	if err != nil {
		return false
	}
	defer conn.Close()
	<-boot

	if strings.Contains(conn.Close().Error(), "Message: Authentication accepted") {
		return true
	} else {
		return false
	}
}
