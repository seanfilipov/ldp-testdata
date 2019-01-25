package loader

import (
	"database/sql"
)

func updateGroups(id string, json map[string]interface{}, tx *sql.Tx) error {
	if json != nil {
		groupName := json["group"].(string)
		description := json["desc"].(string)
		_, err := exec(tx, sqlUpdateGroups, id, groupName, description)
		return err
	} else {
		_, err := exec(tx, sqlUpdateGroupsEmpty, id)
		return err
	}
}

var sqlUpdateGroups string = trimSql("" +
	"  INSERT INTO groups AS g                           \n" +
	"      (id, group_name, description)                 \n" +
	"      VALUES ($1, $2, $3)                           \n" +
	"      ON CONFLICT (id) DO UPDATE                    \n" +
	"      SET group_name = EXCLUDED.group_name,         \n" +
	"          description = EXCLUDED.description        \n" +
	"      WHERE g.group_name <> EXCLUDED.group_name OR  \n" +
	"            g.description <> EXCLUDED.description;  \n")

var sqlUpdateGroupsEmpty string = trimSql("" +
	"  INSERT INTO groups                \n" +
	"      (id)                          \n" +
	"      VALUES ($1)                   \n" +
	"      ON CONFLICT (id) DO NOTHING;  \n")
