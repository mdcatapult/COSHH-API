package main

import (
	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/middleware"
	_ "github.com/lib/pq"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/db"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/server"
	"log"
)

func main() {

	if err := db.Connect("db"); err != nil {
		log.Fatal("Failed to start DB", err)
	}

	if err := server.Start(":8080", middleware.ValidateJWT); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
