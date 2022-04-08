package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
}

var db *sqlx.DB

func main() {

	var err error
	db, err = sqlx.Connect("postgres", "postgres://postgres@localhost/coshh?sslmode=disable")
	if err != nil {
		log.Fatal("Could not connect to the database, " + err.Error())
		return
	}

	r := gin.Default()
	r.GET("/coshh", getCOSHH)
	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func getCOSHH(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	var res []Chemical
	query := `SELECT * FROM coshh`
	if err := db.Select(&res, query); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
