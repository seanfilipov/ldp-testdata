package ldpadmin

/*
func loadGroups(id string, json map[string]interface{}, tx *sql.Tx,
	opts *LoadOptions) error {
	if json != nil {
		groupName := json["group"].(string)
		description := json["desc"].(string)
		_, err := exec(tx, opts, sqlLoadGroups, id, groupName,
			description)
		return err
	} else {
		_, err := exec(tx, opts, sqlLoadGroupsEmpty, id)
		return err
	}
}

var sqlLoadGroups string = trimSql("" +
	"  INSERT INTO groups AS g                           \n" +
	"      (id, group_name, description)                 \n" +
	"      VALUES ($1,                                   \n" +
	"              $2,                                   \n" +
	"              $3)                                   \n" +
	"      ON CONFLICT (id) DO UPDATE                    \n" +
	"      SET group_name = EXCLUDED.group_name,         \n" +
	"          description = EXCLUDED.description        \n" +
	"      WHERE g.group_name <> EXCLUDED.group_name OR  \n" +
	"            g.description <> EXCLUDED.description;  \n")

var sqlLoadGroupsEmpty string = trimSql("" +
	"  INSERT INTO groups                \n" +
	"      (id)                          \n" +
	"      VALUES ($1)                   \n" +
	"      ON CONFLICT (id) DO NOTHING;  \n")
*/
