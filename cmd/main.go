package main

import (
	"log"

	_ "github.com/lib/pq"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/db"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/server"
)

func main() {

	if err := db.Connect("db"); err != nil {
		log.Fatal("Failed to start DB", err)
	}

	if err := server.Start(":8080"); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
