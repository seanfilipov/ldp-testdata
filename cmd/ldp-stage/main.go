package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/folio-org/ldp/cmd/internal/ldputil"
)

func main() {
	config, err := ldputil.ReadConfig()
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	sourcedir := config.Get("extract", "dir")
	targetdir := config.Get("stage", "dir")
	// TODO create new temp directory under targetdir

	var st stage

	st.start, err = os.Create(targetdir + "/00001-start.sql")
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer st.start.Close()
	fmt.Fprintf(st.start, "%s\n", sqlStartTransaction())

	st.deleteStage, err = os.Create(targetdir + "/00002-delete-stage.sql")
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer st.deleteStage.Close()
	fmt.Fprintf(st.deleteStage, `
DELETE FROM stage.users;
DELETE FROM stage.groups;
`)

	st.groups, err = os.Create(targetdir + "/00010-groups.sql")
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer st.groups.Close()
	fmt.Fprintf(st.groups,
		"COPY stage.groups (id, group_name, description) FROM stdin;\n")

	st.usersNaGroups, err = os.Create(targetdir + "/00020-users-na-groups.sql")
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer st.usersNaGroups.Close()
	fmt.Fprintf(st.usersNaGroups, `
INSERT INTO stage.groups (id) VALUES
`)

	st.users, err = os.Create(targetdir + "/00021-users.sql")
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer st.users.Close()
	fmt.Fprintf(st.users,
		"COPY stage.users (id, username, barcode, user_type, active, "+
			"patron_group_id) FROM stdin;\n")

	st.mergeAll, err = os.Create(targetdir + "/90000-merge-all.sql")
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer st.mergeAll.Close()
	fmt.Fprintf(st.mergeAll, `
INSERT INTO public.groups (id, group_name, description)
    SELECT id, group_name, description
        FROM stage.groups
    ON CONFLICT (id) DO
        UPDATE SET group_name = EXCLUDED.group_name,
	           description = EXCLUDED.description;
`)
	fmt.Fprintf(st.mergeAll, `
INSERT INTO public.users (id, username, barcode, user_type, active, patron_group_id)
    SELECT id, username, barcode, user_type, active, patron_group_id
        FROM stage.users
    ON CONFLICT (id) DO
        UPDATE SET username = EXCLUDED.username,
	           barcode = EXCLUDED.barcode,
	           user_type = EXCLUDED.user_type,
	           active = EXCLUDED.active,
	           patron_group_id = EXCLUDED.patron_group_id;
`)

	st.end, err = os.Create(targetdir + "/99999-end.sql")
	if err != nil {
		ldputil.PrintError(err)
		return
	}
	defer st.end.Close()
	fmt.Fprintf(st.end, "%s\n", sqlCommit())

	err = stageAll("groups", sourcedir+"/groups.json", &st)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	err = stageAll("users", sourcedir+"/users.json", &st)
	if err != nil {
		ldputil.PrintError(err)
		return
	}

	/*
		for x := 1; x <= 20; x++ {
			err = stageAll("tmp_loans_locations",
				sourcedir+fmt.Sprintf("/circulation.loans.json.%v",
					x),
				tx)
			if err != nil {
				ldputil.PrintError(err)
				return
			}
		}

		for x := 1; x <= 20; x++ {
			err = stageAll("loans",
				sourcedir+fmt.Sprintf("/loan-storage.loans.json.%v",
					x),
				tx)
			if err != nil {
				ldputil.PrintError(err)
				return
			}
		}
	*/

	fmt.Fprintf(st.groups, "\\.\n")
	fmt.Fprintf(st.usersNaGroups, `
ON CONFLICT (id) DO NOTHING;
`)
	fmt.Fprintf(st.users, "\\.\n")
}

func stageAll(jtype string, filename string, st *stage) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	// Skip past first tokens.
	for x := 0; x < 3; x++ {
		_, err := dec.Token()
		if err != nil {
			return err
		}
	}

	// Read and stage array elements.
	count := 0
	for dec.More() {

		var i interface{}
		err := dec.Decode(&i)
		if err != nil {
			return err
		}

		err = stageOne(st, jtype, i.(map[string]interface{}))
		if err != nil {
			return err
		}

		count++
	}

	fmt.Printf("%s %v %s\n", jtype, count, filename)

	return nil
}

func stageOne(st *stage, jtype string, j map[string]interface{}) error {
	switch jtype {
	case "groups":
		err := stageGroup(st, jtype, j)
		if err != nil {
			return err
		}
	case "users":
		err := stageUser(st, jtype, j)
		if err != nil {
			return err
		}
		/*
			case "tmp_loans_locations":
				err := stageTmpLoanLocation(tx, jtype, j)
				if err != nil {
					return err
				}
			case "loans":
				err := stageLoan(tx, jtype, j)
				if err != nil {
					return err
				}
		*/
	default:
		return fmt.Errorf("unknown type \"%v\"", jtype)
	}

	return nil
}

func stageGroup(st *stage, jtype string, j map[string]interface{}) error {

	id := j["id"].(string)
	groupName := j["group"].(string)
	description := j["desc"].(string)

	fmt.Fprintf(st.groups, "%v\t%v\t%v\n", id, groupName, description)

	//_, err := tx.Exec(
	//"INSERT INTO groups AS g "+
	//"(id, group_name, description) "+
	//"VALUES ($1, $2, $3) "+
	//"ON CONFLICT (id) DO "+
	//"UPDATE SET group_name = EXCLUDED.group_name, "+
	//"description = EXCLUDED.description "+
	//"WHERE g.group_name <> EXCLUDED.group_name OR "+
	//"g.description <> EXCLUDED.description",
	//id, groupName, description)
	//if err != nil {
	//return err
	//}

	return nil
}

func stageUser(st *stage, jtype string, j map[string]interface{}) error {

	id := j["id"].(string)
	username := j["username"].(string)
	barcode := j["barcode"].(string)
	userType := j["type"].(string)
	active := j["active"].(string)
	patronGroupId := j["patronGroup"].(string)

	if st.usersNaGroupsCount != 0 {
		fmt.Fprintf(st.usersNaGroups, ",\n")
	}
	fmt.Fprintf(st.usersNaGroups, "('%v')", patronGroupId)
	st.usersNaGroupsCount++

	fmt.Fprintf(st.users, "%v\t%v\t%v\t%v\t%v\t%v\n", id, username,
		barcode, userType, active, patronGroupId)

	//_, err := tx.Exec(
	//"INSERT INTO users AS u "+
	//"(id, username, barcode, user_type, active, "+
	//"patron_group_id) "+
	//"VALUES ($1, $2, $3, $4, $5, $6) "+
	//"ON CONFLICT (id) DO "+
	//"UPDATE SET username = EXCLUDED.username, "+
	//"barcode = EXCLUDED.barcode, "+
	//"user_type = EXCLUDED.user_type, "+
	//"active = EXCLUDED.active, "+
	//"patron_group_id = EXCLUDED.patron_group_id "+
	//"WHERE u.username <> EXCLUDED.username OR "+
	//"u.barcode <> EXCLUDED.barcode OR "+
	//"u.user_type <> EXCLUDED.user_type OR "+
	//"u.active <> EXCLUDED.active OR "+
	//"u.patron_group_id <> EXCLUDED.patron_group_id",
	//id, username, barcode, userType, active, patronGroupId)
	//if err != nil {
	//return err
	//}

	return nil
}

/*
func stageTmpLoanLocation(tx *sql.Tx, jtype string,
	j map[string]interface{}) error {

	loanId := j["id"].(string)
	item := j["item"].(map[string]interface{})
	location := item["location"]
	locationName := location.(map[string]interface{})["name"].(string)

	//_, err := tx.Exec(
	//"INSERT INTO tmp_loans_locations AS t "+
	//"(loan_id, location_name) "+
	//"VALUES ($1, $2) "+
	//"ON CONFLICT (loan_id) DO "+
	//"UPDATE SET location_name = EXCLUDED.location_name "+
	//"WHERE t.location_name <> EXCLUDED.location_name",
	//loanId, locationName)
	//if err != nil {
	//return err
	//}

	return nil
}

func stageLoan(tx *sql.Tx, jtype string, j map[string]interface{}) error {

	id := j["id"].(string)
	userId := j["userId"].(string)
	itemId := j["itemId"].(string)
	action := j["action"].(string)

	status := j["status"].(map[string]interface{})
	statusName := status["name"].(string)

	loanDateStr := j["loanDate"].(string)
	dueDateStr := j["dueDate"].(string)

	layout := "2006-01-02T15:04:05Z"
	loanDate, _ := time.Parse(layout, loanDateStr)
	dueDate, _ := time.Parse(layout, dueDateStr)

	//_, err := tx.Exec(
	//"INSERT INTO loans AS l "+
	//"(id, user_id, item_id, action, status_name, "+
	//"loan_date, due_date) "+
	//"VALUES ($1, $2, $3, $4, $5, $6, $7) "+
	//"ON CONFLICT (id) DO "+
	//"UPDATE SET user_id = EXCLUDED.user_id, "+
	//"item_id = EXCLUDED.item_id, "+
	//"action = EXCLUDED.action, "+
	//"status_name = EXCLUDED.status_name, "+
	//"loan_date = EXCLUDED.loan_date, "+
	//"due_date = EXCLUDED.due_date "+
	//"WHERE l.user_id <> EXCLUDED.user_id OR "+
	//"l.item_id <> EXCLUDED.item_id OR "+
	//"l.action <> EXCLUDED.action OR "+
	//"l.status_name <> EXCLUDED.status_name OR "+
	//"l.loan_date <> EXCLUDED.loan_date OR "+
	//"l.due_date <> EXCLUDED.due_date",
	//id, userId, itemId, action, statusName, loanDate, dueDate)
	//if err != nil {
	//return err
	//}

	return nil
}
*/

func sqlStartTransaction() string {
	return `
START TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;
`
}

func sqlCommit() string {
	return `
COMMIT;
`
}

func sqlDefaultSchemaLoad() string {
	return `
SET search_path = stage;
`
}

func sqlDefaultSchemaPublic() string {
	return `
SET search_path = public;
`
}

/*
func sqlTableGroupsCreate() string {
	return `
CREATE TABLE groups (
    id           UUID NOT NULL PRIMARY KEY,
    group_name   TEXT NOT NULL UNIQUE DEFAULT 'NOT AVAILABLE ' || nextval('na_groups'),
        CHECK (group_name <> ''),
    description  TEXT NOT NULL DEFAULT 'NOT AVAILABLE'
);
`
}
*/

/*
func sqlTableGroupsAddConstraints() string {
	return `
ALTER TABLE groups ALTER COLUMN id SET NOT NULL;
ALTER TABLE groups ADD PRIMARY KEY (id);
ALTER TABLE groups ALTER COLUMN group_name SET NOT NULL;
ALTER TABLE groups ADD UNIQUE (group_name);
ALTER TABLE groups ADD CONSTRAINT group_name_chk CHECK (group_name <> '');
ALTER TABLE groups ALTER COLUMN description SET NOT NULL;
`
}
*/

type stage struct {
	start              *os.File
	deleteStage        *os.File
	groups             *os.File
	usersNaGroups      *os.File
	usersNaGroupsCount int64
	users              *os.File
	mergeAll           *os.File
	end                *os.File
}
