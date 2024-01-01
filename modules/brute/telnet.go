package modules

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func BruteTelnet(host string, port int, user, password string) bool {
	connection, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
	if err != nil {
		//fmt.Printf("Failed to dial %s:%d: %s\n", host, port, err)
		return false
	}
	defer connection.Close()

	reader := bufio.NewReader(connection)

	serverMessage, err := reader.ReadString('\n')
	if err != nil {
		//fmt.Printf("Failed to read from %s:%d: %s\n", host, port, err)
		return false
	}

	fmt.Fprintf(connection, "%s\n", user)

	serverMessage, err = reader.ReadString('\n')
	if err != nil {
		//fmt.Printf("Failed to read from %s:%d: %s\n", host, port, err)
		return false
	}

	fmt.Fprintf(connection, "%s\n", password)

	serverMessage, err = reader.ReadString('\n')
	if err != nil {
		//fmt.Printf("Failed to read from %s:%d: %s\n", host, port, err)
		return false
	}

	if strings.Contains(serverMessage, "Login successful") {
		//fmt.Printf("Attempt %s:%d:%s:%s successful\n", host, port, user, password)
		return true
	} else {
		//fmt.Printf("Attempt %s:%d:%s:%s failed: %s\n", host, port, user, password, serverMessage)
		return false
	}
}
