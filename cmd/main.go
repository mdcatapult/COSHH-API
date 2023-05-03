package main

import (
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/db"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/server"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	if err := db.Connect("db"); err != nil {
		log.Fatal("Failed to start DB", err)
	}

	if err := server.Start(":8080"); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
