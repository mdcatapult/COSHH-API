/*
* Copyright $today.year Medicines Discovery Catapult
* Licensed under the Apache License, Version 2.0 (the "Licence");
* you may not use this file except in compliance with the Licence.
* You may obtain a copy of the Licence at
*     http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the Licence is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the Licence for the specific language governing permissions and
* limitations under the Licence.
*/

package main

import (
	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/middleware"
	_ "github.com/lib/pq"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/db"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/server"
	"log"
)

func main() {

	if err := db.Connect(); err != nil {
		log.Fatal("Failed to start DB", err)
	}

	if err := server.Start(":8080", middleware.ValidateJWT); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
