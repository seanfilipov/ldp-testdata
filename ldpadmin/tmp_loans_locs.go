package ldpadmin

import "encoding/json"

func (l *Loader) loadTmpLoansLocs(dec *json.Decoder) error {
	err := l.sqlTruncateStage("tmp_loans_locations")
	if err != nil {
		return err
	}
	stmt, err := l.sqlCopyStage("tmp_loans_locations",
		"loan_id", "location_name")
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
		loanId := j["id"].(string)
		item := j["item"].(map[string]interface{})
		location := item["location"]
		locationName :=
			location.(map[string]interface{})["name"].(string)
		_, err = l.sqlCopyExec(stmt, loanId, locationName)
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
	// Upsert tmp_loans_locations
	_, err = l.sqlExec("" +
		"INSERT INTO normal.tmp_loans_locations AS t\n" +
		"    (loan_id, location_name)\n" +
		"    SELECT lt.loan_id,\n" +
		"           lt.location_name\n" +
		"        FROM loading.tmp_loans_locations AS lt\n" +
		"    ON CONFLICT (loan_id) DO UPDATE\n" +
		"    SET location_name = EXCLUDED.location_name\n" +
		"    WHERE t.location_name <> EXCLUDED.location_name;\n")
	if err != nil {
		return err
	}
	// Upsert locations
	_, err = l.sqlExec("" +
		"INSERT INTO locations AS l\n" +
		"    (location_key, location_name)\n" +
		"    SELECT 'id-' ||\n" +
		"               replace(lower(lt.location_name), ' ', '-'),\n" +
		"           lt.location_name\n" +
		"        FROM (\n" +
		"            SELECT DISTINCT location_name\n" +
		"                FROM loading.tmp_loans_locations\n" +
		"             ) AS lt\n" +
		"    ON CONFLICT (location_key) DO UPDATE\n" +
		"    SET location_name = EXCLUDED.location_name\n" +
		"    WHERE l.location_name <> EXCLUDED.location_name;\n")
	if err != nil {
		return err
	}
	err = l.sqlTruncateStage("tmp_loans_locations")
	if err != nil {
		return err
	}
	return nil
}
