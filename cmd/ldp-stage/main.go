package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lib/pq"
)

func main() {
	config, err := readConfig()
	if err != nil {
		printError(err)
		return
	}

	extractDir := config.Get("extract", "dir")

	stagedb, err := openDatabase(
		config.Get("stage-database", "host"),
		config.Get("stage-database", "port"),
		config.Get("stage-database", "user"),
		config.Get("stage-database", "password"),
		config.Get("stage-database", "dbname"))
	if err != nil {
		printError(err)
		return
	}
	defer stagedb.Close()

	stagetx, err := stagedb.Begin()
	if err != nil {
		printError(err)
		return
	}
	defer stagetx.Rollback()

	_, err = stagetx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		printError(err)
		return
	}

	err = stage("groups", extractDir+"/groups.json", stagetx)
	if err != nil {
		printError(err)
		return
	}

	err = stage("users", extractDir+"/users.json", stagetx)
	if err != nil {
		printError(err)
		return
	}

	for x := 1; x <= 20; x++ {
		err = stage("tmp_locations",
			extractDir+fmt.Sprintf("/circulation.loans.json.%v",
				x),
			stagetx)
		if err != nil {
			printError(err)
			return
		}
	}

	for x := 1; x <= 20; x++ {
		err = stage("loans",
			extractDir+fmt.Sprintf("/loan-storage.loans.json.%v",
				x),
			stagetx)
		if err != nil {
			printError(err)
			return
		}
	}

	err = stagetx.Commit()
	if err != nil {
		printError(err)
		return
	}
	fmt.Printf("Updates committed to staging database\n")
}

func stage(jtype string, filename string, tx *sql.Tx) error {

	fmt.Printf("Reading %s from %s\n", jtype, filename)

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

	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
