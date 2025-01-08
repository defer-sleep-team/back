package main

import (
	"flag"
	"log"

	db "github.com/defer-sleep-team/Aether_backend/database/internal/dbworker"
	rout "github.com/defer-sleep-team/Aether_backend/database/router"
)

func main() {
	key := flag.String("connect_url", "connect_url", "the key used to connxt to db")
	flag.Parse()

	D, err := db.NewDBConnection(*key)
	if err != nil {
		log.Fatal("Error creating database: ", err)
	}
	defer D.Close()

	log.Print("Started server successfully")
	err = rout.Rout(D)
	if err != nil {
		log.Fatal("Error creating router: ", err)
	}

}
