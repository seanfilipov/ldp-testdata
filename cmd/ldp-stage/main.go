package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/folio-org/ldp/cmd/internal/ldputil"
	"github.com/lib/pq"
)

func main() {
	config, err := ldputil.ReadConfig()
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	extractDir := config.Get("extract", "dir")

	stagedb, err := ldputil.OpenDatabase(
		config.Get("stage-database", "host"),
		config.Get("stage-database", "port"),
		config.Get("stage-database", "user"),
		config.Get("stage-database", "password"),
		config.Get("stage-database", "dbname"))
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer stagedb.Close()

	stagetx, err := stagedb.Begin()
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer stagetx.Rollback()

	_, err = stagetx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = stage("groups", extractDir+"/groups.json", stagetx)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = stage("users", extractDir+"/users.json", stagetx)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	for x := 1; x <= 20; x++ {
		err = stage("tmp_loans_locations",
			extractDir+fmt.Sprintf("/circulation.loans.json.%v",
				x),
			stagetx)
		if err != nil {
			ldputil.PrintError(err)
			return
		}
	}

	for x := 1; x <= 20; x++ {
		err = stage("loans",
			extractDir+fmt.Sprintf("/loan-storage.loans.json.%v",
				x),
			stagetx)
		if err != nil {
			ldputil.PrintError(err)
			return
		}
	}

	err = stagetx.Commit()
	if err != nil {
		ldputil.PrintError(err)
		return
	}
}

func stage(jtype string, filename string, tx *sql.Tx) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	stmt, err := tx.Prepare(pq.CopyInSchema(
		"public", "stage",
		"jtype", "jid", "j"))
	if err != nil {
		return err
	}
	defer stmt.Close()

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

		jid := i.(map[string]interface{})["id"].(string)

		j, err := json.Marshal(i)
		if err != nil {
			fmt.Println("error:", err)
		}

		_, err = stmt.Exec(jtype, jid, string(j))
		if err != nil {
			return err
		}

		count++
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	fmt.Printf("%s %v %s\n", jtype, count, filename)

	return nil
}
