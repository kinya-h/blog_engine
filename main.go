package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kinya-h/blog_engine/api"
	"github.com/kinya-h/blog_engine/db"
)

func main() {

	conn, err := sql.Open("mysql", "root:root@/blog?parseTime=true")

	if err != nil {
		fmt.Println("AN ERROR OCCURED ", err)
	}
	fmt.Print("CONNECTED SUCCESSFULLY\n")

	queries := db.New(conn)

	server := api.NewServer(queries)
	server.Start()

}
