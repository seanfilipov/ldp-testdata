package ldpadmin

import (
	"encoding/json"
	"time"
)

func (l *Loader) loadLoans(dec *json.Decoder) error {
	err := l.sqlTruncateStage("loans")
	if err != nil {
		return err
	}
	stmt, err := l.sqlCopyStage("loans",
		"loan_id", "user_id", "location_id", "item_id", "action",
		"status_name", "loan_date", "due_date")
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
		loanId := json["id"].(string)
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
		_, err = l.sqlCopyExec(stmt, loanId, userId, "", itemId, action,
			statusName, loanDate, dueDate)
	}
	_, err = l.sqlCopyExec(stmt)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	err = l.sqlMergePlaceholders("users", "user_id", "loans", "user_id")
	if err != nil {
		return err
	}
	_, err = l.sqlExec("" +
		"INSERT INTO loans AS l\n" +
		"    (loan_id, user_key, location_id, item_id, action,\n" +
		"            status_name, loan_date, due_date)\n" +
		"    SELECT sl.loan_id,\n" +
		"           ( SELECT u.user_key\n" +
		"                 FROM users AS u\n" +
		"                 WHERE sl.user_id = u.user_id\n" +
		"                 ORDER BY record_time DESC LIMIT 1\n" +
		"           ),\n" +
		"           'id-' || replace(lower(tll.location_name), ' ',\n" +
		"                   '-') AS location_id,\n" +
		"           sl.item_id,\n" +
		"           sl.action,\n" +
		"           sl.status_name,\n" +
		"           sl.loan_date,\n" +
		"           sl.due_date\n" +
		"        FROM internal.loans AS sl\n" +
		"            LEFT JOIN norm.tmp_loans_locations AS tll\n" +
		"                ON sl.loan_id = tll.loan_id\n" +
		"    ON CONFLICT (loan_id) DO UPDATE\n" +
		"    SET user_key = EXCLUDED.user_key,\n" +
		"        location_id = EXCLUDED.location_id,\n" +
		"        item_id = EXCLUDED.item_id,\n" +
		"        action = EXCLUDED.action,\n" +
		"        status_name = EXCLUDED.status_name,\n" +
		"        loan_date = EXCLUDED.loan_date,\n" +
		"        due_date = EXCLUDED.due_date\n" +
		"    WHERE l.user_key <> EXCLUDED.user_key OR\n" +
		"          l.location_id <> EXCLUDED.location_id OR\n" +
		"          l.item_id <> EXCLUDED.item_id OR\n" +
		"          l.action <> EXCLUDED.action OR\n" +
		"          l.status_name <> EXCLUDED.status_name OR\n" +
		"          l.loan_date <> EXCLUDED.loan_date OR\n" +
		"          l.due_date <> EXCLUDED.due_date;\n")
	if err != nil {
		return err
	}
	err = l.sqlTruncateStage("loans")
	if err != nil {
		return err
	}
	return nil
}
