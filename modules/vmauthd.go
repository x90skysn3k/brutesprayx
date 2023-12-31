package modules

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func BruteVMAuthd(host string, port int, user, password string) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		log.Printf("Failed to dial: %v", err)
		return false
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read: %v", err)
		return false
	}
	response := string(buf[:n])
	if strings.Contains(response, "SSL Required") {
		tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
		defer tlsConn.Close()
		conn = tlsConn
	} else {
		conn.SetReadDeadline(time.Time{})
	}

	cmd := fmt.Sprintf("USER %s\r\n", user)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("Failed to write: %v", err)
		return false
	}

	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read: %v", err)
		return false
	}
	response = string(buf[:n])
	if !strings.HasPrefix(response, "331 ") {
		log.Printf("Unexpected response: %s", response)
		return false
	}

	cmd = fmt.Sprintf("PASS %s\r\n", password)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("Failed to write: %v", err)
		return false
	}

	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read: %v", err)
		return false
	}
	response = string(buf[:n])

	if strings.HasPrefix(response, "230 ") {
		return true
	} else if strings.HasPrefix(response, "530 ") {
		return false
	} else {
		log.Printf("Unexpected response: %s", response)
		return false
	}
}
