package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/folio-org/ldp/cmd/internal/ldputil"
	_ "github.com/lib/pq"
)

const missingDataString string = "NOT AVAILABLE"

func main() {
	config, err := ldputil.ReadConfig()
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	stage1db, err := ldputil.OpenDatabase(
		config.Get("stage-database", "host"),
		config.Get("stage-database", "port"),
		config.Get("stage-database", "user"),
		config.Get("stage-database", "password"),
		config.Get("stage-database", "dbname"))
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer stage1db.Close()
	stage1tx, err := stage1db.Begin()
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer stage1tx.Rollback()
	_, err = stage1tx.Exec(
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	stage2db, err := ldputil.OpenDatabase(
		config.Get("stage-database", "host"),
		config.Get("stage-database", "port"),
		config.Get("stage-database", "user"),
		config.Get("stage-database", "password"),
		config.Get("stage-database", "dbname"))
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer stage2db.Close()
	stage2tx, err := stage2db.Begin()
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer stage2tx.Rollback()
	_, err = stage2tx.Exec(
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	ldpdb, err := ldputil.OpenDatabase(
		config.Get("ldp-database", "host"),
		config.Get("ldp-database", "port"),
		config.Get("ldp-database", "user"),
		config.Get("ldp-database", "password"),
		config.Get("ldp-database", "dbname"))
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer ldpdb.Close()
	ldptx, err := ldpdb.Begin()
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer ldptx.Rollback()
	_, err = ldptx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		ldputil.PrintError(err)
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
		ldputil.PrintError(err)
		return
	}

	// TODO Loop through all dim.* tables to see if any missing data
	// are now available, and if so, update the records.

	err = ldptx.Commit()
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = stage1tx.Commit()
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = stage2tx.Commit()
	if err != nil {
		ldputil.PrintError(err)
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
	var maxId int64 = 0
	for rows.Next() {

		var du DataUnit

		var js string
		err := rows.Scan(&du.Id, &du.Jtype, &du.Jid, &js)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(js), &du.J)
		if err != nil {
			return err
		}

		err = loadDataUnit(&du, tx)
		if err != nil {
			return err
		}

		maxId = du.Id
		count++
		if count%100000 == 0 {
			fmt.Println(count)
		}
	}

	if maxId > 0 {
		/*
			err = deleteFromStage(tx, maxId)
			if err != nil {
				return err
			}
		*/
	}

	return nil
}

func loadDataUnit(du *DataUnit, tx *dbtx) error {

	switch du.Jtype {
	case "groups":
		err := storeGroup(du, tx)
		if err != nil {
			return err
		}
	case "users":
		err := loadUser(du, tx)
		if err != nil {
			return err
		}
		/*
			case "loans":
				err := loadLoan(du, tx)
				if err != nil {
					return err
				}
			case "tmp_loans_locations":
				err := updateMirror(du, tx)
				if err != nil {
					return err
				}
				err = loadTmpLocation(du, tx)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unknown data unit type: %v", du.Jtype)
		*/
	}

	return nil
}

func storeGroup(du *DataUnit, tx *dbtx) error {

	//j, err := json.marshal(du.j)
	//if err != nil {
	//fmt.println("error:", err)
	//}

	id := du.J["id"].(string)
	name := du.J["group"].(string)
	description := du.J["desc"].(string)

	_, err := tx.stage2.Exec(
		"INSERT INTO denorm.groups AS g (id, name, description) "+
			"VALUES ($1, $2, $3) "+
			"ON CONFLICT (id) DO "+
			"UPDATE SET name = EXCLUDED.name, "+
			"description = EXCLUDED.description "+
			"WHERE g.name <> EXCLUDED.name OR "+
			"g.description <> EXCLUDED.description",
		id, name, description)
	if err != nil {
		return err
	}

	return nil
}

/*
func updateMirror(du *DataUnit, tx *dbtx) error {

	table, err := tablename(du.jtype)
	if err != nil {
		return err
	}

	j, err := json.marshal(du.j)
	if err != nil {
		fmt.println("error:", err)
	}

	_, err = tx.stage2.exec(
		"insert into "+table+" (jid, j) "+
			"values ($1, $2) "+
			"on conflict (jid) do "+
			"update set j = excluded.j",
		du.jid, string(j))
	if err != nil {
		return err
	}

	return nil
}
*/

func loadTmpLocation(du *DataUnit, tx *dbtx) error {

	mockid, name := mockLocation(du.J)

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

func loadLoan(du *DataUnit, tx *dbtx) error {

	loanId := du.Jid
	userId := du.J["userId"].(string)
	loanDate := du.J["loanDate"].(string)

	loanlocJ, err := lookupMirror(tx, "tmp_loans_locations", loanId)
	if err != nil {
		return err
	}
	if loanlocJ == nil {
		return fmt.Errorf("Loan %v is missing location data",
			loanId)
	}

	mockid, _ := mockLocation(loanlocJ)

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

func mockLocation(j map[string]interface{}) (string, string) {
	item := j["item"]
	location := item.(map[string]interface{})["location"]
	name := location.(map[string]interface{})["name"].(string)
	mockid := "id-" +
		strings.Replace(strings.ToLower(name), " ", "-", -1)
	return mockid, name
}

func loadUser(du *DataUnit, tx *dbtx) error {

	userName := du.J["username"].(string)
	active := du.J["active"].(string)

	groupJid := du.J["patronGroup"].(string)

	//groupJ, err := lookupMirror(tx, "groups", groupJid)
	groupName, groupDesc, found, err := recallGroup(tx, groupJid)
	if err != nil {
		return err
	}
	if !found {
		// TODO Store a record of the incomplete user and
		// unknown group ID in dim.users.
		//err := storeUser(...)
		groupName = missingDataString
		groupDesc = missingDataString
	}
	_ = groupDesc

	_, err = tx.ldp.Exec(
		"INSERT INTO users AS u "+
			"(user_id, user_name, active, group_name) "+
			"VALUES ($1, $2, $3, $4) "+
			"ON CONFLICT (user_id) DO "+
			"UPDATE SET user_name = EXCLUDED.user_name, "+
			"active = EXCLUDED.active, "+
			"group_name = EXCLUDED.group_name "+
			"WHERE u.user_name <> EXCLUDED.user_name OR "+
			"u.active <> EXCLUDED.active OR "+
			"u.group_name <> EXCLUDED.group_name",
		du.Jid,
		userName,
		active,
		groupName)
	if err != nil {
		return err
	}

	return nil
}

func recallGroup(tx *dbtx, id string) (string, string, bool, error) {

	rows, err := tx.stage2.Query(
		"SELECT name, description "+
			"FROM denorm.groups "+
			"WHERE id = $1",
		id)
	if err != nil {
		return "", "", false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", "", false, nil
	}

	var name string
	var description string
	err = rows.Scan(&name, &description)
	if err != nil {
		return "", "", false, err
	}

	return name, description, true, nil
}

func lookupMirror(tx *dbtx, jtype string, jid string) (map[string]interface{},
	error) {

	table, err := tableName(jtype)
	if err != nil {
		return nil, err
	}

	rows, err := tx.stage2.Query(
		"SELECT j "+
			"FROM "+table+" "+
			"WHERE jid = $1",
		jid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var j map[string]interface{}
	var js string
	err = rows.Scan(&js)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(js), &j)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func deleteFromStage(tx *dbtx, id int64) error {

	_, err := tx.stage1.Exec(
		"DELETE FROM stage WHERE id <= $1", id)
	if err != nil {
		return err
	}

	return nil
}

func tableName(jtype string) (string, error) {
	switch jtype {
	case "groups",
		//"loans",
		//"users",
		"tmp_loans_locations":
		return jtype, nil
	default:
		return "", fmt.Errorf("Data unit type \"%s\" is unknown", jtype)
	}
}

type DataUnit struct {
	Id    int64
	Jtype string
	Jid   string
	J     map[string]interface{}
}
