package ldpadmin

/*
func loadUsers(id string, json map[string]interface{}, tx *sql.Tx,
	opts *LoadOptions) error {
	if json != nil {
		username := json["username"].(string)
		barcode := json["barcode"].(string)
		userType := json["type"].(string)
		active := json["active"].(string)
		patronGroupId := json["patronGroup"].(string)
		err := loadGroups(patronGroupId, nil, tx, opts)
		if err != nil {
			return err
		}
		_, err = exec(tx, opts, sqlLoadUsers, id, username, barcode,
			userType, active, patronGroupId)
		// d_users
		_, err = exec(tx, opts, sqlLoadDUsers, id, username, barcode,
			userType, active, patronGroupId)
		return err
	} else {
		_, err := exec(tx, opts, sqlLoadUsersEmpty, id)
		return err
	}
}

var sqlLoadUsers string = trimSql("" +
	"  INSERT INTO users AS u                                    \n" +
	"      (id, username, barcode, user_type, active,            \n" +
	"              patron_group_id)                              \n" +
	"      VALUES ($1,                                           \n" +
	"              $2,                                           \n" +
	"              $3,                                           \n" +
	"              $4,                                           \n" +
	"              $5,                                           \n" +
	"              $6)                                           \n" +
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

var sqlLoadDUsers string = trimSql("" +
	"  INSERT INTO d_users AS u                                      \n" +
	"      (id, username, barcode, user_type, active,                \n" +
	"              group_name, group_description)                    \n" +
	"      SELECT $1,                                                \n" +
	"             $2,                                                \n" +
	"             $3,                                                \n" +
	"             $4,                                                \n" +
	"             $5,                                                \n" +
	"             g.group_name,                                      \n" +
	"             g.description                                      \n" +
	"          FROM groups g                                         \n" +
	"          WHERE g.id = $6                                       \n" +
	"      ON CONFLICT (id) DO UPDATE                                \n" +
	"      SET username = EXCLUDED.username,                         \n" +
	"          barcode = EXCLUDED.barcode,                           \n" +
	"          user_type = EXCLUDED.user_type,                       \n" +
	"          active = EXCLUDED.active,                             \n" +
	"          group_name = EXCLUDED.group_name,                     \n" +
	"          group_description = EXCLUDED.group_description        \n" +
	"      WHERE u.username <> EXCLUDED.username OR                  \n" +
	"            u.barcode <> EXCLUDED.barcode OR                    \n" +
	"            u.user_type <> EXCLUDED.user_type OR                \n" +
	"            u.active <> EXCLUDED.active OR                      \n" +
	"            u.group_name <> EXCLUDED.group_name OR              \n" +
	"            u.group_description <> EXCLUDED.group_description;  \n")

var sqlLoadUsersEmpty string = trimSql("" +
	"  INSERT INTO users                 \n" +
	"      (id)                          \n" +
	"      VALUES ($1)                   \n" +
	"      ON CONFLICT (id) DO NOTHING;  \n")
*/
