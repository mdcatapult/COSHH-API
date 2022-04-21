package chemical

import "time"

type Chemical struct {
	CasNumber       string    `json:"casNumber" db:"cas_number"`
	Name            string    `json:"name" db:"chemical_name"`
	PhotoPath       string    `json:"photoPath" db:"photo_path"`
	MatterState     string    `json:"matterState" db:"matter_state"`
	Quantity        string    `json:"quantity" db:"quantity"`
	Added           time.Time `json:"added" db:"added"`
	Expiry          time.Time `json:"expiry" db:"expiry"`
	SafetyDataSheet string    `json:"safetyDataSheet" db:"safety_data_sheet"`
	CoshhLink       *string   `json:"coshhLink" db:"coshh_link"`
	Location        *string   `json:"location" db:"lab_location"`
	StorageTemp     string    `json:"storageTemp" db:"storage_temp"`
	IsArchived      bool      `json:"isArchived" db:"is_archived"`
	Hazards         []string  `json:"hazards" db:"-"`
	DBHazards       *string   `json:"-" db:"hazards"`
}
