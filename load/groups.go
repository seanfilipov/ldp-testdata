package load

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Update(jsonType string, json map[string]interface{}, tx *sql.Tx) error {
	switch jsonType {
	case "groups":
		return UpdateGroups(json, tx)
	default:
		return fmt.Errorf("unknown type \"%v\"", jsonType)
	}
}

func UpdateGroups(json map[string]interface{}, tx *sql.Tx) error {
	id := json["id"].(string)
	groupName := json["group"].(string)
	description := json["desc"].(string)
	sql1 := "INSERT INTO groups AS g\n" +
		"    (id, group_name, description)\n" +
		"    VALUES ($1, $2, $3)\n" +
		"    ON CONFLICT (id) DO\n" +
		"        UPDATE SET group_name = EXCLUDED.group_name,\n" +
		"                   description = EXCLUDED.description\n" +
		"            WHERE g.group_name <> EXCLUDED.group_name OR\n" +
		"                  g.description <> EXCLUDED.description;\n"
	_, err := Exec(tx, sql1, id, groupName, description)
	return err
}

func Exec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	fmt.Fprintf(os.Stderr, "%s{", query)
	for x, a := range args {
		if x != 0 {
			fmt.Fprintf(os.Stderr, ", ")
		}
		fmt.Fprintf(os.Stderr, "\"%s\"", a)
	}
	fmt.Fprintf(os.Stderr, "}\n\n")
	return tx.Exec(query, args...)
}
