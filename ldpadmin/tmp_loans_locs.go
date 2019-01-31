package ldpadmin

/*
func loadTmpLoansLocs(id string, json map[string]interface{},
	tx *sql.Tx, opts *LoadOptions) error {
	if json != nil {
		loanId := json["id"].(string)
		item := json["item"].(map[string]interface{})
		location := item["location"]
		locationName :=
			location.(map[string]interface{})["name"].(string)
		_, err := exec(tx, opts, sqlLoadTmpLoansLocs, loanId,
			locationName)
		// d_locations
		_, err = exec(tx, opts, sqlLoadDLocations, locationName,
			locationName)
		return err
	} else {
		_, err := exec(tx, opts, sqlLoadTmpLoansLocsEmpty, id)
		return err
	}
}

var sqlLoadTmpLoansLocs string = trimSql("" +
	"  INSERT INTO tmp_loans_locations AS t                  \n" +
	"      (loan_id, location_name)                          \n" +
	"      VALUES ($1,                                       \n" +
	"              $2)                                       \n" +
	"      ON CONFLICT (loan_id) DO UPDATE                   \n" +
	"      SET location_name = EXCLUDED.location_name        \n" +
	"      WHERE t.location_name <> EXCLUDED.location_name;  \n")

var sqlLoadDLocations string = trimSql("" +
	"  INSERT INTO d_locations AS l                          \n" +
	"      (id, location_name)                               \n" +
	"      SELECT 'id-' || replace(lower($1), ' ', '-'),     \n" +
	"             $2                                         \n" +
	"      ON CONFLICT (id) DO UPDATE                        \n" +
	"      SET location_name = EXCLUDED.location_name        \n" +
	"      WHERE l.location_name <> EXCLUDED.location_name;  \n")

var sqlLoadTmpLoansLocsEmpty string = trimSql("" +
	"  INSERT INTO tmp_loans_locations        \n" +
	"      (loan_id)                          \n" +
	"      VALUES ($1)                        \n" +
	"      ON CONFLICT (loan_id) DO NOTHING;  \n")
*/
