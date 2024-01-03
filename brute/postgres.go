package brute

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func BrutePostgres(host string, port int, user, password string) bool {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		//fmt.Printf("Attempt %s:%s:%s failed: %s\n", host, user, password, err)
		return false
	}
	defer db.Close()

	//fmt.Printf("Attempt POSTGRES on %s:%s:%s SUCCESS\n", host, user, password)
	return true
}
