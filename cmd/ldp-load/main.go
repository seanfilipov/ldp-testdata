package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/nassibnassar/goconfig/ini"
)

func main() {
	config, err := readConfig()
	if err != nil {
		printError(err)
		return
	}

	stage1db, err := openDatabase(
		config.Get("stage-database", "host"),
		config.Get("stage-database", "port"),
		config.Get("stage-database", "user"),
		config.Get("stage-database", "password"),
		config.Get("stage-database", "dbname"))
	if err != nil {
		printError(err)
		return
	}
	defer stage1db.Close()
	stage1tx, err := stage1db.Begin()
	if err != nil {
		printError(err)
		return
	}
	defer stage1tx.Rollback()
	_, err = stage1tx.Exec(
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		printError(err)
		return
	}

	stage2db, err := openDatabase(
		config.Get("stage-database", "host"),
		config.Get("stage-database", "port"),
		config.Get("stage-database", "user"),
		config.Get("stage-database", "password"),
		config.Get("stage-database", "dbname"))
	if err != nil {
		printError(err)
		return
	}
	defer stage2db.Close()
	stage2tx, err := stage2db.Begin()
	if err != nil {
		printError(err)
		return
	}
	defer stage2tx.Rollback()
	_, err = stage2tx.Exec(
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		printError(err)
		return
	}

	ldpdb, err := openDatabase(
		config.Get("ldp-database", "host"),
		config.Get("ldp-database", "port"),
		config.Get("ldp-database", "user"),
		config.Get("ldp-database", "password"),
		config.Get("ldp-database", "dbname"))
	if err != nil {
		printError(err)
		return
	}
	defer ldpdb.Close()
	ldptx, err := ldpdb.Begin()
	if err != nil {
		printError(err)
		return
	}
	defer ldptx.Rollback()
	_, err = ldptx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		printError(err)
		return
	}

	tx := &dbtx{
		stage1: stage1tx,
		stage2: stage2tx,
		ldp:    ldptx,
		locset: make(map[string]string),
	}

	err = loadAllStage(tx)
	if err != nil {
		printError(err)
		return
	}

	err = ldptx.Commit()
	if err != nil {
		printError(err)
		return
	}

	err = stage1tx.Commit()
	if err != nil {
		printError(err)
		return
	}

	err = stage2tx.Commit()
	if err != nil {
		printError(err)
		return
	}
}

type dbtx struct {
	stage1 *sql.Tx           // Outer select loop
	stage2 *sql.Tx           // Secondary lookups
	ldp    *sql.Tx           // Target database
	locset map[string]string // Temporary memory of locations
}

func loadAllStage(tx *dbtx) error {

	rows, err := tx.stage1.Query(
		"SELECT id, jtype, jid, j " +
			"FROM stage " +
			"ORDER BY id")
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int = 0
	var sr DataUnit
	for rows.Next() {

		var js string
		err := rows.Scan(&sr.Id, &sr.Jtype, &sr.Jid, &js)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(js), &sr.J)
		if err != nil {
			return err
		}

		err = loadStageRow(&sr, tx)
		if err != nil {
			return err
		}

		count++
		if count%100000 == 0 {
			fmt.Println(count)
		}
	}

	/*
		if sr.Id > 0 {
			err = deleteStageRows(tx, sr.Id)
			if err != nil {
				return err
			}
		}
	*/

	return nil
}

func loadStageRow(sr *DataUnit, tx *dbtx) error {

	if sr.Jtype == "users" {
		err := loadUser(sr, tx)
		if err != nil {
			return err
		}
	}

	if sr.Jtype == "tmp_locations" {
		err := loadTmpLocation(sr, tx)
		if err != nil {
			return err
		}
	}

	if sr.Jtype == "loans" {
		err := loadLoan(sr, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadTmpLocation(sr *DataUnit, tx *dbtx) error {

	mockid, name := mockLocation(sr)

	if tx.locset[mockid] == "" {
		_, err := tx.ldp.Exec(
			"INSERT INTO locations (location_id, "+
				"location_name) "+
				"VALUES ($1, $2) "+
				"ON CONFLICT (location_id) DO "+
				"UPDATE SET location_name = "+
				"EXCLUDED.location_name",
			mockid, name)
		if err != nil {
			return err
		}
		tx.locset[mockid] = mockid
	}

	return nil
}

func loadLoan(sr *DataUnit, tx *dbtx) error {

	loanId := sr.Jid
	userId := sr.J["userId"].(string)
	loanDate := sr.J["loanDate"].(string)

	locsr, err := lookup(tx, "tmp_locations", loanId)
	if err != nil {
		return err
	}
	if locsr == nil {
		return fmt.Errorf("Loan %v is missing location data",
			loanId)
	}

	mockid, _ := mockLocation(locsr)

	_, err = tx.ldp.Exec(
		"INSERT INTO loans (loan_id, user_id, location_id, "+
			"loan_date) "+
			"VALUES ($1, $2, $3, $4) "+
			"ON CONFLICT (loan_id) DO "+
			"UPDATE SET user_id = EXCLUDED.user_id, "+
			"location_id = EXCLUDED.location_id, "+
			"loan_date = EXCLUDED.loan_date",
		loanId, userId, mockid, loanDate)
	if err != nil {
		return err
	}

	return nil
}

func mockLocation(sr *DataUnit) (string, string) {
	item := sr.J["item"]
	location := item.(map[string]interface{})["location"]
	name := location.(map[string]interface{})["name"].(string)
	mockid := "id-" +
		strings.Replace(strings.ToLower(name), " ", "-", -1)
	return mockid, name
}

func loadUser(sr *DataUnit, tx *dbtx) error {

	groupJid := sr.J["patronGroup"].(string)

	group, err := lookup(tx, "groups", groupJid)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("User %v references unknown group %v",
			sr.Jid, groupJid)
	}

	_, err = tx.ldp.Exec(
		"INSERT INTO users (user_id, group_name) "+
			"VALUES ($1, $2) "+
			"ON CONFLICT (user_id) DO "+
			"UPDATE SET group_name = EXCLUDED.group_name",
		sr.Jid,
		group.J["group"].(string))
	if err != nil {
		return err
	}

	return nil
}

func lookup(tx *dbtx, jtype string, jid string) (*DataUnit, error) {
	rows, err := tx.stage2.Query(
		"SELECT id, jtype, jid, j "+
			"FROM stage "+
			"WHERE jtype = $1 AND jid = $2 "+
			"ORDER BY id "+
			"LIMIT 1",
		jtype, jid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	sr, err := scan(rows)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func scan(rows *sql.Rows) (*DataUnit, error) {
	var sr DataUnit
	var js string
	err := rows.Scan(&sr.Id, &sr.Jtype, &sr.Jid, &js)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(js), &sr.J)
	if err != nil {
		return nil, err
	}

	return &sr, nil
}

func deleteStageRows(tx *dbtx, id int64) error {

	_, err := tx.stage1.Exec(
		"DELETE FROM stage WHERE id <= $1", id)
	if err != nil {
		return err
	}

	return nil
}

type DataUnit struct {
	Id    int64
	Jtype string
	Jid   string
	J     map[string]interface{}
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
}

func readConfig() (*ini.Config, error) {
	f := os.Getenv("LDP_CONFIG_FILE")
	if f == "" {
		return ini.NewConfig(), nil
	}
	c, err := ini.NewConfigFile(f)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading configuration file: %v", err)
	}
	return c, nil
}

func openDatabase(host, port, user, password, dbname string) (*sql.DB, error) {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s "+
			"sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to test for connection errors.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
