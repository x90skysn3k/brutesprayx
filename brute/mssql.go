package brute

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func BruteMSSQL(host string, port int, user, password string) bool {
	connString := fmt.Sprintf("server=%s:%d;user id=%s;password=%s;database=master", host, port, user, password)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("mssql", connString)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		return false
	}
	return true
}
