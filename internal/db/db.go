package db

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-envconfig"
	"gitlab.mdcatapult.io/informatics/coshh/coshh-api/internal/chemical"
)

var db *sqlx.DB

const SCHEMA = "coshh"

type Config struct {
	Port     int    `env:"PORT,required"`
	User     string `env:"USER,default=postgres,required"`
	Password string `env:"PASSWORD,required"`
	DbName   string `env:"DBNAME,required"`
	Host     string `env:"HOST,required"`
}

func Connect(host string) error {
	var (
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "informatics"
		retries  = 3
		schema   = SCHEMA
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
			fmt.Printf("Connected to database: %s, schema: %s\n", dbname, schema)

			break
		}

		fmt.Println(err)
		fmt.Println("Failed to connect to DB, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	return err
}

func SelectAllChemicals() ([]chemical.Chemical, error) {
	chemicals := make([]chemical.Chemical, 0)
	query := fmt.Sprintf(`
		SELECT 
			c.id,
			c.cas_number,
			c.chemical_name,
			c.chemical_number,
			c.matter_state,
			c.quantity,
			c.added,
			c.expiry,
			c.safety_data_sheet,
			c.coshh_link,
			c.lab_location,
		    c.cupboard,
			c.storage_temp,
			c.is_archived,
		    c.project_specific,
			string_agg(CAST(c2h.hazard AS VARCHAR(255)), ',') AS hazards 
		FROM %s.chemical c
		LEFT JOIN %s.chemical_to_hazard c2h ON c.id = c2h.id
		GROUP BY c.id`,
		SCHEMA, // DO NOT allow user input in raw SQL
		SCHEMA, // it can expose SQL injection
	)

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

func SelectAllCupboards() ([]string, error) {
	returnValue := make([]string, 0)

	query := fmt.Sprintf(`
		SELECT
		    DISTINCT COALESCE(c.cupboard, '')
		FROM %s.chemical c
	`, SCHEMA, // DO NOT allow user input in raw SQL
	)

	if err := db.Select(&returnValue, query); err != nil {
		fmt.Println("OI!!!!")
		return nil, err
	}

	sort.Strings(returnValue)
	return returnValue, nil
}

func UpdateChemical(chemical chemical.Chemical) error {
	query := fmt.Sprintf(`
	UPDATE %s.chemical
	SET 
		cas_number = :cas_number,
		chemical_name = :chemical_name,
		chemical_number = :chemical_number,
		matter_state = :matter_state,
		quantity = :quantity,
		added = :added,
		expiry = :expiry,
		safety_data_sheet = :safety_data_sheet,
		coshh_link = :coshh_link,
		lab_location = :lab_location,
	    cupboard = :cupboard,
		storage_temp = :storage_temp,
		is_archived = :is_archived,
		project_specific = :project_specific

	WHERE id = :id
`, SCHEMA, // DO NOT allow user input in raw sql
	)

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
	query := fmt.Sprintf(`INSERT INTO %s.chemical (
		cas_number,
		chemical_name,
		chemical_number,
		matter_state,
		quantity,
		added,
		expiry,
		safety_data_sheet,
		coshh_link,
		lab_location,
        cupboard,
		storage_temp,
		is_archived,
        project_specific
	)VALUES (
		:cas_number,
		:chemical_name,
		:chemical_number,
		:matter_state,
		:quantity,
		:added,
		:expiry,
		:safety_data_sheet,
		:coshh_link,
		:lab_location,
		:cupboard,
		:storage_temp,
		:is_archived,
		:project_specific
	) RETURNING id`,
		SCHEMA, // DO NOT allow user input in raw SQL
	)

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

func DeleteHazards(chemical chemical.Chemical) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`DELETE FROM %s.chemical_to_hazard WHERE id = $1;`, SCHEMA) // DO NOT allow user input in raw SQL
	_, err = tx.Exec(query, chemical.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func InsertHazards(chemical chemical.Chemical) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = insertHazards(tx, chemical, chemical.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertHazards(tx *sqlx.Tx, chemical chemical.Chemical, id int64) error {
	// chemicalToHazard represents a row in chemical_to_hazard
	type chemicalToHazard struct {
		Id     int64  `db:"id"`
		Hazard string `db:"hazard"`
	}

	chemicalToHazards := make([]chemicalToHazard, 0)

	query := fmt.Sprintf(`INSERT INTO %s.chemical_to_hazard (id, hazard) VALUES (:id, :hazard)`, SCHEMA) // DO NOT allow user input in raw SQL

	for _, hazard := range chemical.Hazards {
		chemicalToHazards = append(chemicalToHazards, chemicalToHazard{
			Id:     id,
			Hazard: hazard,
		})
	}

	_, err := tx.NamedExec(query, chemicalToHazards)
	return err
}
