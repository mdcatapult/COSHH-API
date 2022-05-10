package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/chemical"
	"io/ioutil"
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
	CoshhLink:       nil,
	Location:        nil,
	StorageTemp:     "+4",
	IsArchived:      false,
	Hazards:         []string{"Explosive", "Flammable"},
	DBHazards:       nil,
}

var client = &http.Client{}

func TestPostChemical(t *testing.T) {
	jsonChemical, err := json.Marshal(chem)

	if err != nil {
		fmt.Println("error")
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/chemical", bytes.NewBuffer(jsonChemical))

	assert.Nil(t, err)

	fmt.Println("request error", err)

	response, err := client.Do(req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)

	bodyBytes, err := ioutil.ReadAll(response.Body)

	assert.Nil(t, err)

	var responseChemical chemical.Chemical
	err = json.Unmarshal(bodyBytes, &responseChemical)

	fmt.Println("result", string(bodyBytes))
	assert.Nil(t, err)

	assert.Equal(t, chem, responseChemical)
}

func TestGetChemical(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/chemicals", nil)
	assert.Nil(t, err)

	response, err := client.Do(req)
	assert.Nil(t, err)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	var responseChemicals []chemical.Chemical

	err = json.Unmarshal(bodyBytes, &responseChemicals)
	found := false
	for _, ch := range responseChemicals {
		if reflect.DeepEqual(ch, chem) {
			found = true
			break
		}
	}
	assert.True(t, found)
}
