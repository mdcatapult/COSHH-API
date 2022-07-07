package server

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/chemical"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/db"
)

var labs_csv = "/mnt/labs.csv"

type Config struct {
	LabsCSV  string `env:"LABS_CSV,required"`
}

func Start(port string) error {

	ctx := context.Background()
	var config Config

	if err := envconfig.Process(ctx, &config); err != nil {
		fmt.Println("Env vars unset or incorrect, using default config")
	} else {
		fmt.Println("Using config from env vars")
		labs_csv = config.LabsCSV
	}
	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/chemicals", getChemicals)
	r.PUT("/chemical", updateChemical)
	r.POST("/chemical", insertChemical)

	r.GET("/labs", getLabs)

	r.PUT("/hazards", updateHazards)

	return r.Run(port)
}

func getChemicals(c *gin.Context) {

	chemicals, err := db.SelectAllChemicals()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, chemicals)
}

func updateChemical(c *gin.Context) {

	var chemical chemical.Chemical
	if err := c.BindJSON(&chemical); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := db.UpdateChemical(chemical); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, chemical)
}

func insertChemical(c *gin.Context) {

	var chemical chemical.Chemical
	if err := c.BindJSON(&chemical); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := db.InsertChemical(chemical)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	chemical.Id = id

	c.JSON(http.StatusOK, chemical)
}

func getLabs(c *gin.Context) {
	labsFile, err := os.Open(labs_csv)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer labsFile.Close()

	csvReader := csv.NewReader(labsFile)
	labs, err := csvReader.Read()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, labs)
}

func updateHazards(c *gin.Context) {
	fmt.Printf("in updateHazards")
	var chemical chemical.Chemical
	if err := c.BindJSON(&chemical); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := db.DeleteHazards(chemical)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	err = db.InsertHazards(chemical)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, chemical)
}

func corsMiddleware() gin.HandlerFunc {
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
