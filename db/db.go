package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-envconfig"
	"gitlab.mdcatapult.io/informatics/software-engineering/coshh/chemical"
	"strings"
	"time"
)

var db *sqlx.DB

type Config struct {
	Port     int    `env:"PORT,required"`
	User string `env:"USER,default=postgres,required"`
	Password string `env:"PASSWORD,required"`
	DbName string `env:"DBNAME,required"`
	Host string `env:"HOST,required"`
}

func Connect(host string) error {
	var (
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "informatics"
		retries  = 3
		schema = "coshh"
	)

	ctx := context.Background()
	var config Config

	if err := envconfig.Process(ctx, &config); err != nil {
		fmt.Println("Env vars unset or incorrect, using default config")
	} else {
		fmt.Println("Using config from env vars")
		host = config.Host
		port = config.Port
		user = config.User
		password = config.Password
		dbname = config.DbName
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	for i := 1; i < retries; i++ {
		db, err = sqlx.Connect("postgres", psqlInfo)
		if err == nil {
			_, err = db.Exec(fmt.Sprintf("set search_path=%s", schema))
			if err != nil {
				fmt.Printf("Failed to set search path to schema: %s\n", schema)
				return err
			}
			break
		}

		fmt.Println("Failed to connect to DB, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	return err
}

func SelectAllChemicals() ([]chemical.Chemical, error) {
	var chemicals []chemical.Chemical
	query := `
		SELECT 
			c.id,
			c.cas_number,
			c.chemical_name,
			c.photo_path,
			c.matter_state,
			c.quantity,
			c.added,
			c.expiry,
			c.safety_data_sheet,
			c.coshh_link,
			c.lab_location,
			c.storage_temp,
			c.is_archived,
			string_agg(CAST(c2h.hazard AS VARCHAR(255)), ',') AS hazards 
		FROM chemical c 
		LEFT JOIN chemical_to_hazard c2h ON c.id = c2h.id 
		GROUP BY c.id`

	if err := db.Select(&chemicals, query); err != nil {
		return nil, err
	}

	for i := range chemicals {
		if chemicals[i].DBHazards != nil {
			chemicals[i].Hazards = strings.Split(*chemicals[i].DBHazards, ",")
		}
	}

	return chemicals, nil
}

func UpdateChemical(chemical chemical.Chemical) error {
	query := `
	UPDATE chemical 
	SET 
		cas_number = :cas_number,
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
	WHERE id = :id
`

	_, err := db.NamedExec(query, chemical)
	return err
}

func InsertChemical(chemical chemical.Chemical) (id int64, err error) {

	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}

	id, err = insertChemical(tx, chemical)
	if err != nil {
		return 0, err
	}

	if err := insertHazards(tx, chemical, id); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func insertChemical(tx *sqlx.Tx, chemical chemical.Chemical) (id int64, err error) {
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
	) RETURNING id`

	rows, err := tx.NamedQuery(query, chemical)
	if err != nil {
		return
	}

	rows.Next()
	if err := rows.Scan(&id); err != nil {
		return 0, err
	}

	if err := rows.Close(); err != nil {
		return 0, err
	}

	return

}

func insertHazards(tx *sqlx.Tx, chemical chemical.Chemical, id int64) error {
	// chemicalToHazard represents a row in chemical_to_hazard
	type chemicalToHazard struct {
		Id     int64  `db:"id"`
		Hazard string `db:"hazard"`
	}

	chemicalToHazards := make([]chemicalToHazard, 0)

	query := `INSERT INTO chemical_to_hazard (id, hazard) VALUES (:id, :hazard)`

	for _, hazard := range chemical.Hazards {
		chemicalToHazards = append(chemicalToHazards, chemicalToHazard{
			Id:     id,
			Hazard: hazard,
		})
	}

	_, err := tx.NamedExec(query, chemicalToHazards)
	return err
}
