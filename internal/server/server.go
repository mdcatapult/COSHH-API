package server

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/chemical"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/db"
)

// TODO these shouldn't be hard-coded
var labs_csv = "/mnt/labs.csv"
var projects_csv = "/mnt/projects.csv"

type Config struct {
	LabsCSV     string `env:"LABS_CSV,required"`
	ProjectsCSV string `env:"PROJECTS_CSV, required"`
}

func Start(port string) error {

	ctx := context.Background()
	var config Config

	if err := envconfig.Process(ctx, &config); err != nil {
		fmt.Println("Env vars unset or incorrect, using default config")
	} else {
		fmt.Println("Using config from env vars")
		labs_csv = config.LabsCSV
		projects_csv = config.ProjectsCSV
	}
	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/chemicals", getChemicals)
	r.PUT("/chemical", updateChemical)
	r.POST("/chemical", insertChemical)

	r.GET("/cupboards", getCupboards)

	r.PUT("/hazards", updateHazards)

	r.GET("/labs", getLabs)

	r.GET("/projects", getProjects)

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

func getCupboards(c *gin.Context) {
	values, err := db.SelectAllCupboards()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, values)
}

func updateHazards(c *gin.Context) {
	var chemical chemical.Chemical
	if err := c.BindJSON(&chemical); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := db.DeleteHazards(chemical)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if len(chemical.Hazards) > 0 {
		err = db.InsertHazards(chemical)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}

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

func getProjects(c *gin.Context) {
	projectsFile, err := os.Open(projects_csv)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer projectsFile.Close()

	csvReader := csv.NewReader(projectsFile)
	projects, err := csvReader.ReadAll()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, projects)
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
