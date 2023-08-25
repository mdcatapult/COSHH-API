package server

import (
	"context"
	"encoding/csv"
	adapter "github.com/gwatts/gin-adapter"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/chemical"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/db"
)

var config Config

type Config struct {
	LabsCSV       string `env:"LABS_CSV,default=/mnt/labs.csv"`
	ProjectsCSV   string `env:"PROJECTS_CSV,default=/mnt/projects.csv"`
	Auth0Audience string `env:"AUTH0_AUDIENCE,required"`
	Auth0Domain   string `env:"AUTH0_DOMAIN,required"`
}

type (
	jwtValidator func(audience string, domain string) func(next http.Handler) http.Handler
)

func Start(port string, validator jwtValidator) error {

	ctx := context.Background()

	if err := envconfig.Process(ctx, &config); err != nil {
		log.Println("Server start env vars unset or incorrect, using default config")
	} else {
		log.Println("Server using config from env vars")
	}
	r := gin.Default()
	r.Use(corsMiddleware())

	log.Printf("Server using Audience=%s, Domain=%s", config.Auth0Audience, config.Auth0Domain)
	r.GET("/chemicals", getChemicals)
	r.PUT("/chemical", adapter.Wrap(validator(config.Auth0Audience, config.Auth0Domain)), updateChemical)

	r.POST("/chemical", adapter.Wrap(validator(config.Auth0Audience, config.Auth0Domain)), insertChemical)

	r.GET("/cupboards", getCupboards)

	r.PUT("/hazards", adapter.Wrap(validator(config.Auth0Audience, config.Auth0Domain)), updateHazards)

	r.GET("/labs", getLabs)

	r.GET("/lab-cupboards", getCupboardsForLab)

	r.GET("/projects", getProjects)

	//This route is here to allow standalone testing of authentication using curl
	r.GET("/protected", adapter.Wrap(validator(config.Auth0Audience, config.Auth0Domain)), protectedRoute)

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

func getCupboardsForLab(c *gin.Context) {
	lab := c.Query("lab")
	chemicals, err := db.GetCupboardsForLab(lab)
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

	labsFile, err := os.Open(config.LabsCSV)
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
	projectsFile, err := os.Open(config.ProjectsCSV)
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

// Purely for auth testing
func protectedRoute(c *gin.Context) {
	c.JSON(http.StatusOK, "You have successfully authenticated")
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
