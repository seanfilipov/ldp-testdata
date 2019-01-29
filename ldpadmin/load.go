package ldpadmin

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func Load(jsonType string, dec *json.Decoder, tx *sql.Tx,
	opts *LoadOptions) error {
	switch jsonType {
	case "loans":
		return loadLoans(dec, tx, opts)
	default:
		return fmt.Errorf("unknown type \"%v\"", jsonType)
	}
}

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

type LoadOptions struct {
	// Debug enables debugging output if set to true.
	Debug bool
}
