package ldpadmin

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

func loadLoans(dec *json.Decoder, tx *sql.Tx, opts *LoadOptions) error {
	stmt, err := tx.Prepare(pq.CopyInSchema("load", "f_loans",
		"id", "user_id", "location_id", "item_id", "action",
		"status_name", "loan_date", "due_date"))
	if err != nil {
		return err
	}
	for dec.More() {
		var i interface{}
		err := dec.Decode(&i)
		if err != nil {
			return err
		}
		json := i.(map[string]interface{})
		id := json["id"].(string)
		userId := json["userId"].(string)
		itemId := json["itemId"].(string)
		action := json["action"].(string)
		status := json["status"].(map[string]interface{})
		statusName := status["name"].(string)
		loanDateStr := json["loanDate"].(string)
		dueDateStr := json["dueDate"].(string)
		layout := "2006-01-02T15:04:05Z"
		loanDate, _ := time.Parse(layout, loanDateStr)
		dueDate, _ := time.Parse(layout, dueDateStr)
		_, err = stmt.Exec(id, userId, "", itemId, action, statusName,
			loanDate, dueDate)
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"INSERT INTO d_users\n" +
		"    (id)\n" +
		"    SELECT user_id AS id\n" +
		"        FROM load.f_loans\n" +
		"    ON CONFLICT (id) DO NOTHING;\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"INSERT INTO f_loans AS fl\n" +
		"    (id, user_id, location_id, item_id, action,\n" +
		"            status_name, loan_date, due_date)\n" +
		"    SELECT lfl.id,\n" +
		"           lfl.user_id,\n" +
		"           'id-' || replace(lower(tll.location_name), ' ',\n" +
		"                   '-') AS location_id,\n" +
		"           lfl.item_id,\n" +
		"           lfl.action,\n" +
		"           lfl.status_name,\n" +
		"           lfl.loan_date,\n" +
		"           lfl.due_date\n" +
		"        FROM load.f_loans AS lfl\n" +
		"            LEFT JOIN tmp_loans_locations AS tll\n" +
		"                ON lfl.id = tll.loan_id\n" +
		"    ON CONFLICT (id) DO UPDATE\n" +
		"    SET user_id = EXCLUDED.user_id,\n" +
		"        location_id = EXCLUDED.location_id,\n" +
		"        item_id = EXCLUDED.item_id,\n" +
		"        action = EXCLUDED.action,\n" +
		"        status_name = EXCLUDED.status_name,\n" +
		"        loan_date = EXCLUDED.loan_date,\n" +
		"        due_date = EXCLUDED.due_date\n" +
		"    WHERE fl.user_id <> EXCLUDED.user_id OR\n" +
		"          fl.location_id <> EXCLUDED.location_id OR\n" +
		"          fl.item_id <> EXCLUDED.item_id OR\n" +
		"          fl.action <> EXCLUDED.action OR\n" +
		"          fl.status_name <> EXCLUDED.status_name OR\n" +
		"          fl.loan_date <> EXCLUDED.loan_date OR\n" +
		"          fl.due_date <> EXCLUDED.due_date;\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"TRUNCATE load.f_loans;")
	if err != nil {
		return err
	}
	return nil
}
