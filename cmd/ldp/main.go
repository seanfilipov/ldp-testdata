package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/folio-org/ldp/ldpadmin"
	_ "github.com/lib/pq"
	"github.com/nassibnassar/goconfig/ini"
)

func main() {
	initFlag := flag.Bool("init", false,
		"initialize a database with the LDP schema")
	loadFlag := flag.Bool("load", false, "load data into a database")
	configFlag := flag.String("config", "", "configuration file")
	dbFlag := flag.String("db", "default-database",
		"database selected from configuration file")
	debugFlag := flag.Bool("debug", false, "enable debugging output")
	dirFlag := flag.String("dir", "", "source data directory")
	flag.Parse()

	if !*initFlag && !*loadFlag {
		ldputil.printError(fmt.Errorf(
			"no command flag specified (use -init or -load)"))
		return
	}

	if *initFlag && *loadFlag {
		ldputil.printError(fmt.Errorf(
			"-init and -load cannot be used together"))
		return
	}

	if *initFlag {
		err := cmdInit(configFlag, dbFlag, debugFlag)
		if err != nil {
			ldputil.printError(err)
			return
		}
		return
	}

	if *loadFlag {
		err := cmdLoad(configFlag, dbFlag, debugFlag, dirFlag)
		if err != nil {
			ldputil.printError(err)
			return
		}
		return
	}
}

func cmdInit(configFlag *string, dbFlag *string, debugFlag *bool) error {
	opts := &ldpadmin.InitOptions{
		Debug: *debugFlag,
	}
	config, err := ldputil.readConfig(*configFlag)
	if err != nil {
		return err
	}
	db, err := ldputil.openDatabase(
		config.Get(*dbFlag, "host"),
		config.Get(*dbFlag, "port"),
		config.Get(*dbFlag, "user"),
		config.Get(*dbFlag, "password"),
		config.Get(*dbFlag, "dbname"))
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Printf("-- Starting initialization of database '%s'\n", *dbFlag)
	err = ldpadmin.Initialize(db, opts)
	if err != nil {
		return err
	}
	fmt.Printf("-- Initialization complete in database '%s'\n", *dbFlag)
	return nil
}

func cmdLoad(configFlag *string, dbFlag *string, debugFlag *bool,
	dirFlag *string) error {
	opts := &ldpadmin.LoaderOptions{
		Debug: *debugFlag,
	}

	config, err := ldputil.readConfig(*configFlag)
	if err != nil {
		return err
	}

	extractDir := *dirFlag

	db, err := ldputil.openDatabase(
		config.Get(*dbFlag, "host"),
		config.Get(*dbFlag, "port"),
		config.Get(*dbFlag, "user"),
		config.Get(*dbFlag, "password"),
		config.Get(*dbFlag, "dbname"))
	if err != nil {
		return err
	}
	defer db.Close()

	ld, err := ldpadmin.NewLoader(db, opts)
	if err != nil {
		return err
	}

	fmt.Printf("-- Starting load to database '%s'\n", *dbFlag)

	err = loadFile("groups", extractDir+"/groups.json", ld)
	if err != nil {
		return err
	}

	err = loadFile("users", extractDir+"/users.json", ld)
	if err != nil {
		return err
	}

	for x := 1; x <= 20; x++ {
		err = loadFile("tmp_loans_locations",
			extractDir+fmt.Sprintf("/circulation.loans.json.%v",
				x),
			ld)
		if err != nil {
			return err
		}
	}

	for x := 1; x <= 20; x++ {
		err = loadFile("loans",
			extractDir+fmt.Sprintf("/loan-storage.loans.json.%v",
				x),
			ld)
		if err != nil {
			return err
		}
	}

	ld.Close()

	fmt.Printf("-- Load complete in database '%s'\n", *dbFlag)
	return nil
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

func printError(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
}

func readConfig(filename string) (*ini.Config, error) {
	var fn string
	if filename != "" {
		fn = filename
	} else {
		fn = os.Getenv("LDP_CONFIG_FILE")
		if fn == "" {
			return ini.NewConfig(), nil
		}
	}
	c, err := ini.NewConfigFile(fn)
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
