package ldpadmin

import "encoding/json"

func (l *Loader) loadGroups(dec *json.Decoder) error {
	err := l.sqlTruncateStage("groups")
	if err != nil {
		return err
	}
	stmt, err := l.sqlCopyStage("groups",
		"group_id", "group_name", "description")
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
		groupId := j["id"].(string)
		groupName := j["group"].(string)
		description := j["desc"].(string)
		_, err = l.sqlCopyExec(stmt, groupId, groupName, description)
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
	// Upsert group records
	_, err = l.sqlExec("" +
		"INSERT INTO normal.groups AS g\n" +
		"    (group_id, group_name, description)\n" +
		"    SELECT lg.group_id,\n" +
		"           lg.group_name,\n" +
		"           lg.description\n" +
		"        FROM loading.groups AS lg\n" +
		"    ON CONFLICT (group_id) DO UPDATE\n" +
		"    SET group_name = EXCLUDED.group_name,\n" +
		"        description = EXCLUDED.description\n" +
		"    WHERE g.group_name <> EXCLUDED.group_name OR\n" +
		"          g.description <> EXCLUDED.description;\n")
	if err != nil {
		return err
	}
	err = l.sqlTruncateStage("groups")
	if err != nil {
		return err
	}
	return nil
}
