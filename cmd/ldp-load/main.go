package main

import (
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
	opts := &ldpadmin.LoaderOptions{
		Debug: *debugFlag,
	}

	config, err := ldputil.ReadConfig(*configFlag)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	extractDir := config.Get("extract", "dir")

	db, err := ldputil.OpenDatabase(
		config.Get(*dbFlag, "host"),
		config.Get(*dbFlag, "port"),
		config.Get(*dbFlag, "user"),
		config.Get(*dbFlag, "password"),
		config.Get(*dbFlag, "dbname"))
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer db.Close()

	ld, err := ldpadmin.NewLoader(db, opts)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	fmt.Printf("-- Starting load to database '%s'\n", *dbFlag)

	err = loadFile("groups", extractDir+"/groups.json", ld)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = loadFile("users", extractDir+"/users.json", ld)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	for x := 1; x <= 20; x++ {
		err = loadFile("tmp_loans_locations",
			extractDir+fmt.Sprintf("/circulation.loans.json.%v",
				x),
			ld)
		if err != nil {
			ldputil.PrintError(err)
			return
		}
	}

	for x := 1; x <= 20; x++ {
		err = loadFile("loans",
			extractDir+fmt.Sprintf("/loan-storage.loans.json.%v",
				x),
			ld)
		if err != nil {
			ldputil.PrintError(err)
			return
		}
	}

	ld.Close()

	fmt.Printf("-- Load complete in database '%s'\n", *dbFlag)
}

func loadFile(jtype string, filename string, ld *ldpadmin.Loader) error {
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

	err = ld.Load(jtype, dec)
	if err != nil {
		return err
	}

	return nil
}
