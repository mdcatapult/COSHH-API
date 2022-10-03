package chemical

import "time"

type Chemical struct {
	Id              int64      `json:"id" db:"id"`
	CasNumber       *string    `json:"casNumber" db:"cas_number"`
	Name            string     `json:"name" db:"chemical_name"`
	ChemicalNumber  *string    `json:"chemicalNumber" db:"chemical_number"`
	MatterState     *string    `json:"matterState" db:"matter_state"`
	Quantity        *string    `json:"quantity" db:"quantity"`
	Added           *time.Time `json:"added" db:"added"`
	Expiry          *time.Time `json:"expiry" db:"expiry"`
	SafetyDataSheet *string    `json:"safetyDataSheet" db:"safety_data_sheet"`
	CoshhLink       *string    `json:"coshhLink" db:"coshh_link"`
	Location        *string    `json:"location" db:"lab_location"`
	Cupboard        *string    `json:"cupboard" db:"cupboard"`
	StorageTemp     string     `json:"storageTemp" db:"storage_temp"`
	IsArchived      bool       `json:"isArchived" db:"is_archived"`
	Hazards         []string   `json:"hazards" db:"-"`
	DBHazards       *string    `json:"-" db:"hazards"`
}
