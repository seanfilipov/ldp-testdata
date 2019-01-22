-- Extension for crosstab() requires superuser to install:
-- CREATE EXTENSION IF NOT EXISTS tablefunc;

CREATE TABLE groups (
    id           UUID NOT NULL PRIMARY KEY,
    group_name   TEXT NOT NULL UNIQUE,
        CHECK (group_name <> ''),
    description  TEXT NOT NULL DEFAULT ''
);

CREATE TABLE users (
    id               UUID NOT NULL PRIMARY KEY,
    -- username         TEXT NOT NULL UNIQUE,
    username         TEXT NOT NULL, -- TODO fix test data
        CHECK (username <> ''),
    barcode          TEXT NOT NULL DEFAULT '',
    user_type        TEXT NOT NULL DEFAULT '',
    active           BOOLEAN NOT NULL,
    patron_group_id  UUID NOT NULL
);

CREATE TABLE loans (
    id           UUID NOT NULL PRIMARY KEY,
    user_id      UUID NOT NULL,
    item_id      UUID NOT NULL,
    action       TEXT NOT NULL DEFAULT '',
    status_name  TEXT NOT NULL DEFAULT '',
    loan_date    TIMESTAMP NOT NULL,
    due_date     TIMESTAMP NOT NULL
);

CREATE TABLE tmp_loans_locations (
    loan_id        UUID NOT NULL PRIMARY KEY,
    location_name  TEXT NOT NULL,
        CHECK (location_name <> '')
);

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


