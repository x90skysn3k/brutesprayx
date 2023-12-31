package modules

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func BruteMYSQL(host string, port int, user, password string) bool {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=5s", user, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", connString)
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
