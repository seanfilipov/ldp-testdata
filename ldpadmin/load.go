package ldpadmin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Loader struct {
	db     *sql.DB
	tx     *sql.Tx
	locktx *sql.Tx
	opts   *LoaderOptions
}

type LoaderOptions struct {
	// Debug enables debugging output if set to true.
	Debug bool
}

func (l *Loader) Load(jsonType string, dec *json.Decoder) error {
	switch jsonType {
	case "loans":
		return l.loadLoans(dec)
	default:
		return fmt.Errorf("unknown type \"%v\"", jsonType)
	}
}

/*
func LoadOLD(jsonType string, json map[string]interface{}, tx *sql.Tx,
	opts *LoadOptions) error {
	id := json["id"].(string)
	switch jsonType {
	case "groups":
		return loadGroups(id, json, tx, opts)
	case "users":
		return loadUsers(id, json, tx, opts)
	//case "loans":
	//return loadLoans(id, json, tx, opts)
	case "tmp_loans_locations":
		return loadTmpLoansLocs(id, json, tx, opts)
	default:
		return fmt.Errorf("unknown type \"%v\"", jsonType)
	}
}
*/

func NewLoader(db *sql.DB, opts *LoaderOptions) (*Loader, error) {
	// Start transaction for exclusive lock
	locktx, err := db.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	// Request lock and block until obtained
	_, err = locktx.Exec("" +
		"LOCK TABLE loading.exlock;")
	if err != nil {
		return nil, err
	}
	// Start transaction for main loader
	tx, err := db.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	// Use defaults if opts == nil
	var lopts *LoaderOptions
	if opts == nil {
		lopts = &LoaderOptions{}
	} else {
		lopts = opts
	}
	// Instantiate new loader
	nl := &Loader{
		db:     db,
		tx:     tx,
		locktx: locktx,
		opts:   lopts,
	}
	return nl, nil
}

func (l *Loader) Close() error {
	// Commit all changes
	err := l.tx.Commit()
	if err != nil {
		return err
	}
	// Release exclusive lock
	err = l.locktx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// OLD
/*
type LoadOptions struct {
	// Debug enables debugging output if set to true.
	Debug bool
}
*/
