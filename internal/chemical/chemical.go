/*
* Copyright $today.year Medicines Discovery Catapult
* Licensed under the Apache License, Version 2.0 (the "Licence");
* you may not use this file except in compliance with the Licence.
* You may obtain a copy of the Licence at
*     http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the Licence is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the Licence for the specific language governing permissions and
* limitations under the Licence.
 */

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
	Owner           *string    `json:"owner" db:"chemical_owner"`
	LastUpdatedBy   *string    `json:"lastUpdatedBy" db:"last_updated_by"`
}
