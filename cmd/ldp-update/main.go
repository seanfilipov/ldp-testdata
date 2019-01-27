package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/folio-org/ldp/cmd/internal/ldputil"
	"github.com/folio-org/ldp/ldpadmin"
)

func main() {
	configFlag := flag.String("config", "", "configuration file")
	dbFlag := flag.String("db", "ldp-database",
		"database selected from configuration file")
	debugFlag := flag.Bool("debug", false, "enable debugging output")
	flag.Parse()
	opts := &ldpadmin.UpdateOptions{
		Debug: *debugFlag,
	}

	config, err := ldputil.ReadConfig(*configFlag)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	extractDir := config.Get("extract", "dir")

	fmt.Printf("-- Starting update to database '%s'\n", *dbFlag)

	pgdb, err := ldputil.OpenDatabase(
		config.Get(*dbFlag, "host"),
		config.Get(*dbFlag, "port"),
		config.Get(*dbFlag, "user"),
		config.Get(*dbFlag, "password"),
		config.Get(*dbFlag, "dbname"))
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer pgdb.Close()

	tx, err := pgdb.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer tx.Rollback()

	err = updateAll("groups", extractDir+"/groups.json", tx, opts)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = updateAll("users", extractDir+"/users.json", tx, opts)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	for x := 1; x <= 20; x++ {
		err = updateAll("tmp_loans_locations",
			extractDir+fmt.Sprintf("/circulation.loans.json.%v",
				x),
			tx, opts)
		if err != nil {
			ldputil.PrintError(err)
			return
		}
	}

	for x := 1; x <= 20; x++ {
		err = updateAll("loans",
			extractDir+fmt.Sprintf("/loan-storage.loans.json.%v",
				x),
			tx, opts)
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

	fmt.Printf("-- Updates committed to database '%s'\n", *dbFlag)
}

func updateAll(jtype string, filename string, tx *sql.Tx,
	opts *ldpadmin.UpdateOptions) error {
	fmt.Printf("-- Updating from %s\n", filename)

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

	// Read and update array elements.
	for dec.More() {

		var i interface{}
		err := dec.Decode(&i)
		if err != nil {
			return err
		}

		err = ldpadmin.Update(jtype, i.(map[string]interface{}), tx,
			opts)
		if err != nil {
			return err
		}
	}

	return nil
}
