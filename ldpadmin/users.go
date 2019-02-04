package ldpadmin

import "encoding/json"

func (l *Loader) loadUsers(dec *json.Decoder) error {
	err := l.sqlTruncateStage("users")
	if err != nil {
		return err
	}
	stmt, err := l.sqlCopyStage("users",
		"user_id", "username", "barcode", "user_type", "active",
		"patron_group_id")
	if err != nil {
		return err
	}
	for dec.More() {
		var i interface{}
		err := dec.Decode(&i)
		if err != nil {
			return err
		}
		j := i.(map[string]interface{})
		userId := j["id"].(string)
		username := j["username"].(string)
		barcode := j["barcode"].(string)
		userType := j["type"].(string)
		active := j["active"].(string)
		patronGroupId := j["patronGroup"].(string)
		_, err = l.sqlCopyExec(stmt, userId, username, barcode,
			userType, active, patronGroupId)
		if err != nil {
			return err
		}
	}
	_, err = l.sqlCopyExec(stmt)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	// Merge placeholders for groups
	err = l.sqlMergePlaceholders("normal.groups", "group_id", "users",
		"patron_group_id")
	if err != nil {
		return err
	}
	// Insert user records except for those with placeholders
	_, err = l.sqlExec("" +
		"INSERT INTO users\n" +
		"    (user_id, username, barcode, user_type, active,\n" +
		"            group_name, group_description)\n" +
		"    SELECT lu.user_id,\n" +
		"           lu.username,\n" +
		"           lu.barcode,\n" +
		"           lu.user_type,\n" +
		"           lu.active,\n" +
		"           g.group_name,\n" +
		"           g.description\n" +
		"        FROM loading.users AS lu\n" +
		"            LEFT JOIN normal.groups AS g\n" +
		"                ON lu.patron_group_id = g.group_id\n" +
		"        WHERE NOT EXISTS\n" +
		"          ( SELECT 1\n" +
		"                FROM users AS u\n" +
		"                WHERE lu.user_id = u.user_id AND\n" +
		"                      u.username = 'NOT AVAILABLE'\n" +
		"          );\n")
	if err != nil {
		return err
	}
	// Update placeholder records
	_, err = l.sqlExec("" +
		"UPDATE users AS u\n" +
		"    SET username = lu.username,\n" +
		"        barcode = lu.barcode,\n" +
		"        user_type = lu.user_type,\n" +
		"        active = lu.active,\n" +
		"        group_name = g.group_name,\n" +
		"        group_description = g.description\n" +
		"    FROM loading.users AS lu\n" +
		"        LEFT JOIN normal.groups AS g\n" +
		"            ON lu.patron_group_id = g.group_id\n" +
		"    WHERE u.user_id = lu.user_id AND\n" +
		"          u.username = 'NOT AVAILABLE'\n;")
	if err != nil {
		return err
	}
	err = l.sqlTruncateStage("users")
	if err != nil {
		return err
	}
	return nil
}
