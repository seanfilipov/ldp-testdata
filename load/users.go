package load

import (
	"database/sql"
)

func updateUsers(id string, json map[string]interface{}, tx *sql.Tx) error {
	if json != nil {
		username := json["username"].(string)
		barcode := json["barcode"].(string)
		userType := json["type"].(string)
		active := json["active"].(string)
		patronGroupId := json["patronGroup"].(string)
		err := updateGroups(patronGroupId, nil, tx)
		if err != nil {
			return err
		}
		_, err = exec(tx, sqlUpdateUsers, id, username, barcode,
			userType, active, patronGroupId)
		return err
	} else {
		_, err := exec(tx, sqlUpdateUsersEmpty, id)
		return err
	}
}

var sqlUpdateUsers string = compact("" +
	"  INSERT INTO users AS u                                    \n" +
	"      (id, username, barcode, user_type, active,            \n" +
	"              patron_group_id)                              \n" +
	"      VALUES ($1, $2, $3, $4, $5, $6)                       \n" +
	"      ON CONFLICT (id) DO UPDATE                            \n" +
	"      SET username = EXCLUDED.username,                     \n" +
	"          barcode = EXCLUDED.barcode,                       \n" +
	"          user_type = EXCLUDED.user_type,                   \n" +
	"          active = EXCLUDED.active,                         \n" +
	"          patron_group_id = EXCLUDED.patron_group_id        \n" +
	"      WHERE u.username <> EXCLUDED.username OR              \n" +
	"            u.barcode <> EXCLUDED.barcode OR                \n" +
	"            u.user_type <> EXCLUDED.user_type OR            \n" +
	"            u.active <> EXCLUDED.active OR                  \n" +
	"            u.patron_group_id <> EXCLUDED.patron_group_id;  \n")

var sqlUpdateUsersEmpty string = compact("" +
	"  INSERT INTO users                 \n" +
	"      (id)                          \n" +
	"      VALUES ($1)                   \n" +
	"      ON CONFLICT (id) DO NOTHING;  \n")
