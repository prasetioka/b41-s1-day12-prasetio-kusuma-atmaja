package connection

import (
	"context"

	"fmt"

	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	var err error
	
	databaseUrl := "postgres://postgres:corolladx86@localhost:5432/personal_web_b41"

	Conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success connect to database")
}
