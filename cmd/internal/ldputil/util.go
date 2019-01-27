package ldputil

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/nassibnassar/goconfig/ini"
)

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
}

func ReadConfig(filename string) (*ini.Config, error) {
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

func OpenDatabase(host, port, user, password, dbname string) (*sql.DB, error) {

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
