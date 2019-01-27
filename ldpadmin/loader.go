package ldpadmin

import (
	"database/sql"
	"fmt"
)

func Update(jsonType string, json map[string]interface{}, tx *sql.Tx,
	opts *UpdateOptions) error {
	id := json["id"].(string)
	switch jsonType {
	case "groups":
		return updateGroups(id, json, tx, opts)
	case "users":
		return updateUsers(id, json, tx, opts)
	case "loans":
		return updateLoans(id, json, tx, opts)
	case "tmp_loans_locations":
		return updateTmpLoansLocs(id, json, tx, opts)
	default:
		return fmt.Errorf("unknown type \"%v\"", jsonType)
	}
}

type UpdateOptions struct {
	// Debug enables debugging output if set to true.
	Debug bool
}
