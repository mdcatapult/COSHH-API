package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
	Location        string    `json:"location" db:"lab_location"`
	StorageTemp     string    `json:"storageTemp" db:"storage_temp"`
	IsArchived      bool      `json:"isArchived" db:"is_archived"`
	Hazards         []string  `json:"hazards" db:"-"`
	DBHazards       string    `json:"-" db:"hazards"`
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
	r.Use(CORSMiddleware())

	r.GET("/chemicals", getChemicals)
	r.PUT("/chemical", updateChemical)
	r.POST("/chemical", insertChemical)

	if err := r.Run(); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func getChemicals(c *gin.Context) {

	var chemicals []Chemical
	query := `
		SELECT 
		c.*, 
		string_agg(CAST(c2h.hazard AS VARCHAR(255)), ',') AS hazards 
		FROM chemical c JOIN chemical_to_hazard c2h ON c.cas_number = c2h.cas_number 
		GROUP BY c.cas_number`

	if err := db.Select(&chemicals, query); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for i := range chemicals {
		chemicals[i].Hazards = strings.Split(chemicals[i].DBHazards, ",")
	}

	c.JSON(http.StatusOK, chemicals)
}

func updateChemical(c *gin.Context) {

	var chemical Chemical
	if err := c.BindJSON(&chemical); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	query := `
		UPDATE chemical 
		SET 
			chemical_name = :chemical_name,
			photo_path = :photo_path,
			matter_state = :matter_state,
			quantity = :quantity,
			added = :added,
			expiry = :expiry,
			safety_data_sheet = :safety_data_sheet,
			coshh_link = :coshh_link,
			lab_location = :lab_location,
			storage_temp = :storage_temp,
			is_archived = :is_archived
		WHERE cas_number = :cas_number
	`

	if _, err := db.NamedExec(query, chemical); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func insertChemical(c *gin.Context) {

	var chemical Chemical
	if err := c.BindJSON(&chemical); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	query := `INSERT INTO chemical (
		cas_number,
		chemical_name,
		photo_path,
		matter_state,
		quantity,
		added,
		expiry,
		safety_data_sheet,
		coshh_link,
		lab_location,
		storage_temp,
		is_archived 
	)VALUES (
		:cas_number,
		:chemical_name,
		:photo_path,
		:matter_state,
		:quantity,
		:added,
		:expiry,
		:safety_data_sheet,
		:coshh_link,
		:lab_location,
		:storage_temp,
		:is_archived
	)`

	_, err := db.NamedExec(query, chemical)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, chemical)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
