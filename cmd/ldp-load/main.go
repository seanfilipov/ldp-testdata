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
	dbFlag := flag.String("db", "default-database",
		"database selected from configuration file")
	debugFlag := flag.Bool("debug", false, "enable debugging output")
	flag.Parse()
	opts := &ldpadmin.LoadOptions{
		Debug: *debugFlag,
	}

	config, err := ldputil.ReadConfig(*configFlag)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	extractDir := config.Get("extract", "dir")

	fmt.Printf("-- Starting load to database '%s'\n", *dbFlag)

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

	/*
		err = loadAll("groups", extractDir+"/groups.json", tx, opts)
		if err != nil {
			ldputil.PrintError(err)
			return
		}

		err = loadAll("users", extractDir+"/users.json", tx, opts)
		if err != nil {
			ldputil.PrintError(err)
			return
		}

		for x := 1; x <= 2; x++ {
			err = loadAll("tmp_loans_locations",
				extractDir+fmt.Sprintf("/circulation.loans.json.%v",
					x),
				tx, opts)
			if err != nil {
				ldputil.PrintError(err)
				return
			}
		}
	*/

	for x := 1; x <= 20; x++ {
		err = loadAll("loans",
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

	fmt.Printf("-- Load complete and committed to database '%s'\n", *dbFlag)
}

func loadAll(jtype string, filename string, tx *sql.Tx,
	opts *ldpadmin.LoadOptions) error {
	fmt.Printf("-- Loading from %s\n", filename)

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

	if jtype == "loans" {
		err = ldpadmin.LoadNEW(jtype, dec, tx, opts)
		if err != nil {
			return err
		}
		return nil
	}

	// Read and load array elements.
	for dec.More() {

		var i interface{}
		err := dec.Decode(&i)
		if err != nil {
			return err
		}

		err = ldpadmin.Load(jtype, i.(map[string]interface{}), tx,
			opts)
		if err != nil {
			return err
		}
	}

	return nil
}
