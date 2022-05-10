package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/chemical"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/db"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/server"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var chem = chemical.Chemical{
	CasNumber:       "12345678",
	Name:            "beans",
	PhotoPath:       "blueberries",
	MatterState:     "liquid",
	Quantity:        "5",
	Added:           time.Time{},
	Expiry:          time.Time{},
	SafetyDataSheet: "",
	StorageTemp:     "+4",
	IsArchived:      false,
	Hazards:         []string{"Explosive", "Flammable"},
}

var client = &http.Client{}

func TestMain(m *testing.M) {
	if err := db.Connect(); err != nil {
		log.Fatal("Failed to start DB", err)
	}
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Failed to start server", err)
		}
	}()
	m.Run()
}

func TestPostChemical(t *testing.T) {
	jsonChemical, err := json.Marshal(chem)
	assert.Nil(t, err, "Failed to marshal into chemical")

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/chemical", bytes.NewBuffer(jsonChemical))
	assert.Nil(t, err, "Failed to build post request")

	response, err := client.Do(req)
	assert.Nil(t, err, "Failed to send POST request")
	assert.Equal(t, http.StatusOK, response.StatusCode)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err, "Failed to read message body")

	var responseChemical chemical.Chemical
	err = json.Unmarshal(bodyBytes, &responseChemical)
	assert.Nil(t, err, "Failed to unmarshal into chemical")
	assert.Equal(t, chem, responseChemical)
}

func TestGetChemical(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/chemicals", nil)
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
	assert.True(t, found)
}

func TestPutChemical(t *testing.T) {
	putChem := chem
	putChem.Name = "bread"
	jsonChemical, err := json.Marshal(putChem)
	assert.Nil(t, err, "Failed to marshal into chemical")

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/chemical", bytes.NewBuffer(jsonChemical))
	assert.Nil(t, err, "Failed to build PUT request")

	response, err := client.Do(req)
	assert.Nil(t, err, "Failed to send PUT request")

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err, "Failed to read message body")

	var responseChemical chemical.Chemical
	err = json.Unmarshal(bodyBytes, &responseChemical)
	assert.Equal(t, putChem, responseChemical)
}
