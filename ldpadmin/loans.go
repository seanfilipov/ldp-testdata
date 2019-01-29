package ldpadmin

import (
	"database/sql"
	"time"
)

func loadLoans(id string, json map[string]interface{}, tx *sql.Tx,
	opts *LoadOptions) error {
	if json != nil {
		userId := json["userId"].(string)
		err := loadUsers(userId, nil, tx, opts)
		if err != nil {
			return err
		}
		itemId := json["itemId"].(string)
		action := json["action"].(string)
		status := json["status"].(map[string]interface{})
		statusName := status["name"].(string)
		loanDateStr := json["loanDate"].(string)
		dueDateStr := json["dueDate"].(string)
		layout := "2006-01-02T15:04:05Z"
		loanDate, _ := time.Parse(layout, loanDateStr)
		dueDate, _ := time.Parse(layout, dueDateStr)
		_, err = exec(tx, opts, sqlLoadLoans, id, userId, itemId,
			action, statusName, loanDate, dueDate)
		// f_loans
		_, err = exec(tx, opts, sqlLoadFLoans, id, userId, itemId, action,
			statusName, loanDate, dueDate, id)
		return err
	} else {
		_, err := exec(tx, opts, sqlLoadLoansEmpty, id)
		return err
	}
}

var sqlLoadLoans string = trimSql("" +
	"  INSERT INTO loans AS l                                      \n" +
	"      (id, user_id, item_id, action, status_name, loan_date,  \n" +
	"              due_date)                                       \n" +
	"      VALUES ($1,                                             \n" +
	"              $2,                                             \n" +
	"              $3,                                             \n" +
	"              $4,                                             \n" +
	"              $5,                                             \n" +
	"              $6,                                             \n" +
	"              $7)                                             \n" +
	"      ON CONFLICT (id) DO UPDATE                              \n" +
	"      SET user_id = EXCLUDED.user_id,                         \n" +
	"          item_id = EXCLUDED.item_id,                         \n" +
	"          action = EXCLUDED.action,                           \n" +
	"          status_name = EXCLUDED.status_name,                 \n" +
	"          loan_date = EXCLUDED.loan_date,                     \n" +
	"          due_date = EXCLUDED.due_date                        \n" +
	"      WHERE l.user_id <> EXCLUDED.user_id OR                  \n" +
	"            l.item_id <> EXCLUDED.item_id OR                  \n" +
	"            l.action <> EXCLUDED.action OR                    \n" +
	"            l.status_name <> EXCLUDED.status_name OR          \n" +
	"            l.loan_date <> EXCLUDED.loan_date OR              \n" +
	"            l.due_date <> EXCLUDED.due_date;                  \n")

var sqlLoadFLoans string = trimSql("" +
	"  INSERT INTO f_loans AS l                                       \n" +
	"      (id, user_id, location_id, item_id, action, status_name,   \n" +
	"              loan_date, due_date)                               \n" +
	"      SELECT $1,                                                 \n" +
	"             $2,                                                 \n" +
	"          'id-' || replace(lower(tll.location_name), ' ', '-'),  \n" +
	"              $3,                                                \n" +
	"              $4,                                                \n" +
	"              $5,                                                \n" +
	"              $6,                                                \n" +
	"              $7                                                 \n" +
	"          FROM tmp_loans_locations tll                           \n" +
	"          WHERE tll.loan_id = $8                                 \n" +
	"      ON CONFLICT (id) DO UPDATE                                 \n" +
	"      SET user_id = EXCLUDED.user_id,                            \n" +
	"          location_id = EXCLUDED.location_id,                    \n" +
	"          item_id = EXCLUDED.item_id,                            \n" +
	"          action = EXCLUDED.action,                              \n" +
	"          status_name = EXCLUDED.status_name,                    \n" +
	"          loan_date = EXCLUDED.loan_date,                        \n" +
	"          due_date = EXCLUDED.due_date                           \n" +
	"      WHERE l.user_id <> EXCLUDED.user_id OR                     \n" +
	"            l.location_id <> EXCLUDED.location_id OR             \n" +
	"            l.item_id <> EXCLUDED.item_id OR                     \n" +
	"            l.action <> EXCLUDED.action OR                       \n" +
	"            l.status_name <> EXCLUDED.status_name OR             \n" +
	"            l.loan_date <> EXCLUDED.loan_date OR                 \n" +
	"            l.due_date <> EXCLUDED.due_date;                     \n")

var sqlLoadLoansEmpty string = trimSql("" +
	"  INSERT INTO loans                 \n" +
	"      (id)                          \n" +
	"      VALUES ($1)                   \n" +
	"      ON CONFLICT (id) DO NOTHING;  \n")
