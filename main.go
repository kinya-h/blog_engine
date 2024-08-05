package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kinya-h/blog_engine/api"
	"github.com/kinya-h/blog_engine/db"
	"github.com/kinya-h/blog_engine/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	config, err := util.LoadConfig(".")
	fmt.Printf("CONFIG:: %+v", config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	}

	conn, err := sql.Open("mysql", "root:root@/blog?parseTime=true")

	if err != nil {
		fmt.Println("AN ERROR OCCURED ", err)
	}
	fmt.Print("CONNECTED SUCCESSFULLY\n")

	queries := db.New(conn)

	server, err := api.NewServer(config, queries)

	if err != nil {
		fmt.Println("AN ERROR OCCURED ", err)

	}

	server.Start()

}
