package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/folio-org/ldp/cmd/internal/ldputil"
	_ "github.com/lib/pq"
)

func main() {
	config, err := ldputil.ReadConfig()
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	extractDir := config.Get("extract", "dir")

	db, err := ldputil.OpenDatabase(
		config.Get("ldp-database", "host"),
		config.Get("ldp-database", "port"),
		config.Get("ldp-database", "user"),
		config.Get("ldp-database", "password"),
		config.Get("ldp-database", "dbname"))
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = loadAll("groups", extractDir+"/groups.json", tx)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = loadAll("users", extractDir+"/users.json", tx)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	for x := 1; x <= 20; x++ {
		err = loadAll("tmp_loans_locations",
			extractDir+fmt.Sprintf("/circulation.loans.json.%v",
				x),
			tx)
		if err != nil {
			ldputil.PrintError(err)
			return
		}
	}

	for x := 1; x <= 20; x++ {
		err = loadAll("loans",
			extractDir+fmt.Sprintf("/loan-storage.loans.json.%v",
				x),
			tx)
		if err != nil {
			ldputil.PrintError(err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		ldputil.PrintError(err)
		return
	}
}

func loadAll(jtype string, filename string, tx *sql.Tx) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	// Skip past first tokens.
	for x := 0; x < 3; x++ {
		_, err := dec.Token()
		if err != nil {
			return err
		}
	}

	// Read and load array elements.
	count := 0
	for dec.More() {

		var i interface{}
		err := dec.Decode(&i)
		if err != nil {
			return err
		}

		err = load(tx, jtype, i.(map[string]interface{}))
		if err != nil {
			return err
		}

		count++
	}

	fmt.Printf("%s %v %s\n", jtype, count, filename)

	return nil
}

func load(tx *sql.Tx, jtype string, j map[string]interface{}) error {
	switch jtype {
	case "groups":
		err := loadGroup(tx, jtype, j)
		if err != nil {
			return err
		}
	case "users":
		err := loadUser(tx, jtype, j)
		if err != nil {
			return err
		}
	case "tmp_loans_locations":
		err := loadTmpLoanLocation(tx, jtype, j)
		if err != nil {
			return err
		}
	case "loans":
		err := loadLoan(tx, jtype, j)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown type \"%v\"", jtype)
	}

	return nil
}

func loadGroup(tx *sql.Tx, jtype string, j map[string]interface{}) error {

	id := j["id"].(string)
	groupName := j["group"].(string)
	description := j["desc"].(string)

	_, err := tx.Exec(
		"INSERT INTO groups AS g "+
			"(id, group_name, description) "+
			"VALUES ($1, $2, $3) "+
			"ON CONFLICT (id) DO "+
			"UPDATE SET group_name = EXCLUDED.group_name, "+
			"description = EXCLUDED.description "+
			"WHERE g.group_name <> EXCLUDED.group_name OR "+
			"g.description <> EXCLUDED.description",
		id, groupName, description)
	if err != nil {
		return err
	}

	return nil
}

func loadUser(tx *sql.Tx, jtype string, j map[string]interface{}) error {

	id := j["id"].(string)
	username := j["username"].(string)
	barcode := j["barcode"].(string)
	userType := j["type"].(string)
	active := j["active"].(string)
	patronGroupId := j["patronGroup"].(string)

	_, err := tx.Exec(
		"INSERT INTO users AS u "+
			"(id, username, barcode, user_type, active, "+
			"patron_group_id) "+
			"VALUES ($1, $2, $3, $4, $5, $6) "+
			"ON CONFLICT (id) DO "+
			"UPDATE SET username = EXCLUDED.username, "+
			"barcode = EXCLUDED.barcode, "+
			"user_type = EXCLUDED.user_type, "+
			"active = EXCLUDED.active, "+
			"patron_group_id = EXCLUDED.patron_group_id "+
			"WHERE u.username <> EXCLUDED.username OR "+
			"u.barcode <> EXCLUDED.barcode OR "+
			"u.user_type <> EXCLUDED.user_type OR "+
			"u.active <> EXCLUDED.active OR "+
			"u.patron_group_id <> EXCLUDED.patron_group_id",
		id, username, barcode, userType, active, patronGroupId)
	if err != nil {
		return err
	}

	return nil
}

func loadTmpLoanLocation(tx *sql.Tx, jtype string,
	j map[string]interface{}) error {

	loanId := j["id"].(string)
	item := j["item"].(map[string]interface{})
	location := item["location"]
	locationName := location.(map[string]interface{})["name"].(string)

	_, err := tx.Exec(
		"INSERT INTO tmp_loans_locations AS t "+
			"(loan_id, location_name) "+
			"VALUES ($1, $2) "+
			"ON CONFLICT (loan_id) DO "+
			"UPDATE SET location_name = EXCLUDED.location_name "+
			"WHERE t.location_name <> EXCLUDED.location_name",
		loanId, locationName)
	if err != nil {
		return err
	}

	return nil
}

func loadLoan(tx *sql.Tx, jtype string, j map[string]interface{}) error {

	id := j["id"].(string)
	userId := j["userId"].(string)
	itemId := j["itemId"].(string)
	action := j["action"].(string)

	status := j["status"].(map[string]interface{})
	statusName := status["name"].(string)

	loanDateStr := j["loanDate"].(string)
	dueDateStr := j["dueDate"].(string)

	layout := "2006-01-02T15:04:05Z"
	loanDate, _ := time.Parse(layout, loanDateStr)
	dueDate, _ := time.Parse(layout, dueDateStr)

	_, err := tx.Exec(
		"INSERT INTO loans AS l "+
			"(id, user_id, item_id, action, status_name, "+
			"loan_date, due_date) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7) "+
			"ON CONFLICT (id) DO "+
			"UPDATE SET user_id = EXCLUDED.user_id, "+
			"item_id = EXCLUDED.item_id, "+
			"action = EXCLUDED.action, "+
			"status_name = EXCLUDED.status_name, "+
			"loan_date = EXCLUDED.loan_date, "+
			"due_date = EXCLUDED.due_date "+
			"WHERE l.user_id <> EXCLUDED.user_id OR "+
			"l.item_id <> EXCLUDED.item_id OR "+
			"l.action <> EXCLUDED.action OR "+
			"l.status_name <> EXCLUDED.status_name OR "+
			"l.loan_date <> EXCLUDED.loan_date OR "+
			"l.due_date <> EXCLUDED.due_date",
		id, userId, itemId, action, statusName, loanDate, dueDate)
	if err != nil {
		return err
	}

	return nil
}
