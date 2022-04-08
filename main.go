package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type Chemical struct {
	CasNumber       string    `json:"casNumber" db:"cas_number"`
	ChemicalName    string    `json:"chemicalName" db:"chemical_name"`
	PhotoPath       string    `json:"photoPath" db:"photo_path"`
	MatterState     string    `json:"matterState" db:"matter_state"`
	Quantity        string    `json:"quantity" db:"quantity"`
	Added           time.Time `json:"added" db:"added"`
	Expiry          time.Time `json:"expiry" db:"expiry"`
	SafetyDataSheet string    `json:"safetyDataSheet" db:"safety_data_sheet"`
	CoshhLink       string    `json:"coshhLink" db:"coshh_link"`
	Hazards         []string  `json:"hazards" db:"hazards"`
}

var db *sqlx.DB

func main() {

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "coshh"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Could not connect to the database, " + err.Error())
		return
	}

	r := gin.Default()
	r.GET("/chemicals", getChemicals)
	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func getChemicals(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	var res []Chemical
	query := `SELECT c.cas_number, c.chemical_name, string_agg(CAST(c2h.hazard AS VARCHAR(255)), ',') as hazards FROM chemical c JOIN chemical_to_hazard c2h ON c.cas_number = c2h.cas_number GROUP BY c.cas_number`
	if err := db.Select(&res, query); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
