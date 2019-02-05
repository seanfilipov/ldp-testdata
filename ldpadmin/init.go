package ldpadmin

import (
	"context"
	"database/sql"
)

type InitOptions struct {
	// Debug enables debugging output if set to true.
	Debug bool
}

func Initialize(db *sql.DB, opts *InitOptions) error {
	// Start transaction
	tx, err := db.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE SCHEMA normal;\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE SCHEMA loading;\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE loading.exlock ();\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE normal.groups (\n" +
		"    group_id     UUID NOT NULL PRIMARY KEY,\n" +
		"    group_name   TEXT NOT NULL DEFAULT 'NOT AVAILABLE',\n" +
		"	CHECK (group_name <> ''),\n" +
		"    description  TEXT NOT NULL DEFAULT 'NOT AVAILABLE'\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE loading.groups (\n" +
		"    group_id     UUID NOT NULL PRIMARY KEY,\n" +
		"    group_name   TEXT NOT NULL,\n" +
		"    description  TEXT NOT NULL\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE normal.tmp_loans_locations (\n" +
		"    loan_id        UUID NOT NULL PRIMARY KEY,\n" +
		"    location_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE',\n" +
		"        CHECK (location_name <> '')\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE loading.tmp_loans_locations (\n" +
		"    loan_id        UUID NOT NULL PRIMARY KEY,\n" +
		"    location_name  TEXT NOT NULL\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE users (\n" +
		"    user_key           BIGSERIAL NOT NULL PRIMARY KEY,\n" +
		"        CHECK (user_key > 0),\n" +
		"    user_id            UUID NOT NULL,\n" +
		"    username           TEXT NOT NULL\n" +
		"                           DEFAULT 'NOT AVAILABLE',\n" +
		"        CHECK (username <> ''),\n" +
		"    barcode            TEXT NOT NULL\n" +
		"                           DEFAULT 'NOT AVAILABLE',\n" +
		"    user_type          TEXT NOT NULL\n" +
		"                           DEFAULT 'NOT AVAILABLE',\n" +
		"    active             BOOLEAN NOT NULL\n" +
		"                           DEFAULT FALSE,\n" +
		"    group_name         TEXT NOT NULL\n" +
		"                           DEFAULT 'NOT AVAILABLE',\n" +
		"	CHECK (group_name <> ''),\n" +
		"    group_description  TEXT NOT NULL\n" +
		"                           DEFAULT 'NOT AVAILABLE',\n" +
		"    record_effective   TIMESTAMP NOT NULL\n" +
		"                           DEFAULT CURRENT_TIMESTAMP\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE INDEX ON users (user_id);\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE loading.users (\n" +
		"    user_id          UUID NOT NULL PRIMARY KEY,\n" +
		"    username         TEXT NOT NULL,\n" +
		"    barcode          TEXT NOT NULL,\n" +
		"    user_type        TEXT NOT NULL,\n" +
		"    active           BOOLEAN NOT NULL,\n" +
		"    patron_group_id  UUID NOT NULL\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE locations (\n" +
		"    location_key   TEXT NOT NULL PRIMARY KEY,\n" +
		"    location_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE'\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE loans (\n" +
		"    loan_key      BIGSERIAL NOT NULL PRIMARY KEY,\n" +
		"        CHECK (loan_key > 0),\n" +
		"    loan_id       UUID NOT NULL UNIQUE,\n" +
		"    user_key      BIGINT NOT NULL\n" +
		"                      REFERENCES users (user_key),\n" +
		"    location_key  TEXT NOT NULL\n" +
		"                      REFERENCES locations (location_key)\n" +
		"          DEFAULT '00000000-0000-0000-0000-000000000000',\n" +
		"    item_id       UUID NOT NULL\n" +
		"          DEFAULT '00000000-0000-0000-0000-000000000000',\n" +
		"    action        TEXT NOT NULL DEFAULT 'NOT AVAILABLE',\n" +
		"    status_name   TEXT NOT NULL DEFAULT 'NOT AVAILABLE',\n" +
		"    loan_date     TIMESTAMP NOT NULL DEFAULT 'epoch',\n" +
		"    due_date      TIMESTAMP NOT NULL DEFAULT 'epoch'\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE INDEX ON loans (loan_date);\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE TABLE loading.loans (\n" +
		"    loan_id       UUID NOT NULL PRIMARY KEY,\n" +
		"    user_id       UUID NOT NULL,\n" +
		"    location_key  TEXT NOT NULL,\n" +
		"    item_id       UUID NOT NULL,\n" +
		"    action        TEXT NOT NULL,\n" +
		"    status_name   TEXT NOT NULL,\n" +
		"    loan_date     TIMESTAMP NOT NULL,\n" +
		"    due_date      TIMESTAMP NOT NULL\n" +
		");\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"GRANT SELECT ON ALL TABLES IN SCHEMA public TO ldp;\n")
	if err != nil {
		return err
	}
	_, err = tx.Exec("" +
		"CREATE FUNCTION circ_detail(start_date DATE,\n" +
		"                            end_date DATE)\n" +
		"    RETURNS TABLE(location_name TEXT, group_name TEXT,\n" +
		"                  ct BIGINT)\n" +
		"    AS $$\n" +
		"SELECT loc.location_name AS location_name,\n" +
		"       u.group_name AS group_name,\n" +
		"       count(l.loan_key) AS ct\n" +
		"    FROM (\n" +
		"        SELECT loan_key, user_key, location_key\n" +
		"            FROM loans\n" +
		"            WHERE loan_date >= start_date AND\n" +
		"                  loan_date <= end_date\n" +
		"    ) l\n" +
		"        LEFT JOIN locations AS loc\n" +
		"            ON l.location_key = loc.location_key\n" +
		"        LEFT JOIN users AS u ON l.user_key = u.user_key\n" +
		"    GROUP BY loc.location_name, u.group_name\n" +
		"    ORDER BY loc.location_name, u.group_name;\n" +
		"$$\n" +
		"LANGUAGE SQL\n" +
		"IMMUTABLE;\n")
	if err != nil {
		return err
	}
	// Commit all changes
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
