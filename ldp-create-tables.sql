-- Extension for crosstab() requires superuser to install:
-- CREATE EXTENSION IF NOT EXISTS tablefunc;

CREATE SEQUENCE na_groups;

CREATE TABLE groups (
    id           UUID NOT NULL PRIMARY KEY,
    group_name   TEXT NOT NULL UNIQUE
            DEFAULT 'NOT AVAILABLE [' || nextval('na_groups') || ']',
	CHECK (group_name <> ''),
    description  TEXT NOT NULL DEFAULT 'NOT AVAILABLE'
);

INSERT INTO groups (id) VALUES ('00000000-0000-0000-0000-000000000000');

CREATE SEQUENCE na_users;

CREATE TABLE users (
    id               UUID NOT NULL PRIMARY KEY,
    -- username         TEXT NOT NULL UNIQUE,
    username         TEXT NOT NULL  -- TODO fix test data
            DEFAULT 'NOT AVAILABLE [' || nextval('na_users') || ']',
        CHECK (username <> ''),
    barcode          TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    user_type        TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    active           BOOLEAN NOT NULL DEFAULT FALSE,
    patron_group_id  UUID NOT NULL REFERENCES groups (id)
            DEFAULT '00000000-0000-0000-0000-000000000000'
);

INSERT INTO users (id) VALUES ('00000000-0000-0000-0000-000000000000');

CREATE TABLE loans (
    id           UUID NOT NULL PRIMARY KEY,
    user_id      UUID NOT NULL REFERENCES users (id)
            DEFAULT '00000000-0000-0000-0000-000000000000',
    item_id      UUID NOT NULL
            DEFAULT '00000000-0000-0000-0000-000000000000',
    action       TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    status_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    loan_date    TIMESTAMP NOT NULL DEFAULT 'epoch',
    due_date     TIMESTAMP NOT NULL DEFAULT 'epoch'
);

INSERT INTO loans (id) VALUES ('00000000-0000-0000-0000-000000000000');

CREATE TABLE tmp_loans_locations (
    loan_id        UUID NOT NULL PRIMARY KEY,
    location_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
        CHECK (location_name <> '')
);

INSERT INTO tmp_loans_locations (loan_id)
    VALUES ('00000000-0000-0000-0000-000000000000');

CREATE VIEW dim_users AS
SELECT u.id,
       u.username,
       u.barcode,
       u.user_type,
       u.active,
       g.group_name,
       g.description group_description
    FROM users u
        LEFT JOIN groups g ON u.patron_group_id = g.id;

CREATE VIEW dim_locations AS
SELECT 'id-' || replace(lower(tll.location_name), ' ', '-') id,
       tll.location_name
    FROM (
        SELECT DISTINCT location_name FROM tmp_loans_locations
    ) tll;

CREATE VIEW fact_loans AS
SELECT l.id,
       l.user_id,
       'id-' || replace(lower(tll.location_name), ' ', '-') location_id,
       l.item_id,
       l.action,
       l.status_name,
       l.loan_date,
       l.due_date
    FROM loans l
        LEFT JOIN tmp_loans_locations tll ON l.id = tll.loan_id;


