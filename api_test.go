package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/chemical"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/db"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/server"
)

var chem = chemical.Chemical{
	CasNumber:       stringPtr("12345678"),
	Name:            "beans",
	PhotoPath:       stringPtr("blueberries"),
	MatterState:     stringPtr("liquid"),
	Quantity:        stringPtr("5"),
	Added:           &time.Time{},
	Expiry:          &time.Time{},
	SafetyDataSheet: stringPtr(""),
	StorageTemp:     "+4",
	IsArchived:      false,
	Hazards:         []string{"Explosive", "Flammable"},
}

func stringPtr(v string) *string {
	return &v
}

var client = &http.Client{}

func TestMain(m *testing.M) {
	if err := db.Connect("localhost"); err != nil {
		log.Fatal("Failed to start DB", err)
	}
	go func() {
		if err := server.Start(":8081"); err != nil {
			log.Fatal("Failed to start server", err)
		}
	}()

	// wait for server to start
	time.Sleep(time.Second * 2)

	status := m.Run()
	os.Exit(status)
}

func TestPostChemical(t *testing.T) {
	jsonChemical, err := json.Marshal(chem)
	assert.Nil(t, err, "Failed to marshal into chemical")

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/chemical", bytes.NewBuffer(jsonChemical))
	assert.Nil(t, err, "Failed to build post request")

	response, err := client.Do(req)
	assert.Nil(t, err, "Failed to send POST request")
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err, "Failed to read message body")

	var responseChemical chemical.Chemical
	err = json.Unmarshal(bodyBytes, &responseChemical)
	chem.Id = responseChemical.Id
	assert.Nil(t, err, "Failed to unmarshal into chemical")
	assert.Equal(t, chem, responseChemical)
}

func TestGetChemical(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/chemicals", nil)
	assert.Nil(t, err, "Failed to build GET request")

	response, err := client.Do(req)
	assert.Nil(t, err, "Failed to send GET request")

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err, "Failed to read message body")
	var responseChemicals []chemical.Chemical

	err = json.Unmarshal(bodyBytes, &responseChemicals)
	assert.Nil(t, err, "Failed to unmarshal into chemical")
	found := false
	for _, ch := range responseChemicals {
		if reflect.DeepEqual(ch, chem) {
			found = true
			break
		}
	}
	assert.True(t, found, "Error: values are not the same")
}

func TestPutChemical(t *testing.T) {
	putChem := chem
	putChem.Name = "bread"
	jsonChemical, err := json.Marshal(putChem)
	assert.Nil(t, err, "Failed to marshal into chemical")

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8081/chemical", bytes.NewBuffer(jsonChemical))
	assert.Nil(t, err, "Failed to build PUT request")

	response, err := client.Do(req)
	assert.Nil(t, err, "Failed to send PUT request")

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err, "Failed to read message body")

	var responseChemical chemical.Chemical
	err = json.Unmarshal(bodyBytes, &responseChemical)
	assert.Equal(t, putChem, responseChemical)
}

func TestPutHazards(t *testing.T) {
	putChem := chem
	putChem.Hazards = []string{"Corrosive", "Serious health hazard"}
	jsonChemical, err := json.Marshal(putChem)
	assert.Nil(t, err, "Failed to marshal into chemical")

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8081/hazards", bytes.NewBuffer(jsonChemical))
	assert.Nil(t, err, "Failed to build PUT request")

	response, err := client.Do(req)
	assert.Nil(t, err, "Failed to send PUT request")

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err, "Failed to read message body")

	var responseChemical chemical.Chemical
	err = json.Unmarshal(bodyBytes, &responseChemical)
	assert.Equal(t, putChem, responseChemical)
}
