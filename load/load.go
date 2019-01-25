package load

import (
	"database/sql"
	"fmt"
)

func Update(jsonType string, json map[string]interface{}, tx *sql.Tx) error {
	id := json["id"].(string)
	switch jsonType {
	case "groups":
		return updateGroups(id, json, tx)
	case "users":
		return updateUsers(id, json, tx)
	default:
		return fmt.Errorf("unknown type \"%v\"", jsonType)
	}
}
