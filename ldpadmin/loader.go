package ldpadmin

import (
	"database/sql"
	"fmt"
)

func Load(jsonType string, json map[string]interface{}, tx *sql.Tx,
	opts *LoadOptions) error {
	id := json["id"].(string)
	switch jsonType {
	case "groups":
		return loadGroups(id, json, tx, opts)
	case "users":
		return loadUsers(id, json, tx, opts)
	case "loans":
		return loadLoans(id, json, tx, opts)
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
